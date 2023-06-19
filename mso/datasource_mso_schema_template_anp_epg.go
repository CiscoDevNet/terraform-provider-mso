package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaTemplateAnpEpg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateAnpEpgRead,

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
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bd_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bd_template_name": &schema.Schema{
				Type:     schema.TypeString,
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
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"useg_epg": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"intra_epg": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"intersite_multicast_source": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"proxy_arp": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"preferred_group": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"epg_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateAnpEpgRead(d *schema.ResourceData, m interface{}) error {
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
	stateANP := d.Get("anp_name")
	stateEPG := d.Get("name")
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				apiANP := models.StripQuotes(anpCont.S("name").String())
				if apiANP == stateANP {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						apiEPG := models.StripQuotes(epgCont.S("name").String())
						if apiEPG == stateEPG {
							d.SetId(apiEPG)
							d.Set("schema_id", schemaId)
							d.Set("name", apiEPG)
							d.Set("template_name", apiTemplate)
							d.Set("display_name", models.StripQuotes(epgCont.S("displayName").String()))
							d.Set("description", models.StripQuotes(epgCont.S("description").String()))
							d.Set("intra_epg", models.StripQuotes(epgCont.S("intraEpg").String()))
							d.Set("useg_epg", epgCont.S("uSegEpg").Data().(bool))
							if epgCont.Exists("mCastSource") {
								d.Set("intersite_multicast_source", epgCont.S("mCastSource").Data().(bool))
							}
							if epgCont.Exists("proxyArp") {
								d.Set("proxy_arp", epgCont.S("proxyArp").Data().(bool))
							}
							d.Set("preferred_group", epgCont.S("preferredGroup").Data().(bool))

							vrfRef := models.StripQuotes(epgCont.S("vrfRef").String())
							re_vrf := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
							match_vrf := re_vrf.FindStringSubmatch(vrfRef)
							d.Set("vrf_name", match_vrf[3])
							d.Set("vrf_schema_id", match_vrf[1])
							d.Set("vrf_template_name", match_vrf[2])

							bdRef := models.StripQuotes(epgCont.S("bdRef").String())
							re_bd := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
							match_bd := re_bd.FindStringSubmatch(bdRef)
							d.Set("bd_name", match_bd[3])
							d.Set("bd_schema_id", match_bd[1])
							d.Set("bd_template_name", match_bd[2])

							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the EPG %s", stateEPG)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
