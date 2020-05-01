package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceMSOTemplateBD() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateBDRead,

		// Importer: &schema.ResourceImporter{
		//     State: resourceMSOSchemaSiteImport,
		// },

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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
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
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateBD := d.Get("name")
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
					log.Printf("Unique BD %v", bdCont)
					d.Set("name", apiBD)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(bdCont.S("displayName").String()))
					d.Set("layer2_unknown_unicast", models.StripQuotes(bdCont.S("l2UnknownUnicast").String()))
					if bdCont.Exists("intersiteBumTrafficAllow") {
						d.Set("intersite_bum_traffic", bdCont.S("intersiteBumTrafficAllow").Data().(bool))
					}

					if bdCont.Exists("optimize_wan_bandwidth") {
						d.Set("optimize_wan_bandwidth", bdCont.S("optimizeWanBandwidth").Data().(bool))
					}

					if bdCont.Exists("layer3_multicast") {
						d.Set("layer3_multicast", bdCont.S("l3MCast").Data().(bool))
					}

					if bdCont.Exists("layer2_stretch") {
						d.Set("layer2_stretch", bdCont.S("l2Stretch").Data().(bool))
					}

					vrfRef := models.StripQuotes(bdCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

					if bdCont.Exists("dhcpLabel") {
						dhcpPolMap := make(map[string]interface{})
						dhcpPolMap["name"] = models.StripQuotes(bdCont.S("dhcpLabel", "name").String())
						dhcpPolMap["version"] = models.StripQuotes(bdCont.S("dhcpLabel", "version").String())
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
						d.Set("dhcp_policy", dhcpPolMap)
					} else {
						d.Set("dhcp_policy", make(map[string]interface{}))
					}
					found = true
					break
				}

			}
		}

	}

	if !found {
		return fmt.Errorf("Unable to find the BD %s", stateBD)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
