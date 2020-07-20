package mso

import (
	"fmt"
	"strings"

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
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"epg_name": &schema.Schema{
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

			"expressions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},

						"operator": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"equals",
								"notEquals",
								"in",
								"notIn",
								"keyExist",
								"keyNotExist",
							}, false),
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

func datasourceSchemaSiteApnEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	found := false
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	template := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	name := d.Get("name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	siteCount, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	for i := 0; i < siteCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}

		currentSite := models.StripQuotes(siteCont.S("siteId").String())
		currentTemp := models.StripQuotes(siteCont.S("templateName").String())

		if currentTemp == template && currentSite == siteID {
			anpCount, err := siteCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("No Anp found")
			}

			for j := 0; j < anpCount; j++ {
				anpCont, err := siteCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}

				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				tokens := strings.Split(anpRef, "/")
				currentAnpName := tokens[len(tokens)-1]
				if currentAnpName == anpName {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("No Epg found")
					}

					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}

						epgRef := models.StripQuotes(epgCont.S("epgRef").String())
						tokensEpg := strings.Split(epgRef, "/")
						currentEpgName := tokensEpg[len(tokensEpg)-1]
						if currentEpgName == epgName {
							selectorCount, err := epgCont.ArrayCount("selectors")
							if err != nil {
								return fmt.Errorf("No selectors found")
							}

							for s := 0; s < selectorCount; s++ {
								selectorCont, err := epgCont.ArrayElement(s, "selectors")
								if err != nil {
									return err
								}

								currentName := models.StripQuotes(selectorCont.S("name").String())
								if currentName == name {
									found = true
									d.SetId(name)
									d.Set("name", currentName)
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
			d.Set("site_id", siteID)
			d.Set("template_name", template)
			break
		}
	}
	if found {
		d.Set("schema_id", schemaID)
	} else {
		d.SetId("")
		return fmt.Errorf("No Site ANP EPG selector found for given name")
	}
	return nil
}
