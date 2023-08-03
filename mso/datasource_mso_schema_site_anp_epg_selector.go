package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteAnpEpgSelector() *schema.Resource {
	return &schema.Resource{
		Read: datasourceSchemaSiteApnEpgSelectorRead,

		Schema: map[string]*schema.Schema{
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
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"expressions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceSchemaSiteApnEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	name := d.Get("name").(string)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("schema_id", schemaId)
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	anpCont, err := getSiteAnp(anp, siteCont)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", anp)
	}

	epgCont, err := getSiteEpg(epg, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", epg)
	}

	selectorCount, err := epgCont.ArrayCount("selectors")
	if err != nil {
		return fmt.Errorf("No selectors found")
	}

	found := false
	for s := 0; s < selectorCount; s++ {
		selectorCont, err := epgCont.ArrayElement(s, "selectors")
		if err != nil {
			return err
		}

		currentName := models.StripQuotes(selectorCont.S("name").String())
		if currentName == name {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/anps/%s/epgs/%s/selectors/%s", schemaId, siteId, templateName, anp, epg, name))
			d.Set("name", name)
			exps := selectorCont.S("expressions").Data().([]interface{})

			expressionsList := make([]interface{}, 0, 1)
			for _, val := range exps {
				tp := val.(map[string]interface{})
				expressionsMap := make(map[string]interface{})

				expressionsMap["key"] = tp["key"]

				expressionsMap["operator"] = tp["operator"]

				if tp["value"] != nil {
					expressionsMap["value"] = tp["value"]
				}
				expressionsList = append(expressionsList, expressionsMap)
			}
			d.Set("expressions", expressionsList)
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the Site Anp Epg Selector %s", name)
	}
	return nil
}
