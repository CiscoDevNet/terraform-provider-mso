package mso

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceSchemaTemplateExternalEPGSelector() *schema.Resource {
	return &schema.Resource{
		Read: datasourceSchemaTemplateExternalEPGSelectorRead,

		Schema: map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"external_epg_name": &schema.Schema{
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

func datasourceSchemaTemplateExternalEPGSelectorRead(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)

	found := false

	schemaID := d.Get("schema_id").(string)
	template := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	name := d.Get("name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No templates found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return fmt.Errorf("Error fetching template")
		}

		tempName := models.StripQuotes(tempCont.S("name").String())
		if tempName == template {
			extrEpgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("no externalEpgs found")
			}

			for j := 0; j < extrEpgCount; j++ {
				extrEpgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return fmt.Errorf("Error fetching external Epg")
				}

				extrEpgName := models.StripQuotes(extrEpgCont.S("name").String())
				if extrEpgName == externalEpgName {
					selectorCount, err := extrEpgCont.ArrayCount("selectors")
					if err != nil {
						return fmt.Errorf("No selectors found")
					}

					for k := 0; k < selectorCount; k++ {
						selectorCont, err := extrEpgCont.ArrayElement(k, "selectors")
						if err != nil {
							return fmt.Errorf("Error fetching selector")
						}

						selectorName := models.StripQuotes(selectorCont.S("name").String())
						if selectorName == name {
							d.SetId(name)
							d.Set("name", selectorName)
							exps := selectorCont.S("expressions").Data().([]interface{})

							expressionList := make([]interface{}, 0, 1)
							for _, val := range exps {
								tp := val.(map[string]interface{})
								expMap := make(map[string]interface{})

								expMap["key"] = "ipAddress"
								expMap["operator"] = "equals"
								if tp["value"] != nil {
									expMap["value"] = tp["value"]
								}
								expressionList = append(expressionList, expMap)
							}
							d.Set("expressions", expressionList)
							found = true
							break
						}
					}
				}
				if found {
					d.Set("external_epg_name", externalEpgName)
					break
				}
			}
		}
		if found {
			d.Set("template_name", tempName)
			break
		}
	}
	if found {
		d.Set("schema_id", schemaID)
	} else {
		d.SetId("")
		return fmt.Errorf("External Epg Selector not found for given name")
	}
	return nil
}
