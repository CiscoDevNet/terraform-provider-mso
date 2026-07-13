package mso

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	domainTypeToFamily = map[string]string{
		"physicalDomain": "physical",
		"vmmDomain":      "vmm",
	}
	dnPrefixToFamily = map[string]string{
		"uni/phys-": "physical",
		"uni/vmmp-": "vmm",
	}
	familyToIsPhysicalDomain = map[string]bool{
		"physical": true,
		"vmm":      false,
	}
	fabricInterfaceDNTemplates = map[string]string{
		"port": "topology/pod-%s/paths-%s/pathep-[%s]",
		"dpc":  "topology/pod-%s/paths-%s/pathep-[%s]",
		"vpc":  "topology/pod-%s/protpaths-%s/pathep-[%s]",
	}
	// NDO stores fabricToDeviceConnectivity.path as the full interfaceDn even when
	// the caller supplied only the bracketed segment (e.g. "eth1/20" or a policy
	// group name). The schema contract is the short form, so Read strips the DN
	// down to whatever is inside the trailing pathep-[...].
	pathepBracketRe = regexp.MustCompile(`pathep-\[([^\]]+)\]`)
)

func resourceMSOServiceDeviceClusterSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOServiceDeviceClusterSiteCreate,
		Read:   resourceMSOServiceDeviceClusterSiteRead,
		Update: resourceMSOServiceDeviceClusterSiteUpdate,
		Delete: resourceMSOServiceDeviceClusterSiteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOServiceDeviceClusterSiteImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the service device template that contains the Service Device Cluster.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the Service Device Cluster to configure on the site.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the site on which to configure the Service Device Cluster.",
			},

			"domain_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vmmDomain", "physicalDomain",
				}, false),
				RequiredWith:  []string{"domain_name"},
				ConflictsWith: []string{"domain_dn"},
				Description:   "The type of domain associated with the Service Device Cluster on the site. Must be used together with `domain_name` and cannot be combined with `domain_dn`.",
			},
			"vmm_domain_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VMware", "Microsoft", "Redhat",
				}, false),
				ConflictsWith: []string{"domain_dn"},
				Description:   "The VMM domain provider type. Required when `domain_type` is `vmmDomain` and must not be set when `domain_type` is `physicalDomain`.",
			},
			"domain_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				RequiredWith:  []string{"domain_type"},
				ConflictsWith: []string{"domain_dn"},
				Description:   "The name of the domain associated with the Service Device Cluster on the site. Must be used together with `domain_type` and cannot be combined with `domain_dn`.",
			},
			"domain_dn": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 1000),
				ConflictsWith: []string{"domain_type", "vmm_domain_type", "domain_name"},
				Description:   "The distinguished name of the domain associated with the Service Device Cluster on the site. Must start with `uni/phys-` for a physical domain or `uni/vmmp-` for a VMM domain. Cannot be combined with `domain_type`, `vmm_domain_type`, or `domain_name`.",
			},

			"high_availability_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "notAvailable",
				ValidateFunc: validation.StringInSlice([]string{
					"activeActive", "activeStandby", "notAvailable",
				}, false),
				Description: "The high availability mode of the Service Device Cluster on the site. Changing this forces the Service Device Cluster on the site to be destroyed and re-created, since transitioning between `activeActive` (per-interface domains) and the other modes (device-level domain) is a structural change to the device entry.",
			},
			"promiscuous_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether promiscuous mode is enabled on the Service Device Cluster on the site.",
			},
			"trunking_port": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the Service Device Cluster on the site uses a trunking port.",
			},
			"vlan": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 4094),
				Description:  "The VLAN ID on the Service Device Cluster on the site. Set at the device level only when `high_availability_mode` is `activeStandby`; for other modes use the interface-level `vlan` (regular L3 devices) or the per-path `vlan` inside `fabric_to_device_connectivity` (`activeActive`).",
			},

			"interfaces": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				MinItems: 1,
				// ForceNew on the outer TypeList[Resource] only triggers on count
				// changes (SDK v1 does not propagate ForceNew from the outer block
				// into inner attributes when Elem is a *Resource). Keeping it here
				// forces Destroy+Create when the interface set is reshaped
				// (add/remove/rename), which NDO server-side validation requires
				// for interface-identity changes. Content-only changes within an
				// existing interface (vlan, fabric paths, pbr_destinations, etc.)
				// flow through the Update path, which emits a wholesale
				// /interfaces replace patch.
				Description: "The list of interfaces of the Service Device Cluster on the site.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
							Description:  "The name of the interface.",
						},
						"vlan": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4094),
							Description:  "The VLAN ID of the interface. Must not be set when the matching cluster `interface_properties` binds to an `external_epg_uuid` (L3out interface); NDO rejects a VLAN on L3out interfaces.",
						},

						// Per-interface domain is used only for L1 device clusters in
						// activeActive HA mode, where NDO requires each interface to
						// carry its own physical domain. VMM domains are never valid
						// at interface scope, so the interface schema exposes only
						// domain_name (physical domain name; provider composes
						// "uni/phys-<name>") or domain_dn (literal DN, must start with
						// "uni/phys-"). Mutual exclusion and prefix validation are
						// enforced in CustomizeDiff because SDK v2 ignores
						// RequiredWith/ConflictsWith inside nested blocks.
						"domain_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
							Description:  "The name of the physical domain associated with the interface. Mutually exclusive with `domain_dn`. Only valid when `high_availability_mode` is `activeActive`; in that mode every interface must configure its own physical domain (the device-level domain attributes are derived from the first interface and any device-level values in config are ignored).",
						},
						"domain_dn": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
							Description:  "The distinguished name of the physical domain associated with the interface. Must start with `uni/phys-`. Mutually exclusive with `domain_name`. Only valid when `high_availability_mode` is `activeActive`.",
						},

						"fabric_to_device_connectivity": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of fabric-to-device connectivity paths for the interface. Allowed only when the device uses a physical domain. Mutually exclusive with `vm_information`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The pod ID of the fabric path.",
									},
									"node_id": {
										Type:        schema.TypeList,
										Required:    true,
										MinItems:    1,
										MaxItems:    2,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The node ID(s) of the fabric path. Provide a single element for `port_type` `port` and `dpc`, and two elements for `port_type` `vpc`.",
									},
									"path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The path on the node. For `port` `port_type` this is the interface (e.g. `eth1/1`). For `dpc` and `vpc` `port_type` this is the policy group name.",
									},
									"port_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"port", "vpc", "dpc",
										}, false),
										Description: "The type of port used for the fabric path.",
									},
									"tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The tag of the fabric path.",
									},
									"vlan": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 4094),
										Description:  "The VLAN ID carried on this fabric path. Used when `high_availability_mode` is `activeActive`, where each fabric path carries its own access VLAN.",
									},
								},
							},
						},

						"vm_information": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of VM information entries for the interface. Allowed only when the device uses a VMM domain. Mutually exclusive with `fabric_to_device_connectivity`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the VM.",
									},
									"vnic_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the vNIC on the VM.",
									},
									"pod_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The pod ID of the fabric path the VM interface attaches to.",
									},
									"node_id": {
										Type:        schema.TypeList,
										Optional:    true,
										MinItems:    0,
										MaxItems:    2,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The node ID(s) of the fabric path the VM interface attaches to, as a list of strings. Provide a single element for `port_type` `port` and `dpc`, and two elements for `port_type` `vpc`.",
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										// NDO server-defaults the VM interface path to a
										// DN with empty segments (e.g.
										// "topology/pod-/paths-/pathep-[]") when the caller
										// omits the pod/node/path/port_type fields. The
										// Read regex only extracts non-empty pathep-[...]
										// contents, so the raw defaulted DN leaks into
										// state and diffs against an unset HCL value on
										// every plan. Suppress the cosmetic diff when HCL
										// leaves path unset and state holds an empty
										// pathep-[] DN.
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return new == "" && strings.HasSuffix(old, "pathep-[]")
										},
										Description: "The path on the node the VM interface attaches to. For `port_type` `port` this is the interface (e.g. `eth1/1`). For `port_type` `dpc` and `vpc` this is the policy group name.",
									},
									"port_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"port", "vpc", "dpc",
										}, false),
										Description: "The type of port used for the VM interface's fabric path. Allowed values are `port`, `vpc`, `dpc`.",
									},
								},
							},
						},

						"enhanced_lag_policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the enhanced LAG policy associated with the interface. Only valid when the device uses a VMM domain.",
						},

						"pbr_destinations": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of policy-based redirect (PBR) destinations for the interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPAddress,
										Description:  "The IP address of the PBR destination. Required for L3 device clusters; omit for L1 clusters with `high_availability_mode` `activeActive` or `activeStandby`, which carry only `mac` and `tag`.",
									},
									"mac": {
										Type:     schema.TypeString,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// Newer NDO releases echo the MAC back in upper case; a
											// case-only difference between HCL and what NDO returned
											// is not a real diff.
											return strings.EqualFold(old, new)
										},
										Description: "The MAC address of the PBR destination.",
									},
									"pod_id": {
										Type:     schema.TypeString,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// NDO injects "1" for PBR destinations that omitted this field; treat that default as equivalent to unset.
											return new == "" && old == "1"
										},
										Description: "The pod ID of the PBR destination. NDO defaults this to `1` when omitted, and the provider treats that default as equivalent to leaving the attribute unset.",
									},
									"additional_tracking_ip": {
										Type:     schema.TypeString,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// NDO injects "0.0.0.0" for L3 PBR destinations that omitted this field; treat that default as equivalent to unset.
											return new == "" && old == "0.0.0.0"
										},
										Description: "The additional IP address used for tracking the PBR destination. NDO defaults this to `0.0.0.0` for L3 PBR destinations when omitted, and the provider treats that default as equivalent to leaving the attribute unset.",
									},
									"weight": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 10),
										Description:  "The weight of the PBR destination.",
									},
									"is_backup": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the PBR destination is a backup destination.",
									},
									"tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The tag of the PBR destination.",
									},
								},
							},
						},
					},
				},
			},
		},
		CustomizeDiff: resourceMSOServiceDeviceClusterSiteCustomizeDiff,
	}
}

func resourceMSOServiceDeviceClusterSiteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Beginning Import: %v", d.Id())
	if err := resourceMSOServiceDeviceClusterSiteRead(d, m); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOServiceDeviceClusterSiteCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Beginning Create")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	siteId := d.Get("site_id").(string)
	name := d.Get("name").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteId, "deviceTemplate", "sites")
	if err != nil {
		return fmt.Errorf("site %q is not present in template %q: %v", siteId, templateId, err)
	}

	siteCont := templateCont.S("deviceTemplate", "sites").Index(siteIndex)

	deviceIndex, deviceErr := GetPolicyIndexByKeyAndValue(siteCont, "name", name, "devices")
	if deviceErr == nil {
		updatePath := fmt.Sprintf("/deviceTemplate/sites/%d/devices/%d", siteIndex, deviceIndex)
		payloadCont := container.New()
		payloadCont.Array()
		appendServiceDeviceClusterSiteReplacePatches(payloadCont, updatePath, d, true)
		if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont); err != nil {
			return err
		}
	} else {
		payload := buildServiceDeviceClusterSitePayload(d)
		payloadModel := models.GetPatchPayload("add", fmt.Sprintf("/deviceTemplate/sites/%d/devices/-", siteIndex), payload)
		if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel); err != nil {
			return err
		}
	}

	d.SetId(fmt.Sprintf("templateId/%s/site/%s/ServiceDeviceCluster/%s", templateId, siteId, name))
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Create Complete: %v", d.Id())
	return resourceMSOServiceDeviceClusterSiteRead(d, m)
}

func resourceMSOServiceDeviceClusterSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, siteId, name, err := resolveServiceDeviceClusterSiteIdentity(d)
	if err != nil {
		return err
	}

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster Site - template not found, removing from state: %v", d.Id())
		d.SetId("")
		return nil
	}

	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteId, "deviceTemplate", "sites")
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster Site - site %q not found in template, removing from state: %v", siteId, d.Id())
		d.SetId("")
		return nil
	}

	siteCont := templateCont.S("deviceTemplate", "sites").Index(siteIndex)

	deviceCont, err := GetPolicyByName(siteCont, name, "devices")
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster Site - device %q not found on site, removing from state: %v", name, d.Id())
		d.SetId("")
		return nil
	}

	if err := setServiceDeviceClusterSiteData(d, deviceCont, templateId, siteId); err != nil {
		return err
	}
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOServiceDeviceClusterSiteUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	siteId := d.Get("site_id").(string)
	name := d.Get("name").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteId, "deviceTemplate", "sites")
	if err != nil {
		return fmt.Errorf("site %q is not present in template %q: %v", siteId, templateId, err)
	}

	siteCont := templateCont.S("deviceTemplate", "sites").Index(siteIndex)

	deviceIndex, err := GetPolicyIndexByKeyAndValue(siteCont, "name", name, "devices")
	if err != nil {
		return fmt.Errorf("device %q not found on site %q: %v", name, siteId, err)
	}

	updatePath := fmt.Sprintf("/deviceTemplate/sites/%d/devices/%d", siteIndex, deviceIndex)
	payloadCont := container.New()
	payloadCont.Array()
	appendServiceDeviceClusterSiteReplacePatches(payloadCont, updatePath, d, false)

	if err := doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont); err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Update Complete: %v", d.Id())
	return resourceMSOServiceDeviceClusterSiteRead(d, m)
}

func resourceMSOServiceDeviceClusterSiteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	siteId := d.Get("site_id").(string)
	name := d.Get("name").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster Site - template not found during delete: %v", d.Id())
		d.SetId("")
		return nil
	}

	siteIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "siteId", siteId, "deviceTemplate", "sites")
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster Site - site %q not found during delete: %v", siteId, d.Id())
		d.SetId("")
		return nil
	}

	siteCont := templateCont.S("deviceTemplate", "sites").Index(siteIndex)

	deviceIndex, err := GetPolicyIndexByKeyAndValue(siteCont, "name", name, "devices")
	if err != nil {
		log.Printf("[DEBUG] MSO Service Device Cluster Site - device %q not found during delete: %v", name, d.Id())
		d.SetId("")
		return nil
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/deviceTemplate/sites/%d/devices/%d", siteIndex, deviceIndex))
	if _, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO Service Device Cluster Site Resource - Delete Complete")
	return nil
}

func resourceMSOServiceDeviceClusterSiteCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	ha := d.Get("high_availability_mode").(string)

	var deviceFamily string
	if ha != "activeActive" {
		domainType, vmmDomainType, domainName, domainDN := readDeviceLevelDomainAttrs(d)
		family, err := validateDomainAttributes("", domainType, vmmDomainType, domainName, domainDN)
		if err != nil {
			return err
		}
		deviceFamily = family
	}

	interfacesRaw, ok := d.GetOk("interfaces")
	if !ok {
		return nil
	}
	for _, rawInterface := range interfacesRaw.([]interface{}) {
		interfaceData := rawInterface.(map[string]interface{})
		name, _ := interfaceData["name"].(string)

		// activeActive HA pins the interface family to physical (VMM domains
		// are not valid at interface scope), so effectiveFamily and
		// familyScope are constants in that branch. In every other HA mode
		// the interface inherits the device-level domain family.
		effectiveFamily := deviceFamily
		familyScope := "device"
		if ha == "activeActive" {
			ifaceDomainName, _ := interfaceData["domain_name"].(string)
			ifaceDomainDN, _ := interfaceData["domain_dn"].(string)
			if ifaceDomainName == "" && ifaceDomainDN == "" {
				return fmt.Errorf("interface %q: domain must be configured on every interface when high_availability_mode is \"activeActive\" (set domain_name or domain_dn)", name)
			}
			if err := validateInterfaceDomainAttributes(name, ifaceDomainName, ifaceDomainDN); err != nil {
				return err
			}
			effectiveFamily = "physical"
			familyScope = "interface"
		}

		hasFabric := collectionHasItems(interfaceData["fabric_to_device_connectivity"])
		hasVMM := collectionHasItems(interfaceData["vm_information"])

		if hasFabric && hasVMM {
			return fmt.Errorf("interface %q: fabric_to_device_connectivity and vm_information are mutually exclusive", name)
		}
		if !hasFabric && !hasVMM {
			return fmt.Errorf("interface %q: one of fabric_to_device_connectivity or vm_information must be set", name)
		}

		switch effectiveFamily {
		case "physical":
			if hasVMM {
				return fmt.Errorf("interface %q: vm_information is not allowed when the %s uses a physicalDomain", name, familyScope)
			}
			if enhancedLagPolicy, _ := interfaceData["enhanced_lag_policy"].(string); enhancedLagPolicy != "" {
				return fmt.Errorf("interface %q: enhanced_lag_policy is not allowed when the %s uses a physicalDomain", name, familyScope)
			}
		case "vmm":
			if hasFabric {
				return fmt.Errorf("interface %q: fabric_to_device_connectivity is not allowed when the %s uses a vmmDomain", name, familyScope)
			}
		}

		if hasFabric {
			for _, rawFabricPath := range interfaceData["fabric_to_device_connectivity"].([]interface{}) {
				fabricPathData := rawFabricPath.(map[string]interface{})
				portType, _ := fabricPathData["port_type"].(string)
				nodeIDs := nodeIDsFromAttr(fabricPathData["node_id"])
				if err := validateFabricPortTypeNodeCount(name, portType, nodeIDs); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// validateFabricPortTypeNodeCount enforces the node_id list length that NDO
// requires for each port_type on a fabric_to_device_connectivity entry.
func validateFabricPortTypeNodeCount(interfaceName, portType string, nodeIDs []string) error {
	switch portType {
	case "vpc":
		if len(nodeIDs) != 2 {
			return fmt.Errorf("interface %q: port_type \"vpc\" requires exactly two node_id entries, got %d", interfaceName, len(nodeIDs))
		}
	case "port", "dpc":
		if len(nodeIDs) != 1 {
			return fmt.Errorf("interface %q: port_type %q requires exactly one node_id entry, got %d", interfaceName, portType, len(nodeIDs))
		}
	}
	return nil
}

func collectionHasItems(value interface{}) bool {
	switch v := value.(type) {
	case []interface{}:
		return len(v) > 0
	case *schema.Set:
		return v != nil && v.Len() > 0
	}
	return false
}

func resolveServiceDeviceClusterSiteIdentity(d *schema.ResourceData) (string, string, string, error) {
	if d.Id() == "" {
		return d.Get("template_id").(string), d.Get("site_id").(string), d.Get("name").(string), nil
	}
	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return "", "", "", err
	}
	siteId, err := GetPolicyNameFromResourceId(d.Id(), "site")
	if err != nil {
		return "", "", "", err
	}
	name, err := GetPolicyNameFromResourceId(d.Id(), "ServiceDeviceCluster")
	if err != nil {
		return "", "", "", err
	}
	return templateId, siteId, name, nil
}

func composeDomainDn(domainType, vmmDomainType, domainName string) string {
	switch domainType {
	case "physicalDomain":
		return fmt.Sprintf("uni/phys-%s", domainName)
	case "vmmDomain":
		return fmt.Sprintf("uni/vmmp-%s/dom-%s", vmmDomainType, domainName)
	}
	return ""
}

// composeFabricInterfaceDn assembles the NDO interfaceDn URL form. NDO uses a
// hyphen-separated node segment in the URL path (paths-101 or protpaths-101-102),
// which is different from the comma-separated form NDO uses for the JSON nodeID
// field on the same payload entry.
func composeFabricInterfaceDn(podID string, nodeIDs []string, path, portType string) string {
	tmpl, ok := fabricInterfaceDNTemplates[portType]
	if !ok {
		return ""
	}
	return fmt.Sprintf(tmpl, podID, strings.Join(nodeIDs, "-"), path)
}

// nodeIDsFromAttr normalises a TypeList[TypeString] value (which arrives as
// []interface{} from the SDK) into a []string.
func nodeIDsFromAttr(raw interface{}) []string {
	items, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func resolveDomainFamily(domainType, domainDN string) string {
	if f, ok := domainTypeToFamily[domainType]; ok {
		return f
	}
	for prefix, f := range dnPrefixToFamily {
		if strings.HasPrefix(domainDN, prefix) {
			return f
		}
	}
	return ""
}

// validateDomainAttributes runs the positive validation rules shared by the device-
// level domain attribute group and the per-interface domain attribute group. scope
// is prefixed onto error messages (use "" for device-level, "interface \"name\": "
// for per-interface). Returns the resolved family ("physical"/"vmm") on success.
//
// Some of these rules are also enforced at the device level by the schema's
// RequiredWith/ConflictsWith; we re-state them here so the per-interface attribute
// group (where the SDK ignores RequiredWith/ConflictsWith inside nested blocks)
// gets the same enforcement.
func validateDomainAttributes(scope, domainType, vmmDomainType, domainName, domainDN string) (string, error) {
	if domainDN != "" && (domainType != "" || vmmDomainType != "" || domainName != "") {
		// Schema ConflictsWith blocks user-supplied conflicts at the device level.
		// Here we tolerate the redundant set when both forms resolve to the same
		// DN, which happens after Read populates both the literal DN and the
		// decomposed triple and the prior state propagates back through Computed
		// on the next plan.
		if composeDomainDn(domainType, vmmDomainType, domainName) != domainDN {
			return "", fmt.Errorf("%sdomain_dn is mutually exclusive with domain_type, vmm_domain_type, and domain_name", scope)
		}
	}
	if domainName != "" && domainType == "" {
		return "", fmt.Errorf("%sdomain_name requires domain_type to be set", scope)
	}
	if vmmDomainType != "" && domainType == "physicalDomain" {
		return "", fmt.Errorf("%svmm_domain_type must not be set when domain_type is %q", scope, "physicalDomain")
	}
	if vmmDomainType != "" && domainType != "vmmDomain" {
		return "", fmt.Errorf("%svmm_domain_type requires domain_type to be set to %q", scope, "vmmDomain")
	}
	switch domainType {
	case "vmmDomain":
		if vmmDomainType == "" {
			return "", fmt.Errorf("%svmm_domain_type is required when domain_type is %q", scope, "vmmDomain")
		}
		if domainName == "" {
			return "", fmt.Errorf("%sdomain_name is required when domain_type is set", scope)
		}
	case "physicalDomain":
		if domainName == "" {
			return "", fmt.Errorf("%sdomain_name is required when domain_type is set", scope)
		}
	}
	if domainType == "" && domainDN == "" {
		return "", fmt.Errorf("%sone of domain_type or domain_dn must be set", scope)
	}
	family := resolveDomainFamily(domainType, domainDN)
	if family == "" {
		return "", fmt.Errorf("%sdomain_dn %q must start with %q or %q", scope, domainDN, "uni/phys-", "uni/vmmp-")
	}
	return family, nil
}

// domainSource is the minimal subset of *schema.ResourceData / *schema.ResourceDiff
// needed by the HA-aware domain resolver, so it can be reused from CustomizeDiff,
// payload builders, and the patch emitter.
type domainSource interface {
	Get(key string) interface{}
	HasChange(key string) bool
}

// readDeviceLevelDomainAttrs returns the four device-level domain attribute
// values to feed into validateDomainAttributes / resolveDomainFromValues, with
// stale Computed back-fills filtered out.
//
// All four attrs are Optional+Computed and Read writes them every refresh, so
// on an update the side the user did NOT supply in config back-fills from
// prior state. When the user switches the domain via only one form
// (e.g. drops the type/name pair and sets a fresh domain_dn, or vice versa)
// the stale Computed back-fill from the other form still points at the
// previous domain. Without filtering, validateDomainAttributes' mutual-
// exclusion check trips on a value the user did not write, and
// resolveDomainFromValues' "literal DN wins" rule submits the stale DN.
//
// The filter keys off d.HasChange: if only one form changed in this plan,
// the other form's d.Get value is the stale back-fill and is zeroed out.
func readDeviceLevelDomainAttrs(d domainSource) (domainType, vmmDomainType, domainName, domainDN string) {
	domainType = d.Get("domain_type").(string)
	vmmDomainType = d.Get("vmm_domain_type").(string)
	domainName = d.Get("domain_name").(string)
	domainDN = d.Get("domain_dn").(string)
	dnChanged := d.HasChange("domain_dn")
	tripleChanged := d.HasChange("domain_type") || d.HasChange("vmm_domain_type") || d.HasChange("domain_name")
	switch {
	case dnChanged && !tripleChanged:
		domainType, vmmDomainType, domainName = "", "", ""
	case tripleChanged && !dnChanged:
		domainDN = ""
	}
	return
}

// resolveDomainFromValues applies the prefer-literal-DN-else-compose rule and
// returns (resolvedDomainDN, family). The composition rule reflects how NDO
// stores `domainDn` in JSON: callers may supply either a literal DN or the
// type/name/vmm triple, and we always submit the literal DN form.
func resolveDomainFromValues(domainType, vmmDomainType, domainName, domainDN string) (string, string) {
	if domainDN == "" {
		domainDN = composeDomainDn(domainType, vmmDomainType, domainName)
	}
	return domainDN, resolveDomainFamily(domainType, domainDN)
}

// resolveInterfaceDomain returns the resolved domain DN for one interface
// map. Interface-scoped domains are only ever physical (activeActive HA / L1),
// so the caller can assume family = "physical" and this helper deals only
// with the literal-DN-else-compose rule.
func resolveInterfaceDomain(interfaceData map[string]interface{}) string {
	domainName, _ := interfaceData["domain_name"].(string)
	domainDN, _ := interfaceData["domain_dn"].(string)
	if domainDN != "" {
		return domainDN
	}
	if domainName != "" {
		return fmt.Sprintf("uni/phys-%s", domainName)
	}
	return ""
}

// validateInterfaceDomainAttributes enforces mutual exclusion between
// domain_name and domain_dn and the "uni/phys-" prefix on domain_dn. Called
// only in activeActive HA mode where per-interface domain applies.
func validateInterfaceDomainAttributes(interfaceName, domainName, domainDN string) error {
	if domainName != "" && domainDN != "" {
		// Tolerate the redundant set when both forms resolve to the same DN;
		// this can happen after Read populates both the literal DN and the
		// name, and the prior state propagates back through Computed on the
		// next plan.
		if fmt.Sprintf("uni/phys-%s", domainName) != domainDN {
			return fmt.Errorf("interface %q: domain_name and domain_dn are mutually exclusive", interfaceName)
		}
	}
	if domainDN != "" && !strings.HasPrefix(domainDN, "uni/phys-") {
		return fmt.Errorf("interface %q: domain_dn %q must start with %q (only physical domains are supported per-interface)", interfaceName, domainDN, "uni/phys-")
	}
	return nil
}

// resolveDeviceLevelDomain returns (resolvedDomainDN, family) to send to NDO at the
// device level. In activeActive HA mode the device-level domain is derived from the
// first interface's domain because the user supplies per-interface domains only;
// any device-level inputs in config are ignored. Interface-scoped domains are
// always physical, so the family is a constant "physical" in that branch. In
// every other HA mode the device-level domain is taken directly from the
// top-level inputs, with stale Computed back-fills filtered (see
// readDeviceLevelDomainAttrs).
func resolveDeviceLevelDomain(d domainSource) (string, string) {
	if d.Get("high_availability_mode").(string) == "activeActive" {
		interfaces, _ := d.Get("interfaces").([]interface{})
		if len(interfaces) > 0 {
			if first, ok := interfaces[0].(map[string]interface{}); ok {
				if dn := resolveInterfaceDomain(first); dn != "" {
					return dn, "physical"
				}
			}
		}
		return "", ""
	}
	domainType, vmmDomainType, domainName, domainDN := readDeviceLevelDomainAttrs(d)
	return resolveDomainFromValues(domainType, vmmDomainType, domainName, domainDN)
}

func buildServiceDeviceClusterSitePayload(d *schema.ResourceData) map[string]interface{} {
	payload := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"highAvailabilityMode": d.Get("high_availability_mode").(string),
		"promiscuousMode":      d.Get("promiscuous_mode").(bool),
		"trunkingPort":         d.Get("trunking_port").(bool),
		"interfaces":           buildServiceDeviceClusterSiteInterfacesPayload(d),
	}
	// Domain emission at device scope:
	//   * non-activeActive HA: emit both `isPhysicalDomain` and `domainDn`.
	//   * activeActive HA: NDO requires the domain in INTERFACE scope only;
	//     emitting `domainDn` at device scope fails server-side validation with
	//     "Physical Domain in device scope is not supported for L1 in
	//     Active/Active HA Mode. Instead, a domain in interface scope needs to
	//     be configured." Emit only `isPhysicalDomain` (derived from
	//     interfaces[0]'s family so the device-level flag stays consistent).
	if resolvedDomainDN, family := resolveDeviceLevelDomain(d); resolvedDomainDN != "" {
		payload["isPhysicalDomain"] = familyToIsPhysicalDomain[family]
		if d.Get("high_availability_mode").(string) != "activeActive" {
			payload["domainDn"] = resolvedDomainDN
		}
	}
	if vlan, ok := d.GetOk("vlan"); ok {
		payload["vlan"] = vlan.(int)
	}
	return payload
}

func buildServiceDeviceClusterSiteInterfacesPayload(d *schema.ResourceData) []map[string]interface{} {
	interfaces, ok := d.Get("interfaces").([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	haIsActiveActive := d.Get("high_availability_mode").(string) == "activeActive"
	payload := make([]map[string]interface{}, 0, len(interfaces))
	for _, rawInterface := range interfaces {
		interfaceData := rawInterface.(map[string]interface{})
		vlan := interfaceData["vlan"].(int)
		entry := map[string]interface{}{
			"name": interfaceData["name"].(string),
		}
		if vlan > 0 {
			entry["vlan"] = vlan
		}
		// Per-interface domainDn is only emitted in activeActive HA mode. In other
		// modes, per-interface domain attrs in config are ignored so the device-level
		// domain remains the single source of truth.
		if haIsActiveActive {
			if ifaceDomainDN := resolveInterfaceDomain(interfaceData); ifaceDomainDN != "" {
				entry["domainDn"] = ifaceDomainDN
			}
		}
		if fabricList, ok := interfaceData["fabric_to_device_connectivity"].([]interface{}); ok && len(fabricList) > 0 {
			entry["fabricToDeviceConnectivity"] = buildFabricToDeviceConnectivityPayload(fabricList)
		}
		if vmList, ok := interfaceData["vm_information"].([]interface{}); ok && len(vmList) > 0 {
			entry["vmmIntfInfo"] = buildVMMIntfInfoPayload(vmList)
		}
		if pbrList, ok := interfaceData["pbr_destinations"].([]interface{}); ok && len(pbrList) > 0 {
			entry["pbrDestinations"] = buildPbrDestinationsPayload(pbrList)
		}
		if enhancedLagPolicy, _ := interfaceData["enhanced_lag_policy"].(string); enhancedLagPolicy != "" {
			entry["enhancedLagPolicy"] = enhancedLagPolicy
		}
		payload = append(payload, entry)
	}
	return payload
}

func buildFabricToDeviceConnectivityPayload(fabricPaths []interface{}) []map[string]interface{} {
	payload := make([]map[string]interface{}, 0, len(fabricPaths))
	for _, rawFabricPath := range fabricPaths {
		fabricPathData := rawFabricPath.(map[string]interface{})
		podID := fabricPathData["pod_id"].(string)
		nodeIDs := nodeIDsFromAttr(fabricPathData["node_id"])
		path := fabricPathData["path"].(string)
		portType := fabricPathData["port_type"].(string)
		interfaceDn := composeFabricInterfaceDn(podID, nodeIDs, path, portType)
		entry := map[string]interface{}{
			"podID":       podID,
			"nodeID":      strings.Join(nodeIDs, ","),
			"path":        interfaceDn,
			"portType":    portType,
			"interfaceDn": interfaceDn,
		}
		if tag, _ := fabricPathData["tag"].(string); tag != "" {
			entry["tag"] = tag
		}
		if pathVlan, ok := fabricPathData["vlan"].(int); ok && pathVlan > 0 {
			entry["vlan"] = pathVlan
		}
		payload = append(payload, entry)
	}
	return payload
}

func buildVMMIntfInfoPayload(vms []interface{}) []map[string]interface{} {
	payload := make([]map[string]interface{}, 0, len(vms))
	for _, rawVM := range vms {
		vmData := rawVM.(map[string]interface{})
		entry := map[string]interface{}{
			"vmName":   vmData["vm_name"].(string),
			"vNicName": vmData["vnic_name"].(string),
		}
		// Fabric attach point on the VMM interface. Emit whenever any field is
		// set so partial configs surface as NDO server-side errors rather than
		// being silently dropped; bare {vm_name, vnic_name} entries stay
		// backwards-compatible.
		podID, _ := vmData["pod_id"].(string)
		nodeIDs := nodeIDsFromAttr(vmData["node_id"])
		path, _ := vmData["path"].(string)
		portType, _ := vmData["port_type"].(string)
		if podID != "" || len(nodeIDs) > 0 || path != "" || portType != "" {
			interfaceDn := composeFabricInterfaceDn(podID, nodeIDs, path, portType)
			entry["podID"] = podID
			entry["nodeID"] = strings.Join(nodeIDs, ",")
			entry["path"] = interfaceDn
			entry["portType"] = portType
			entry["interfaceDn"] = interfaceDn
		}
		payload = append(payload, entry)
	}
	return payload
}

func buildPbrDestinationsPayload(destinations []interface{}) []map[string]interface{} {
	payload := make([]map[string]interface{}, 0, len(destinations))
	for _, rawDestination := range destinations {
		destinationData := rawDestination.(map[string]interface{})
		entry := map[string]interface{}{
			// TODO: derive isAdvancedConfigSet from advanced fields (ip, mac, additionalTrackingIP, weight, isBackUp, tag).
			//       For now always false until the NDO UI / API behaviour for this flag is confirmed.
			"isAdvancedConfigSet": false,
		}
		if ip, _ := destinationData["ip"].(string); ip != "" {
			entry["ip"] = ip
		}
		if mac, _ := destinationData["mac"].(string); mac != "" {
			entry["mac"] = mac
		}
		if podID, _ := destinationData["pod_id"].(string); podID != "" {
			entry["podID"] = podID
		}
		if additionalTrackingIP, _ := destinationData["additional_tracking_ip"].(string); additionalTrackingIP != "" {
			entry["additionalTrackingIP"] = additionalTrackingIP
		}
		if weight, ok := destinationData["weight"].(int); ok && weight > 0 {
			entry["weight"] = weight
		}
		if isBackup, ok := destinationData["is_backup"].(bool); ok {
			entry["isBackUp"] = isBackup
		}
		if tag, _ := destinationData["tag"].(string); tag != "" {
			entry["tag"] = tag
		}
		payload = append(payload, entry)
	}
	return payload
}

// appendServiceDeviceClusterSiteReplacePatches emits one JSON-patch "replace" op per
// owned field. forceAll=true emits every field (used by Create when overwriting an
// existing device entry); forceAll=false emits only fields with d.HasChange (Update).
//
// `high_availability_mode` is ForceNew, so any HA-mode change triggers
// Destroy+Create and never reaches this Update path — /highAvailabilityMode is
// therefore only emitted in the forceAll case. `interfaces` is also marked
// ForceNew on the outer block, but SDK v1 only propagates that to the
// interfaces.# count diff; content-only changes within an existing interface
// (vlan, fabric paths, pbr_destinations, per-interface domain, etc.) reach this
// Update path, so /interfaces is replaced wholesale whenever d.HasChange("interfaces")
// is true.
func appendServiceDeviceClusterSiteReplacePatches(payloadCont *container.Container, updatePath string, d *schema.ResourceData, forceAll bool) {
	domainGroupChanged := forceAll ||
		d.HasChange("domain_type") || d.HasChange("vmm_domain_type") || d.HasChange("domain_name") || d.HasChange("domain_dn")
	if domainGroupChanged {
		if resolvedDomainDN, family := resolveDeviceLevelDomain(d); resolvedDomainDN != "" {
			addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/isPhysicalDomain", updatePath), familyToIsPhysicalDomain[family])
			// activeActive HA: NDO rejects /domainDn at device scope (must be
			// per-interface). Mirror the Create-side rule and skip the device-
			// scope /domainDn patch; per-interface domainDn flows through the
			// /interfaces wholesale replace below.
			if d.Get("high_availability_mode").(string) != "activeActive" {
				addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/domainDn", updatePath), resolvedDomainDN)
			}
		}
	}
	if forceAll {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/highAvailabilityMode", updatePath), d.Get("high_availability_mode").(string))
	}
	if forceAll || d.HasChange("promiscuous_mode") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/promiscuousMode", updatePath), d.Get("promiscuous_mode").(bool))
	}
	if forceAll || d.HasChange("trunking_port") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/trunkingPort", updatePath), d.Get("trunking_port").(bool))
	}
	if forceAll || d.HasChange("vlan") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/vlan", updatePath), d.Get("vlan").(int))
	}
	if forceAll || d.HasChange("interfaces") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/interfaces", updatePath), buildServiceDeviceClusterSiteInterfacesPayload(d))
	}
}

func setServiceDeviceClusterSiteData(d *schema.ResourceData, deviceCont *container.Container, templateId, siteId string) error {
	name := models.StripQuotes(deviceCont.S("name").String())
	d.SetId(fmt.Sprintf("templateId/%s/site/%s/ServiceDeviceCluster/%s", templateId, siteId, name))
	d.Set("template_id", templateId)
	d.Set("site_id", siteId)
	d.Set("name", name)

	domainDn := models.StripQuotes(deviceCont.S("domainDn").String())
	if domainDn == "{}" {
		domainDn = ""
	}
	d.Set("domain_dn", domainDn)
	domainType, vmmDomainType, domainName := decomposeDomainDn(domainDn)
	d.Set("domain_type", domainType)
	d.Set("vmm_domain_type", vmmDomainType)
	d.Set("domain_name", domainName)

	if deviceCont.Exists("highAvailabilityMode") {
		v := models.StripQuotes(deviceCont.S("highAvailabilityMode").String())
		if v == "{}" {
			v = ""
		}
		d.Set("high_availability_mode", v)
	}
	if deviceCont.Exists("promiscuousMode") {
		if v, ok := deviceCont.S("promiscuousMode").Data().(bool); ok {
			d.Set("promiscuous_mode", v)
		}
	}
	if deviceCont.Exists("trunkingPort") {
		if v, ok := deviceCont.S("trunkingPort").Data().(bool); ok {
			d.Set("trunking_port", v)
		}
	}
	if deviceCont.Exists("vlan") {
		if v, ok := deviceCont.S("vlan").Data().(float64); ok {
			d.Set("vlan", int(v))
		}
	}

	interfaceConts, err := deviceCont.S("interfaces").Children()
	if err != nil {
		d.Set("interfaces", []interface{}{})
		return nil
	}
	interfaceList := make([]map[string]interface{}, 0, len(interfaceConts))
	for _, interfaceCont := range interfaceConts {
		entry := map[string]interface{}{
			"name": models.StripQuotes(interfaceCont.S("name").String()),
		}
		if interfaceCont.Exists("vlan") {
			if vlan, ok := interfaceCont.S("vlan").Data().(float64); ok {
				entry["vlan"] = int(vlan)
			}
		}
		if interfaceCont.Exists("domainDn") {
			ifaceDomainDn := models.StripQuotes(interfaceCont.S("domainDn").String())
			if ifaceDomainDn == "{}" {
				ifaceDomainDn = ""
			}
			entry["domain_dn"] = ifaceDomainDn
			// Interface-scoped domains are always physical (activeActive HA),
			// so decomposing the DN reduces to stripping the "uni/phys-" prefix.
			// A non-conforming DN cannot round-trip and leaves domain_name empty.
			if strings.HasPrefix(ifaceDomainDn, "uni/phys-") {
				entry["domain_name"] = strings.TrimPrefix(ifaceDomainDn, "uni/phys-")
			} else {
				entry["domain_name"] = ""
			}
		}
		if interfaceCont.Exists("enhancedLagPolicy") {
			enhancedLagPolicy := models.StripQuotes(interfaceCont.S("enhancedLagPolicy").String())
			if enhancedLagPolicy == "{}" {
				enhancedLagPolicy = ""
			}
			entry["enhanced_lag_policy"] = enhancedLagPolicy
		}
		if interfaceCont.Exists("fabricToDeviceConnectivity") {
			if fabricPathConts, err := interfaceCont.S("fabricToDeviceConnectivity").Children(); err == nil {
				fabricPaths := make([]map[string]interface{}, 0, len(fabricPathConts))
				for _, fabricPathCont := range fabricPathConts {
					fabricPathData := map[string]interface{}{}
					if fabricPathCont.Exists("podID") {
						fabricPathData["pod_id"] = models.StripQuotes(fabricPathCont.S("podID").String())
					}
					if fabricPathCont.Exists("nodeID") {
						// NDO encodes vpc nodes as a comma-separated string in the JSON nodeID field
						// ("101,102"), even though it uses hyphens in the corresponding URL path
						// segment (protpaths-101-102). Split back to a list for state.
						rawNodeID := models.StripQuotes(fabricPathCont.S("nodeID").String())
						var nodeIDs []interface{}
						if rawNodeID != "" && rawNodeID != "{}" {
							for _, n := range strings.Split(rawNodeID, ",") {
								nodeIDs = append(nodeIDs, n)
							}
						}
						fabricPathData["node_id"] = nodeIDs
					}
					if fabricPathCont.Exists("path") {
						rawPath := models.StripQuotes(fabricPathCont.S("path").String())
						if matches := pathepBracketRe.FindStringSubmatch(rawPath); len(matches) >= 2 {
							fabricPathData["path"] = matches[1]
						} else {
							fabricPathData["path"] = rawPath
						}
					}
					if fabricPathCont.Exists("portType") {
						fabricPathData["port_type"] = models.StripQuotes(fabricPathCont.S("portType").String())
					}
					if fabricPathCont.Exists("tag") {
						fabricPathData["tag"] = models.StripQuotes(fabricPathCont.S("tag").String())
					}
					if fabricPathCont.Exists("vlan") {
						if pathVlan, ok := fabricPathCont.S("vlan").Data().(float64); ok {
							fabricPathData["vlan"] = int(pathVlan)
						}
					}
					fabricPaths = append(fabricPaths, fabricPathData)
				}
				entry["fabric_to_device_connectivity"] = fabricPaths
			}
		}
		if interfaceCont.Exists("vmmIntfInfo") {
			if vmConts, err := interfaceCont.S("vmmIntfInfo").Children(); err == nil {
				vms := make([]map[string]interface{}, 0, len(vmConts))
				for _, vmCont := range vmConts {
					vmData := map[string]interface{}{}
					if vmCont.Exists("vmName") {
						vmData["vm_name"] = models.StripQuotes(vmCont.S("vmName").String())
					}
					if vmCont.Exists("vNicName") {
						vmData["vnic_name"] = models.StripQuotes(vmCont.S("vNicName").String())
					}
					if vmCont.Exists("podID") {
						vmData["pod_id"] = models.StripQuotes(vmCont.S("podID").String())
					}
					if vmCont.Exists("nodeID") {
						// Same encoding as fabricToDeviceConnectivity: comma-separated in JSON
						// ("101,102" for vpc), hyphen-separated in the corresponding URL path.
						rawNodeID := models.StripQuotes(vmCont.S("nodeID").String())
						var nodeIDs []interface{}
						if rawNodeID != "" && rawNodeID != "{}" {
							for _, n := range strings.Split(rawNodeID, ",") {
								nodeIDs = append(nodeIDs, n)
							}
						}
						vmData["node_id"] = nodeIDs
					}
					if vmCont.Exists("path") {
						rawPath := models.StripQuotes(vmCont.S("path").String())
						if matches := pathepBracketRe.FindStringSubmatch(rawPath); len(matches) >= 2 {
							vmData["path"] = matches[1]
						} else {
							vmData["path"] = rawPath
						}
					}
					if vmCont.Exists("portType") {
						vmData["port_type"] = models.StripQuotes(vmCont.S("portType").String())
					}
					vms = append(vms, vmData)
				}
				entry["vm_information"] = vms
			}
		}
		if interfaceCont.Exists("pbrDestinations") {
			if destinationConts, err := interfaceCont.S("pbrDestinations").Children(); err == nil {
				destinations := make([]map[string]interface{}, 0, len(destinationConts))
				for _, destinationCont := range destinationConts {
					destinationData := map[string]interface{}{}
					if destinationCont.Exists("ip") {
						destinationData["ip"] = models.StripQuotes(destinationCont.S("ip").String())
					}
					if destinationCont.Exists("mac") {
						destinationData["mac"] = models.StripQuotes(destinationCont.S("mac").String())
					}
					if destinationCont.Exists("podID") {
						destinationData["pod_id"] = models.StripQuotes(destinationCont.S("podID").String())
					}
					if destinationCont.Exists("additionalTrackingIP") {
						destinationData["additional_tracking_ip"] = models.StripQuotes(destinationCont.S("additionalTrackingIP").String())
					}
					if destinationCont.Exists("weight") {
						if weight, ok := destinationCont.S("weight").Data().(float64); ok {
							destinationData["weight"] = int(weight)
						}
					}
					if destinationCont.Exists("isBackUp") {
						if isBackup, ok := destinationCont.S("isBackUp").Data().(bool); ok {
							destinationData["is_backup"] = isBackup
						}
					}
					if destinationCont.Exists("tag") {
						destinationData["tag"] = models.StripQuotes(destinationCont.S("tag").String())
					}
					destinations = append(destinations, destinationData)
				}
				entry["pbr_destinations"] = destinations
			}
		}
		interfaceList = append(interfaceList, entry)
	}
	d.Set("interfaces", interfaceList)
	return nil
}

// decomposeDomainDn splits an APIC domain DN into (domainType, vmmDomainType, domainName).
// Returns ("", "", "") for an unrecognised DN.
func decomposeDomainDn(dn string) (string, string, string) {
	switch {
	case strings.HasPrefix(dn, "uni/phys-"):
		return "physicalDomain", "", strings.TrimPrefix(dn, "uni/phys-")
	case strings.HasPrefix(dn, "uni/vmmp-"):
		rest := strings.TrimPrefix(dn, "uni/vmmp-")
		i := strings.Index(rest, "/dom-")
		if i < 0 {
			return "vmmDomain", rest, ""
		}
		return "vmmDomain", rest[:i], rest[i+len("/dom-"):]
	}
	return "", "", ""
}
