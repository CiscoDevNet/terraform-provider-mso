package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSchemaTemplateExternalEPGSelector() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchemaTemplateExternalEPGSelectorCreate,
		Update: resourceSchemaTemplateExternalEPGSelectorUpdate,
		Read:   resourceSchemaTemplateExternalEPGSelectorRead,
		Delete: resourceSchemaTemplateExternalEPGSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSchemaTemplateExternalEPGSelectorImport,
		},

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

func resourceSchemaTemplateExternalEPGSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	found := false
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	template := get_attribute[2]
	externalEpgName := get_attribute[4]
	name := get_attribute[6]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No templates found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, fmt.Errorf("Error fetching template")
		}

		tempName := models.StripQuotes(tempCont.S("name").String())
		if tempName == template {
			extrEpgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return nil, fmt.Errorf("no externalEpgs found")
			}

			for j := 0; j < extrEpgCount; j++ {
				extrEpgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return nil, fmt.Errorf("Error fetching external Epg")
				}

				extrEpgName := models.StripQuotes(extrEpgCont.S("name").String())
				if extrEpgName == externalEpgName {
					selectorCount, err := extrEpgCont.ArrayCount("selectors")
					if err != nil {
						return nil, fmt.Errorf("No selectors found")
					}

					for k := 0; k < selectorCount; k++ {
						selectorCont, err := extrEpgCont.ArrayElement(k, "selectors")
						if err != nil {
							return nil, fmt.Errorf("Error fetching selector")
						}

						selectorName := models.StripQuotes(selectorCont.S("name").String())
						if selectorName == name {
							d.SetId(get_attribute[6])
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
		d.Set("schema_id", schemaId)
	} else {
		d.SetId("")
		return nil, fmt.Errorf("External Epg Selector not found for given name")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceSchemaTemplateExternalEPGSelectorCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaTemplateextrepgSelectorMap := make(map[string]interface{})

	schemaID := d.Get("schema_id").(string)

	template := d.Get("template_name").(string)

	extrnalEPGName := d.Get("external_epg_name").(string)

	name := d.Get("name").(string)

	expList := make([]interface{}, 0, 1)
	if exp, ok := d.GetOk("expressions"); ok {
		exps := exp.([]interface{})

		for _, val := range exps {
			tp := val.(map[string]interface{})

			expMap := make(map[string]interface{})

			expMap["key"] = "ipAddress"
			expMap["operator"] = "equals"
			if tp["value"] != nil {
				expMap["value"] = tp["value"]
			}

			expList = append(expList, expMap)
		}
	}

	schemaTemplateextrepgSelectorMap["name"] = name
	schemaTemplateextrepgSelectorMap["expressions"] = expList

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/selectors/-", template, extrnalEPGName)

	schemaTemplateExternalEPGSelector := models.NewSchemaTemplateExternalEPGSelector("add", path, schemaTemplateextrepgSelectorMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("/api/v1/schemas/%s", schemaID), schemaTemplateExternalEPGSelector)
	if err != nil {
		return err
	}

	d.SetId(name)
	return resourceSchemaTemplateExternalEPGSelectorRead(d, m)
}

func resourceSchemaTemplateExternalEPGSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaTemplateextrepgSelectorMap := make(map[string]interface{})

	dn := d.Id()

	schemaID := d.Get("schema_id").(string)

	template := d.Get("template_name").(string)

	externalEpgName := d.Get("external_epg_name").(string)

	name := d.Get("name").(string)

	expList := make([]interface{}, 0, 1)
	if exp, ok := d.GetOk("expressions"); ok {
		exps := exp.([]interface{})

		for _, val := range exps {
			tp := val.(map[string]interface{})

			expMap := make(map[string]interface{})
			expMap["key"] = "ipAddress"
			expMap["operator"] = "equals"
			if tp["value"] != nil {
				expMap["value"] = tp["value"]
			}

			expList = append(expList, expMap)
		}
	}

	schemaTemplateextrepgSelectorMap["name"] = name
	schemaTemplateextrepgSelectorMap["expressions"] = expList

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/selectors/%s", template, externalEpgName, dn)

	schemaTemplateExternalEpgSelector := models.NewSchemaTemplateExternalEPGSelector("replace", path, schemaTemplateextrepgSelectorMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaTemplateExternalEpgSelector)
	if err != nil {
		return err
	}

	d.SetId(name)
	return resourceSchemaTemplateExternalEPGSelectorRead(d, m)
}

func resourceSchemaTemplateExternalEPGSelectorRead(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)

	found := false

	dn := d.Id()
	schemaID := d.Get("schema_id").(string)
	template := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)

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
						if selectorName == dn {
							d.SetId(dn)
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
	}
	return nil
}

func resourceSchemaTemplateExternalEPGSelectorDelete(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)

	dn := d.Id()

	schemaID := d.Get("schema_id").(string)

	template := d.Get("template_name").(string)

	externalEpgName := d.Get("external_epg_name").(string)

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s/selectors/%s", template, externalEpgName, dn)

	schemaTemplateExternalEpgSelector := models.NewSchemaTemplateExternalEPGSelector("remove", path, nil)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemaTemplateExternalEpgSelector)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func checkAvailExternalEpgSelector(cont *container.Container, template, externalEpgName, name string) (int, error) {
	found := false
	index := -1

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return index, fmt.Errorf("No templates found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return index, fmt.Errorf("Error fetching template")
		}

		tempName := models.StripQuotes(tempCont.S("name").String())
		if tempName == template {
			extrEpgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return index, fmt.Errorf("no externalEpgs found")
			}

			for j := 0; j < extrEpgCount; j++ {
				extrEpgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return index, fmt.Errorf("Error fetching external Epg")
				}

				extrEpgName := models.StripQuotes(extrEpgCont.S("name").String())
				if extrEpgName == externalEpgName {
					selectorCount, err := extrEpgCont.ArrayCount("selectors")
					if err != nil {
						return index, fmt.Errorf("No selectors found")
					}

					for k := 0; k < selectorCount; k++ {
						selectorCont, err := extrEpgCont.ArrayElement(k, "selectors")
						if err != nil {
							return index, fmt.Errorf("Error fetching selector")
						}

						selectorName := models.StripQuotes(selectorCont.S("name").String())
						if selectorName == name {
							index = k
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
		if found {
			break
		}
	}
	return index, nil
}
