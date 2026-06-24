package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Default:  "notAvailable",
				ValidateFunc: validation.StringInSlice([]string{
					"activeActive", "activeStandby", "notAvailable",
				}, false),
				Description: "The high availability mode of the Service Device Cluster on the site.",
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

			"interfaces": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				MinItems: 1,
				// ForceNew on the outer block is load-bearing for two
				// reasons: NDO server-side validation rejects most in-place
				// reshapes of the interface array (the device entry must
				// be torn down and rebuilt to safely change interface
				// identity), and ForceNew also makes the Optional+Computed
				// attributes inside the block safe — every change goes
				// through Destroy+Create+Read so the prior state can never
				// leak into a new slot's Apply input the way it did on the
				// in-place mso_service_device_cluster TypeList.
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
								},
							},
						},

						"vm_information": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The set of VM information entries for the interface. Allowed only when the device uses a VMM domain. Mutually exclusive with `fabric_to_device_connectivity`.",
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
								},
							},
						},

						"enhanced_lag_policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The UUID of the enhanced LAG policy associated with the interface. Only valid when the device uses a VMM domain.",
						},

						"pbr_destinations": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of policy-based redirect (PBR) destinations for the interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPAddress,
										Description:  "The IP address of the PBR destination.",
									},
									"mac": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The MAC address of the PBR destination.",
									},
									"pod_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The pod ID of the PBR destination.",
									},
									"additional_tracking_ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The additional IP address used for tracking the PBR destination.",
									},
									"weight": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IntBetween(1, 10),
										Description:  "The weight of the PBR destination.",
									},
									"is_backup": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Whether the PBR destination is a backup destination.",
									},
									"tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
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

func resourceMSOServiceDeviceClusterSiteCustomizeDiff(d *schema.ResourceDiff, _ interface{}) error {
	domainType := d.Get("domain_type").(string)
	domainDN := d.Get("domain_dn").(string)
	vmmDomainType := d.Get("vmm_domain_type").(string)

	if domainType == "" && domainDN == "" {
		return fmt.Errorf("one of domain_type or domain_dn must be set")
	}

	switch domainType {
	case "vmmDomain":
		if vmmDomainType == "" {
			return fmt.Errorf("vmm_domain_type is required when domain_type is \"vmmDomain\"")
		}
	case "physicalDomain":
		if vmmDomainType != "" {
			return fmt.Errorf("vmm_domain_type must not be set when domain_type is \"physicalDomain\"")
		}
	}

	family := resolveDomainFamily(domainType, domainDN)
	if family == "" {
		return fmt.Errorf("domain_dn %q must start with %q or %q", domainDN, "uni/phys-", "uni/vmmp-")
	}

	interfacesRaw, ok := d.GetOk("interfaces")
	if !ok {
		return nil
	}
	for _, rawInterface := range interfacesRaw.([]interface{}) {
		interfaceData := rawInterface.(map[string]interface{})
		name, _ := interfaceData["name"].(string)

		hasFabric := collectionHasItems(interfaceData["fabric_to_device_connectivity"])
		hasVMM := collectionHasItems(interfaceData["vm_information"])

		if hasFabric && hasVMM {
			return fmt.Errorf("interface %q: fabric_to_device_connectivity and vm_information are mutually exclusive", name)
		}
		if !hasFabric && !hasVMM {
			return fmt.Errorf("interface %q: one of fabric_to_device_connectivity or vm_information must be set", name)
		}

		switch family {
		case "physical":
			if hasVMM {
				return fmt.Errorf("interface %q: vm_information is not allowed when the device uses a physicalDomain", name)
			}
			if enhancedLagPolicy, _ := interfaceData["enhanced_lag_policy"].(string); enhancedLagPolicy != "" {
				return fmt.Errorf("interface %q: enhanced_lag_policy is not allowed when the device uses a physicalDomain", name)
			}
		case "vmm":
			if hasFabric {
				return fmt.Errorf("interface %q: fabric_to_device_connectivity is not allowed when the device uses a vmmDomain", name)
			}
		}

		if hasFabric {
			for _, rawFabricPath := range interfaceData["fabric_to_device_connectivity"].([]interface{}) {
				fabricPathData := rawFabricPath.(map[string]interface{})
				portType, _ := fabricPathData["port_type"].(string)
				nodeIDs := nodeIDsFromAttr(fabricPathData["node_id"])
				switch portType {
				case "vpc":
					if len(nodeIDs) != 2 {
						return fmt.Errorf("interface %q: port_type \"vpc\" requires exactly two node_id entries, got %d", name, len(nodeIDs))
					}
				case "port", "dpc":
					if len(nodeIDs) != 1 {
						return fmt.Errorf("interface %q: port_type %q requires exactly one node_id entry, got %d", name, portType, len(nodeIDs))
					}
				}
			}
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

func buildServiceDeviceClusterSitePayload(d *schema.ResourceData) map[string]interface{} {
	domainType := d.Get("domain_type").(string)
	vmmDomainType := d.Get("vmm_domain_type").(string)
	domainName := d.Get("domain_name").(string)
	domainDN := d.Get("domain_dn").(string)

	family := resolveDomainFamily(domainType, domainDN)
	resolvedDomainDN := domainDN
	if resolvedDomainDN == "" {
		resolvedDomainDN = composeDomainDn(domainType, vmmDomainType, domainName)
	}

	return map[string]interface{}{
		"name":                 d.Get("name").(string),
		"isPhysicalDomain":     familyToIsPhysicalDomain[family],
		"domainDn":             resolvedDomainDN,
		"highAvailabilityMode": d.Get("high_availability_mode").(string),
		"promiscuousMode":      d.Get("promiscuous_mode").(bool),
		"trunkingPort":         d.Get("trunking_port").(bool),
		"interfaces":           buildServiceDeviceClusterSiteInterfacesPayload(d),
	}
}

func buildServiceDeviceClusterSiteInterfacesPayload(d *schema.ResourceData) []map[string]interface{} {
	interfaces, ok := d.Get("interfaces").([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
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
		if fabricList, ok := interfaceData["fabric_to_device_connectivity"].([]interface{}); ok && len(fabricList) > 0 {
			entry["fabricToDeviceConnectivity"] = buildFabricToDeviceConnectivityPayload(fabricList, vlan)
		}
		if vmSet, ok := interfaceData["vm_information"].(*schema.Set); ok && vmSet.Len() > 0 {
			entry["vmmIntfInfo"] = buildVMMIntfInfoPayload(vmSet)
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

func buildFabricToDeviceConnectivityPayload(fabricPaths []interface{}, vlan int) []map[string]interface{} {
	payload := make([]map[string]interface{}, 0, len(fabricPaths))
	for _, rawFabricPath := range fabricPaths {
		fabricPathData := rawFabricPath.(map[string]interface{})
		podID := fabricPathData["pod_id"].(string)
		nodeIDs := nodeIDsFromAttr(fabricPathData["node_id"])
		path := fabricPathData["path"].(string)
		portType := fabricPathData["port_type"].(string)
		interfaceDn := composeFabricInterfaceDn(podID, nodeIDs, path, portType)
		payload = append(payload, map[string]interface{}{
			"podID":       podID,
			"nodeID":      strings.Join(nodeIDs, ","),
			"path":        interfaceDn,
			"portType":    portType,
			"interfaceDn": interfaceDn,
		})
	}
	return payload
}

func buildVMMIntfInfoPayload(vmSet *schema.Set) []map[string]interface{} {
	vms := vmSet.List()
	payload := make([]map[string]interface{}, 0, len(vms))
	for _, rawVM := range vms {
		vmData := rawVM.(map[string]interface{})
		payload = append(payload, map[string]interface{}{
			"vmName":   vmData["vm_name"].(string),
			"vNicName": vmData["vnic_name"].(string),
		})
	}
	return payload
}

func buildPbrDestinationsPayload(destinations []interface{}) []map[string]interface{} {
	payload := make([]map[string]interface{}, 0, len(destinations))
	for _, rawDestination := range destinations {
		destinationData := rawDestination.(map[string]interface{})
		entry := map[string]interface{}{
			"ip": destinationData["ip"].(string),
			// TODO: derive isAdvancedConfigSet from advanced fields (mac, additionalTrackingIP, weight, isBackUp, tag).
			//       For now always false until the NDO UI / API behaviour for this flag is confirmed.
			"isAdvancedConfigSet": false,
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
// owned field. When forceAll is true every field is emitted (used by Create when an
// existing device entry is found and must be fully overwritten). When forceAll is false
// only fields with d.HasChange are emitted (used by Update).
func appendServiceDeviceClusterSiteReplacePatches(payloadCont *container.Container, updatePath string, d *schema.ResourceData, forceAll bool) {
	domainGroupChanged := forceAll || d.HasChange("domain_type") || d.HasChange("vmm_domain_type") || d.HasChange("domain_name") || d.HasChange("domain_dn")
	if domainGroupChanged {
		domainType := d.Get("domain_type").(string)
		vmmDomainType := d.Get("vmm_domain_type").(string)
		domainName := d.Get("domain_name").(string)
		domainDN := d.Get("domain_dn").(string)

		family := resolveDomainFamily(domainType, domainDN)
		resolvedDomainDN := domainDN
		if resolvedDomainDN == "" {
			resolvedDomainDN = composeDomainDn(domainType, vmmDomainType, domainName)
		}
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/isPhysicalDomain", updatePath), familyToIsPhysicalDomain[family])
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/domainDn", updatePath), resolvedDomainDN)
	}
	if forceAll || d.HasChange("high_availability_mode") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/highAvailabilityMode", updatePath), d.Get("high_availability_mode").(string))
	}
	if forceAll || d.HasChange("promiscuous_mode") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/promiscuousMode", updatePath), d.Get("promiscuous_mode").(bool))
	}
	if forceAll || d.HasChange("trunking_port") {
		addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/trunkingPort", updatePath), d.Get("trunking_port").(bool))
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
