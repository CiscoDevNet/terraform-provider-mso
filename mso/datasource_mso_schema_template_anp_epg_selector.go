package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaTemplateAnpEpgSelector() *schema.Resource {
	return &schema.Resource{
		Read: datasourceMSOSchemaTemplateAnpEpgSelectorRead,

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

func datasourceMSOSchemaTemplateAnpEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)

	selectorName := d.Get("name").(string)
	schemaId := d.Get("schema_id").(string)
	template := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	tempCount, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	found := false
	for i := 0; i < tempCount && !found; i++ {
		tempcont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		tempName := models.StripQuotes(tempcont.S("name").String())
		if tempName == template {
			anpCount, err := tempcont.ArrayCount("anps")
			if err != nil {
				return err
			}
			for j := 0; j < anpCount && !found; j++ {
				anpCont, err := tempcont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				currentanpName := models.StripQuotes(anpCont.S("name").String())
				if currentanpName == anpName {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return err
					}
					for k := 0; k < epgCount && !found; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							selectorCount, err := epgCont.ArrayCount("selectors")
							if err != nil {
								return err
							}
							for l := 0; l < selectorCount; l++ {
								selectorCont, err := epgCont.ArrayElement(l, "selectors")
								if err != nil {
									return err
								}
								currSelectorName := models.StripQuotes(selectorCont.S("name").String())
								if currSelectorName == selectorName {
									found = true
									d.SetId(fmt.Sprintf("%s/templates/%s/anps/%s/epgs/%s/selectors/%s", schemaId, template, anpName, epgName, currSelectorName))
									d.Set("name", currSelectorName)
									d.Set("template_name", tempName)
									d.Set("anp_name", anpName)
									d.Set("epg_name", epgName)
									exps := selectorCont.S("expressions").Data().([]interface{})
									expressionsList := make([]interface{}, 0, 1)
									for _, val := range exps {
										tp := val.(map[string]interface{})
										expressionsMap := make(map[string]interface{})

										if tp["key"] != nil {
											expressionsMap["key"] = tp["key"]
										}

										if tp["operator"] != nil {
											expressionsMap["operator"] = tp["operator"]
										}

										if tp["value"] != nil {
											expressionsMap["value"] = tp["value"]
										}
										expressionsList = append(expressionsList, expressionsMap)
									}
									d.Set("expressions", expressionsList)
									break
								}
							}
						}
					}
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the ANP EPG Selector %s in Template %s of Schema Id %s ", selectorName, template, schemaId)
	}

	return nil
}
