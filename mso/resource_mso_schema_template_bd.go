package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTemplateBD() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateBDCreate,
		Read:   resourceMSOTemplateBDRead,
		Update: resourceMSOTemplateBDUpdate,
		Delete: resourceMSOTemplateBDDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateBDImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"intersite_bum_traffic": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"optimize_wan_bandwidth": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"layer2_stretch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"layer3_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"layer2_unknown_unicast": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"unknown_multicast_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"flood",
					"optimized_flooding",
				}, false),
			},
			"multi_destination_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"flood_in_bd",
					"drop",
					"flood_in_encap",
				}, false),
			},
			"ipv6_unknown_multicast_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"flood",
					"optimized_flooding",
				}, false),
			},
			"arp_flooding": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"virtual_mac_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"unicast_routing": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_policy": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				Description:   "Configure dhcp policy in versions before NDO 3.2",
				Computed:      true,
				ConflictsWith: []string{"dhcp_policies"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dhcp_policies": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "Configure dhcp policies in versions NDO 3.2 and higher",
				ConflictsWith: []string{"dhcp_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ep_move_detection_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"garp",
				}, false),
			},
		}),
	}
}

func resourceMSOTemplateBDImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := get_attribute[2]

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return nil, err
	}

	found := false
	stateBD := get_attribute[4]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())
		if apiTemplate == stateTemplate {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return nil, fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return nil, err
				}
				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {
					d.SetId(apiBD)
					d.Set("name", apiBD)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(bdCont.S("displayName").String()))
					d.Set("description", models.StripQuotes(bdCont.S("description").String()))
					d.Set("layer2_unknown_unicast", models.StripQuotes(bdCont.S("l2UnknownUnicast").String()))
					if models.StripQuotes(bdCont.S("unkMcastAct").String()) == "opt-flood" {
						d.Set("unknown_multicast_flooding", "optimized_flooding")
					} else {
						d.Set("unknown_multicast_flooding", "flood")
					}
					multiDstPktAct := models.StripQuotes(bdCont.S("multiDstPktAct").String())
					if multiDstPktAct == "encap-flood" {
						d.Set("multi_destination_flooding", "flood_in_encap")
					} else if multiDstPktAct == "bd-flood" {
						d.Set("multi_destination_flooding", "flood_in_bd")
					} else {
						d.Set("multi_destination_flooding", "drop")
					}
					v6unkMcastAct := models.StripQuotes(bdCont.S("v6unkMcastAct").String())
					if v6unkMcastAct == "opt-flood" {
						d.Set("ipv6_unknown_multicast_flooding", "optimized_flooding")
					} else {
						d.Set("ipv6_unknown_multicast_flooding", "flood")
					}
					vmac := models.StripQuotes(bdCont.S("vmac").String())
					if vmac != "{}" {
						d.Set("virtual_mac_address", vmac)
					} else {
						d.Set("virtual_mac_address", "")
					}

					epMoveDetectMode := models.StripQuotes(bdCont.S("epMoveDetectMode").String())
					if epMoveDetectMode != "{}" {
						d.Set("ep_move_detection_mode", epMoveDetectMode)
					} else {
						d.Set("ep_move_detection_mode", "none") // set to default value of none when not present
					}

					if bdCont.Exists("intersiteBumTrafficAllow") {
						d.Set("intersite_bum_traffic", bdCont.S("intersiteBumTrafficAllow").Data().(bool))
					}

					if bdCont.Exists("optimizeWanBandwidth") {
						d.Set("optimize_wan_bandwidth", bdCont.S("optimizeWanBandwidth").Data().(bool))
					}

					if bdCont.Exists("l3MCast") {
						d.Set("layer3_multicast", bdCont.S("l3MCast").Data().(bool))
					}

					if bdCont.Exists("l2Stretch") {
						d.Set("layer2_stretch", bdCont.S("l2Stretch").Data().(bool))
					}
					if bdCont.Exists("arpFlood") {
						d.Set("arp_flooding", bdCont.S("arpFlood").Data().(bool))
					}
					if bdCont.Exists("unicastRouting") {
						d.Set("unicast_routing", bdCont.S("unicastRouting").Data().(bool))
					}

					vrfRef := models.StripQuotes(bdCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

					dhcpPolMap := make(map[string]interface{})
					dhcpPoliciesList := make([]interface{}, 0)

					var dhcpCount int
					if bdCont.Exists("dhcpLabels") {
						dhcpCount, err = bdCont.ArrayCount("dhcpLabels")
						if err != nil {
							return nil, err
						}
					} else {
						dhcpCount = 0
					}
					if dhcpCount != 0 {
						if versionInt == -1 {
							dhcpPoliciesNameList, error := getDHCPPolicesNameByRef(dhcpCount, schemaId, stateTemplate, bdCont, msoClient)
							if error != nil {
								return nil, error
							}
							dhcpPoliciesList = dhcpPoliciesNameList
						} else {
							for l := 0; l < dhcpCount; l++ {
								dhcpPolicy, err := bdCont.ArrayElement(l, "dhcpLabels")
								if err != nil {
									return nil, err
								}
								dhcpPolicyMap := make(map[string]interface{})
								dhcpPolicyMap["name"] = models.StripQuotes(dhcpPolicy.S("name").String())
								version, err := strconv.Atoi(models.StripQuotes(dhcpPolicy.S("version").String()))
								if err != nil {
									return nil, err
								}
								dhcpPolicyMap["version"] = version
								if dhcpPolicy.Exists("dhcpOptionLabel") {
									dhcpPolicyMap["dhcp_option_policy_name"] = models.StripQuotes(dhcpPolicy.S("dhcpOptionLabel", "name").String())
									version, err := strconv.Atoi(models.StripQuotes(dhcpPolicy.S("dhcpOptionLabel", "version").String()))
									if err != nil {
										return nil, err
									}
									dhcpPolicyMap["dhcp_option_policy_version"] = version
									if dhcpPolicyMap["dhcp_option_policy_name"] == "{}" {
										dhcpPolicyMap["dhcp_option_policy_name"] = nil
									}
									if dhcpPolicyMap["dhcp_option_policy_version"] == "{}" {
										dhcpPolicyMap["dhcp_option_policy_version"] = nil
									}
								}
								dhcpPoliciesList = append(dhcpPoliciesList, dhcpPolicyMap)
							}
						}
					} else {
						if bdCont.Exists("dhcpLabel") {
							dhcpPolMap["name"] = models.StripQuotes(bdCont.S("dhcpLabel", "name").String())
							dhcpPolMap["version"] = models.StripQuotes(bdCont.S("dhcpLabel", "version").String())
							if dhcpPolMap["version"] == "{}" {
								dhcpPolMap["version"] = nil
							}
							if bdCont.Exists("dhcpLabel", "dhcpOptionLabel") {
								dhcpPolMap["dhcp_option_policy_name"] = models.StripQuotes(bdCont.S("dhcpLabel", "dhcpOptionLabel", "name").String())
								dhcpPolMap["dhcp_option_policy_version"] = models.StripQuotes(bdCont.S("dhcpLabel", "dhcpOptionLabel", "version").String())
								if dhcpPolMap["dhcp_option_policy_name"] == "{}" {
									dhcpPolMap["dhcp_option_policy_name"] = nil
								}
								if dhcpPolMap["dhcp_option_policy_version"] == "{}" {
									dhcpPolMap["dhcp_option_policy_version"] = nil
								}
							}
						}
					}
					d.Set("dhcp_policy", dhcpPolMap)
					d.Set("dhcp_policies", dhcpPoliciesList)
					found = true
					break
				}
			}
		}
	}
	if !found {
		d.SetId("")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOTemplateBDCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	var intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast, unicast_routing, arp_flooding bool
	var layer2_unknown_unicast, vrf_schema_id, vrf_template_name, virtual_mac_address, ep_move_detection_mode, ipv6_unknown_multicast_flooding, multi_destination_flooding, unknown_multicast_flooding string

	if tempVar, ok := d.GetOk("intersite_bum_traffic"); ok {
		intersite_bum_traffic = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("optimize_wan_bandwidth"); ok {
		optimize_wan_bandwidth = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_stretch"); ok {
		layer2_stretch = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer3_multicast"); ok {
		layer3_multicast = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_unknown_unicast"); ok {
		layer2_unknown_unicast = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("unknown_multicast_flooding"); ok {
		unknown_multicast_flooding = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("multi_destination_flooding"); ok {
		multi_destination_flooding = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("ipv6_unknown_multicast_flooding"); ok {
		ipv6_unknown_multicast_flooding = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("unicast_routing"); ok {
		unicast_routing = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("virtual_mac_address"); ok {
		virtual_mac_address = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("ep_move_detection_mode"); ok {
		ep_move_detection_mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("arp_flooding"); ok {
		arp_flooding = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	var dhcpPolMap map[string]interface{}
	if tempVar, ok := d.GetOk("dhcp_policy"); ok {
		dhcp_policy := tempVar.(map[string]interface{})
		dhcpPolMap = make(map[string]interface{})
		dhcpPolMap["name"] = dhcp_policy["name"]
		version, err := strconv.Atoi(dhcp_policy["version"].(string))
		if err != nil {
			return err
		}
		dhcpPolMap["version"] = version

		optionName := dhcp_policy["dhcp_option_policy_name"]
		optionVersion := dhcp_policy["dhcp_option_policy_version"]
		if optionName != nil {
			if optionVersion != nil {
				dhcpOptionMap := make(map[string]interface{})
				dhcpOptionMap["name"] = optionName
				ver, err := strconv.Atoi(optionVersion.(string))
				if err != nil {
					return err
				}
				dhcpOptionMap["version"] = ver
				dhcpPolMap["dhcpOptionLabel"] = dhcpOptionMap
			} else {
				return fmt.Errorf("dhcp_option_policy_version is required with dhcp_option_policy_name")
			}
		}
	} else {
		dhcpPolMap = nil
	}

	dhcpPolList := make([]interface{}, 0)
	if dhcpPolicies, ok := d.GetOk("dhcp_policies"); ok {
		if versionInt == -1 {
			dhcpRefList, err := mapDHCPPoliciesRefByName(schemaID, templateName, dhcpPolicies, msoClient)
			if err != nil {
				return err
			}
			dhcpPolList = dhcpRefList
		} else {
			for _, dhcpPolicy := range dhcpPolicies.(*schema.Set).List() {
				policy := dhcpPolicy.(map[string]interface{})
				dhcpPolicyMap := make(map[string]interface{})
				dhcpPolicyMap["name"] = policy["name"]
				dhcpPolicyMap["version"] = policy["version"]
				if policy["dhcp_option_policy_name"] != "" {
					dhcpOptionMap := make(map[string]interface{})
					dhcpOptionMap["name"] = policy["dhcp_option_policy_name"]
					if policy["version"] != 0 {
						dhcpOptionMap["version"] = policy["dhcp_option_policy_version"]
					} else {
						dhcpOptionMap["version"] = policy["version"]
					}
					dhcpPolicyMap["dhcpOptionLabel"] = dhcpOptionMap
				}
				dhcpPolList = append(dhcpPolList, dhcpPolicyMap)
			}
		}
	} else {
		dhcpPolList = nil
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName
	path := fmt.Sprintf("/templates/%s/bds/-", templateName)
	bdStruct := models.NewTemplateBD("add", path, name, displayName, layer2_unknown_unicast, unknown_multicast_flooding, multi_destination_flooding, ipv6_unknown_multicast_flooding, virtual_mac_address, ep_move_detection_mode, description, intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast, arp_flooding, unicast_routing, vrfRefMap, dhcpPolMap, dhcpPolList)
	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), bdStruct)

	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Create finished successfully", d.Id())
	return resourceMSOTemplateBDRead(d, m)
}

func resourceMSOTemplateBDRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateBD := d.Get("name")

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			bdCount, err := tempCont.ArrayCount("bds")
			if err != nil {
				return fmt.Errorf("Unable to get BD list")
			}
			for j := 0; j < bdCount; j++ {
				bdCont, err := tempCont.ArrayElement(j, "bds")
				if err != nil {
					return err
				}

				apiBD := models.StripQuotes(bdCont.S("name").String())
				if apiBD == stateBD {
					d.SetId(apiBD)
					d.Set("name", apiBD)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(bdCont.S("displayName").String()))
					d.Set("description", models.StripQuotes(bdCont.S("description").String()))
					d.Set("layer2_unknown_unicast", models.StripQuotes(bdCont.S("l2UnknownUnicast").String()))
					if models.StripQuotes(bdCont.S("unkMcastAct").String()) == "opt-flood" {
						d.Set("unknown_multicast_flooding", "optimized_flooding")
					} else {
						d.Set("unknown_multicast_flooding", "flood")
					}
					multiDstPktAct := models.StripQuotes(bdCont.S("multiDstPktAct").String())
					if multiDstPktAct == "encap-flood" {
						d.Set("multi_destination_flooding", "flood_in_encap")
					} else if multiDstPktAct == "bd-flood" {
						d.Set("multi_destination_flooding", "flood_in_bd")
					} else {
						d.Set("multi_destination_flooding", "drop")
					}
					v6unkMcastAct := models.StripQuotes(bdCont.S("v6unkMcastAct").String())
					if v6unkMcastAct == "opt-flood" {
						d.Set("ipv6_unknown_multicast_flooding", "optimized_flooding")
					} else {
						d.Set("ipv6_unknown_multicast_flooding", "flood")
					}

					vmac := models.StripQuotes(bdCont.S("vmac").String())
					if vmac != "{}" {
						d.Set("virtual_mac_address", vmac)
					} else {
						d.Set("virtual_mac_address", "")
					}

					epMoveDetectMode := models.StripQuotes(bdCont.S("epMoveDetectMode").String())
					if epMoveDetectMode != "{}" {
						d.Set("ep_move_detection_mode", epMoveDetectMode)
					} else {
						d.Set("ep_move_detection_mode", "none") // set to default value of none when not present
					}

					if bdCont.Exists("intersiteBumTrafficAllow") {
						d.Set("intersite_bum_traffic", bdCont.S("intersiteBumTrafficAllow").Data().(bool))
					}

					if bdCont.Exists("optimizeWanBandwidth") {
						d.Set("optimize_wan_bandwidth", bdCont.S("optimizeWanBandwidth").Data().(bool))
					}

					if bdCont.Exists("l3MCast") {
						d.Set("layer3_multicast", bdCont.S("l3MCast").Data().(bool))
					}

					if bdCont.Exists("l2Stretch") {
						d.Set("layer2_stretch", bdCont.S("l2Stretch").Data().(bool))
					}
					if bdCont.Exists("arpFlood") {
						d.Set("arp_flooding", bdCont.S("arpFlood").Data().(bool))
					}
					if bdCont.Exists("unicastRouting") {
						d.Set("unicast_routing", bdCont.S("unicastRouting").Data().(bool))
					}

					vrfRef := models.StripQuotes(bdCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

					dhcpPolMap := make(map[string]interface{})
					dhcpPoliciesList := make([]interface{}, 0)

					var dhcpCount int
					if bdCont.Exists("dhcpLabels") {
						dhcpCount, err = bdCont.ArrayCount("dhcpLabels")
						if err != nil {
							return err
						}
					} else {
						dhcpCount = 0
					}

					if dhcpCount != 0 {
						if versionInt == -1 {
							dhcpPoliciesNameList, error := getDHCPPolicesNameByRef(dhcpCount, schemaId, stateTemplate, bdCont, msoClient)
							if error != nil {
								return error
							}
							dhcpPoliciesList = dhcpPoliciesNameList
						} else {
							for l := 0; l < dhcpCount; l++ {
								dhcpPolicy, err := bdCont.ArrayElement(l, "dhcpLabels")
								if err != nil {
									return err
								}
								dhcpPolicyMap := make(map[string]interface{})
								dhcpPolicyMap["name"] = models.StripQuotes(dhcpPolicy.S("name").String())
								var version int
								if dhcpPolicy.Exists("version") {
									version, err = strconv.Atoi(models.StripQuotes(dhcpPolicy.S("version").String()))
								}
								if err != nil {
									return err
								}
								dhcpPolicyMap["version"] = version
								if dhcpPolicy.Exists("dhcpOptionLabel") {
									dhcpPolicyMap["dhcp_option_policy_name"] = models.StripQuotes(dhcpPolicy.S("dhcpOptionLabel", "name").String())
									version, err := strconv.Atoi(models.StripQuotes(dhcpPolicy.S("dhcpOptionLabel", "version").String()))
									if err != nil {
										return err
									}
									dhcpPolicyMap["dhcp_option_policy_version"] = version
								}
								dhcpPoliciesList = append(dhcpPoliciesList, dhcpPolicyMap)
							}
						}
					} else {
						if bdCont.Exists("dhcpLabel") {
							dhcpPolMap["name"] = models.StripQuotes(bdCont.S("dhcpLabel", "name").String())
							dhcpPolMap["version"] = models.StripQuotes(bdCont.S("dhcpLabel", "version").String())
							if dhcpPolMap["version"] == "{}" {
								dhcpPolMap["version"] = nil
							}
							if bdCont.Exists("dhcpLabel", "dhcpOptionLabel") {
								dhcpPolMap["dhcp_option_policy_name"] = models.StripQuotes(bdCont.S("dhcpLabel", "dhcpOptionLabel", "name").String())
								dhcpPolMap["dhcp_option_policy_version"] = models.StripQuotes(bdCont.S("dhcpLabel", "dhcpOptionLabel", "version").String())
								if dhcpPolMap["dhcp_option_policy_name"] == "" {
									dhcpPolMap["dhcp_option_policy_name"] = nil
								}
								if dhcpPolMap["dhcp_option_policy_version"] == "{}" {
									dhcpPolMap["dhcp_option_policy_version"] = nil
								}
							}
						}
					}
					d.Set("dhcp_policy", dhcpPolMap)
					d.Set("dhcp_policies", dhcpPoliciesList)
					found = true
					break
				}

			}
		}

	}

	if !found {
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOTemplateBDUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}
	var intersite_bum_traffic, optimize_wan_bandwidth, layer2_stretch, layer3_multicast, unicast_routing, arp_flooding bool
	var layer2_unknown_unicast, vrf_schema_id, vrf_template_name, virtual_mac_address, ep_move_detection_mode, ipv6_unknown_multicast_flooding, multi_destination_flooding, unknown_multicast_flooding string

	if tempVar, ok := d.GetOk("intersite_bum_traffic"); ok {
		intersite_bum_traffic = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("optimize_wan_bandwidth"); ok {
		optimize_wan_bandwidth = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_stretch"); ok {
		layer2_stretch = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer3_multicast"); ok {
		layer3_multicast = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("layer2_unknown_unicast"); ok {
		layer2_unknown_unicast = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("unknown_multicast_flooding"); ok {
		unknown_multicast_flooding = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("multi_destination_flooding"); ok {
		multi_destination_flooding = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("ipv6_unknown_multicast_flooding"); ok {
		ipv6_unknown_multicast_flooding = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("unicast_routing"); ok {
		unicast_routing = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("virtual_mac_address"); ok {
		virtual_mac_address = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("ep_move_detection_mode"); ok {
		ep_move_detection_mode = tempVar.(string)
	}
	if tempVar, ok := d.GetOk("arp_flooding"); ok {
		arp_flooding = tempVar.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	var dhcpPolMap map[string]interface{}
	if tempVar, ok := d.GetOk("dhcp_policy"); ok {
		dhcp_policy := tempVar.(map[string]interface{})
		dhcpPolMap = make(map[string]interface{})
		dhcpPolMap["name"] = dhcp_policy["name"]
		version, err := strconv.Atoi(dhcp_policy["version"].(string))
		if err != nil {
			return err
		}
		dhcpPolMap["version"] = version

		optionName := dhcp_policy["dhcp_option_policy_name"]
		optionVersion := dhcp_policy["dhcp_option_policy_version"]
		if optionName != nil {
			if optionVersion != nil {
				dhcpOptionMap := make(map[string]interface{})
				dhcpOptionMap["name"] = optionName
				ver, err := strconv.Atoi(optionVersion.(string))
				if err != nil {
					return err
				}
				dhcpOptionMap["version"] = ver
				dhcpPolMap["dhcpOptionLabel"] = dhcpOptionMap
			} else {
				return fmt.Errorf("dhcp_option_policy_version is required with dhcp_option_policy_name")
			}
		}
	} else {
		dhcpPolMap = nil
	}

	dhcpPolList := make([]interface{}, 0)
	if dhcpPolicies, ok := d.GetOk("dhcp_policies"); ok {
		if versionInt == -1 {
			dhcpRefList, err := mapDHCPPoliciesRefByName(schemaID, templateName, dhcpPolicies, msoClient)
			if err != nil {
				return err
			}
			dhcpPolList = dhcpRefList
		} else {
			for _, dhcpPolicy := range dhcpPolicies.(*schema.Set).List() {
				policy := dhcpPolicy.(map[string]interface{})
				dhcpPolicyMap := make(map[string]interface{})
				dhcpPolicyMap["name"] = policy["name"]
				dhcpPolicyMap["version"] = policy["version"]
				if policy["dhcp_option_policy_name"] != "" {
					dhcpOptionMap := make(map[string]interface{})
					dhcpOptionMap["name"] = policy["dhcp_option_policy_name"]
					if policy["version"] != 0 {
						dhcpOptionMap["version"] = policy["dhcp_option_policy_version"]
					} else {
						dhcpOptionMap["version"] = policy["version"]
					}
					dhcpPolicyMap["dhcpOptionLabel"] = dhcpOptionMap
				}
				dhcpPolList = append(dhcpPolList, dhcpPolicyMap)
			}
		}
	} else {
		dhcpPolList = nil
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	basePath := fmt.Sprintf("/templates/%s/bds/%s", templateName, name)
	payloadCon := container.New()
	payloadCon.Array()

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/displayName", basePath), displayName)
	if err != nil {
		return err
	}

	if layer2_unknown_unicast != "" {
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/l2UnknownUnicast", basePath), layer2_unknown_unicast)
		if err != nil {
			return err
		}
	}

	if unknown_multicast_flooding != "" {
		if unknown_multicast_flooding == "optimized_flooding" {
			unknown_multicast_flooding = "opt-flood"
		}
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/unkMcastAct", basePath), unknown_multicast_flooding)
		if err != nil {
			return err
		}
	}

	if multi_destination_flooding != "" {
		if multi_destination_flooding == "flood_in_encap" {
			multi_destination_flooding = "encap-flood"
		} else if multi_destination_flooding == "flood_in_bd" {
			multi_destination_flooding = "bd-flood"
		}
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/multiDstPktAct", basePath), multi_destination_flooding)
		if err != nil {
			return err
		}
	}

	if ipv6_unknown_multicast_flooding != "" {
		if ipv6_unknown_multicast_flooding == "optimized_flooding" {
			ipv6_unknown_multicast_flooding = "opt-flood"
		}
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/v6unkMcastAct", basePath), ipv6_unknown_multicast_flooding)
		if err != nil {
			return err
		}
	}

	if virtual_mac_address != "" {
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/vmac", basePath), virtual_mac_address)
		if err != nil {
			return err
		}
	}

	if ep_move_detection_mode != "" {
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/epMoveDetectMode", basePath), ep_move_detection_mode)
		if err != nil {
			return err
		}
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/description", basePath), description)
	if err != nil {
		return err
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/intersiteBumTrafficAllow", basePath), intersite_bum_traffic)
	if err != nil {
		return err
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/optimizeWanBandwidth", basePath), optimize_wan_bandwidth)
	if err != nil {
		return err
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/l2Stretch", basePath), layer2_stretch)
	if err != nil {
		return err
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/l3MCast", basePath), layer3_multicast)
	if err != nil {
		return err
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/arpFlood", basePath), arp_flooding)
	if err != nil {
		return err
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/unicastRouting", basePath), unicast_routing)
	if err != nil {
		return err
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/vrfRef", basePath), vrfRefMap)
	if err != nil {
		return err
	}

	if len(dhcpPolMap) != 0 {
		err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/dhcpLabel", basePath), dhcpPolMap)
		if err != nil {
			return err
		}
	}

	err = addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/dhcpLabels", basePath), dhcpPolList)
	if err != nil {
		return err
	}

	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/schemas/%s", schemaID), payloadCon)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOTemplateBDRead(d, m)
}

func resourceMSOTemplateBDDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Delete")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	name := d.Get("name").(string)
	templateName := d.Get("template_name").(string)
	patchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/bds/%s", templateName, name))
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), patchPayload)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: Delete finished successfully", d.Id())
	return nil
}

// getDHCPPolicesNameByRef retrieves the DHCP policies name by reference.
//
// Parameters:
// - dhcpCount: The number of DHCP policies to retrieve.
// - schemaId: The ID of the schema template.
// - stateTemplate: The state template.
// - bdCont: The container.
// - msoClient: The client.
// Returns:
// - []interface{}: The list of DHCP policies name.
// - error: An error if the retrieval fails.
func getDHCPPolicesNameByRef(dhcpCount int, schemaId, stateTemplate string, bdCont *container.Container, msoClient *client.Client) ([]interface{}, error) {
	tenantID, err := msoClient.GetTenantIDFromSchemaTemplate(schemaId, stateTemplate)
	if err != nil {
		return nil, err
	}
	dhcpPoliciesRefList := make([]interface{}, 0)
	for l := 0; l < dhcpCount; l++ {
		dhcpPolicy, err := bdCont.ArrayElement(l, "dhcpLabels")
		if err != nil {
			return nil, err
		}
		dhcpPolicyMap := make(map[string]interface{})
		dhcpPolicyMap["relayRef"] = models.StripQuotes(dhcpPolicy.S("ref").String())
		dhcpPolicyMap["optionRef"] = models.StripQuotes(dhcpPolicy.S("dhcpOptionLabel", "ref").String())
		dhcpPoliciesRefList = append(dhcpPoliciesRefList, dhcpPolicyMap)
	}
	return msoClient.GetDHCPPoliciesNameByUUID(tenantID, dhcpPoliciesRefList)
}

// mapDHCPPoliciesRefByName retrieves a list of DHCP policy UUIDs by name.
//
// Parameters:
// - schemaID: the ID of the schema.
// - templateName: the name of the template.
// - dhcpPolicies: the list of DHCP policies.
// - msoClient: the client used to communicate with the MSO.
// Returns:
// - []interface{}: It returns a list of interface{} containing the UUIDs of the DHCP policies,
// - error: An error if the retrieval fails.
func mapDHCPPoliciesRefByName(schemaID, templateName string, dhcpPolicies interface{}, msoClient *client.Client) ([]interface{}, error) {
	tenantID, err := msoClient.GetTenantIDFromSchemaTemplate(schemaID, templateName)
	if err != nil {
		return nil, err
	}
	dhcpPolicyNameList := make([]interface{}, 0)
	for _, dhcpPolicy := range dhcpPolicies.(*schema.Set).List() {
		policy := dhcpPolicy.(map[string]interface{})
		dhcpPolicyNameMap := make(map[string]interface{})
		dhcpPolicyNameMap["relayName"] = policy["name"]
		dhcpPolicyNameMap["optionName"] = policy["dhcp_option_policy_name"]
		dhcpPolicyNameList = append(dhcpPolicyNameList, dhcpPolicyNameMap)
	}
	return msoClient.GetDHCPPoliciesUUIDByName(tenantID, dhcpPolicyNameList)
}
