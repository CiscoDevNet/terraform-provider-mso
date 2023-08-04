package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOTemplateExternalepg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateExternalepgRead,

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
			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"external_epg_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
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
			"anp_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"anp_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"anp_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l3out_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l3out_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l3out_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"selector_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"selector_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOTemplateExternalepgRead(d *schema.ResourceData, m interface{}) error {
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
	stateExternalepg := d.Get("external_epg_name")

	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			externalepgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get Externalepg list")
			}
			for j := 0; j < externalepgCount; j++ {
				externalepgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				apiExternalepg := models.StripQuotes(externalepgCont.S("name").String())
				if apiExternalepg == stateExternalepg {
					d.SetId(fmt.Sprintf("%s/templates/%s/externalEpgs/%s", schemaId, stateTemplate, stateExternalepg))
					d.Set("external_epg_name", apiExternalepg)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(externalepgCont.S("displayName").String()))
					d.Set("external_epg_type", models.StripQuotes(externalepgCont.S("extEpgType").String()))

					vrfRef := models.StripQuotes(externalepgCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])

					anpRef := models.StripQuotes(externalepgCont.S("anpRef").String())
					if anpRef != "{}" && anpRef != "" {
						tokens := strings.Split(anpRef, "/")
						d.Set("anp_name", tokens[len(tokens)-1])
						d.Set("anp_schema_id", tokens[len(tokens)-5])
						d.Set("anp_template_name", tokens[len(tokens)-3])
					} else {
						d.Set("anp_name", "")
						d.Set("anp_schema_id", "")
						d.Set("anp_template_name", "")
					}

					l3outRef := models.StripQuotes(externalepgCont.S("l3outRef").String())
					if l3outRef != "{}" && l3outRef != "" {
						tokens := strings.Split(l3outRef, "/")
						d.Set("l3out_name", tokens[len(tokens)-1])
						d.Set("l3out_schema_id", tokens[len(tokens)-5])
						d.Set("l3out_template_name", tokens[len(tokens)-3])
					} else {
						d.Set("l3out_name", "")
						d.Set("l3out_schema_id", "")
						d.Set("l3out_template_name", "")
					}

					epgType := d.Get("external_epg_type").(string)
					if epgType == "cloud" {
						selList := externalepgCont.S("selectors").Data().([]interface{})
						if len(selList) > 0 {
							selector := selList[0].(map[string]interface{})
							d.Set("selector_name", selector["name"])
							expList := selector["expressions"].([]interface{})
							if len(expList) > 0 {
								exp := expList[0].(map[string]interface{})
								d.Set("selector_ip", exp["value"])
							} else {
								d.Set("selector_ip", "")
							}
						} else {
							d.Set("selector_name", "")
							d.Set("selector_ip", "")
						}
					} else {
						d.Set("selector_name", "")
						d.Set("selector_ip", "")
					}

					found = true
					break
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the External Epg %s in Template %s of Schema Id %s", stateExternalepg, stateTemplate, schemaId)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
