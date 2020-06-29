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

			"template": &schema.Schema{
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},

						"operator": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},

						"value": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
		},
	}
}

func datasourceMSOSchemaTemplateAnpEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	found := false
	msoClient := m.(*client.Client)

	dn := d.Get("name").(string)
	schemaID := d.Get("schema_id").(string)
	template := d.Get("template").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	tempCount, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	for i := 0; i < tempCount; i++ {
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

			for j := 0; j < anpCount; j++ {
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

					for k := 0; k < epgCount; k++ {
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
								if currSelectorName == dn {
									found = true
									d.SetId(dn)
									d.Set("name", currSelectorName)
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
							if found {
								d.Set("epg_name", epgName)
								break
							}
						}
					}
					if found {
						d.Set("anp_name", anpName)
						break
					}
				}
			}
			if found {
				d.Set("template", tempName)
				break
			}
		}
	}
	if found {
		d.Set("schemaID", schemaID)
	} else {
		d.SetId("")
		return fmt.Errorf("unable to find selector for given name")
	}
	return nil
}
