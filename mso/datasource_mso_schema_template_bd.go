package mso

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOTemplateBD() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateBDRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"intersite_bum_traffic": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"optimize_wan_bandwidth": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"layer2_stretch": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"layer3_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"layer2_unknown_unicast": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"unknown_multicast_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"multi_destination_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_unknown_multicast_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"arp_flooding": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"virtual_mac_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"unicast_routing": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dhcp_policy": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Configure dhcp policy in versions before NDO 3.2",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"dhcp_policies": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Configure dhcp policies in versions NDO 3.2 and higher",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"ep_move_detection_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateBDRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	templateName := d.Get("template_name").(string)
	bdName := d.Get("name").(string)

	err = setSchemaTemplateBDAttrs(schemaId, templateName, bdName, cont, d, msoClient)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func setSchemaTemplateBDAttrs(schemaId, templateName, bdName string, cont *container.Container, d *schema.ResourceData, msoClient *client.Client) error {
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())
		if apiTemplate == templateName {
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
				if apiBD == bdName {
					found = true
					d.SetId(fmt.Sprintf("%s/templates/%s/bds/%s", schemaId, templateName, bdName))
					d.Set("name", bdName)
					d.Set("schema_id", schemaId)
					d.Set("template_name", templateName)
					d.Set("display_name", models.StripQuotes(bdCont.S("displayName").String()))
					d.Set("description", models.StripQuotes(bdCont.S("description").String()))
					d.Set("layer2_unknown_unicast", models.StripQuotes(bdCont.S("l2UnknownUnicast").String()))
					d.Set("uuid", models.StripQuotes(bdCont.S("uuid").String()))
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
							dhcpPoliciesNameList, error := getDHCPPolicesNameByRef(dhcpCount, schemaId, templateName, bdCont, msoClient)
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
					break
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the BD %s in Template %s of Schema Id %s ", bdName, templateName, schemaId)
	}

	return nil
}
