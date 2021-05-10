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

func resourceMSOSchemaTemplateAnpEpgSelector() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateAnpEpgSelectorCreate,
		Read:   resourceMSOSchemaTemplateAnpEpgSelectorRead,
		Update: resourceMSOSchemaTemplateAnpEpgSelectorUpdate,
		Delete: resourceMSOSchemaTemplateAnpEpgSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateAnpEpgSelectorImport,
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

func resourceMSOSchemaTemplateAnpEpgSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	found := false
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaID := get_attribute[0]
	dn := get_attribute[8]
	template := get_attribute[2]
	anpName := get_attribute[4]
	epgName := get_attribute[6]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return nil, err
	}
	d.Set("schema_id", schemaID)

	tempCount, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}

	for i := 0; i < tempCount; i++ {
		tempcont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}

		tempName := models.StripQuotes(tempcont.S("name").String())
		if tempName == template {
			d.Set("template_name", tempName)
			anpCount, err := tempcont.ArrayCount("anps")
			if err != nil {
				return nil, err
			}

			for j := 0; j < anpCount; j++ {
				anpCont, err := tempcont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}

				currentanpName := models.StripQuotes(anpCont.S("name").String())
				if currentanpName == anpName {
					d.Set("anp_name", anpName)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, err
					}

					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
						}

						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							d.Set("epg_name", epgName)
							selectorCount, err := epgCont.ArrayCount("selectors")
							if err != nil {
								return nil, err
							}

							for l := 0; l < selectorCount; l++ {
								selectorCont, err := epgCont.ArrayElement(l, "selectors")
								if err != nil {
									return nil, err
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

						}
					}

				}
			}

		}
	}
	if !found {
		d.SetId("")
		return nil, fmt.Errorf("unable to find selector for given name")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateAnpEpgSelectorCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schematemplateanpepgselectorMap := make(map[string]interface{})

	schemaID := d.Get("schema_id").(string)

	template := d.Get("template_name").(string)

	anpName := d.Get("anp_name").(string)

	epgName := d.Get("epg_name").(string)

	name := d.Get("name").(string)

	expList := make([]interface{}, 0, 1)
	if exp, ok := d.GetOk("expressions"); ok {
		exps := exp.([]interface{})

		for _, val := range exps {
			exp := val.(map[string]interface{})

			expMap := make(map[string]interface{})

			expMap["key"] = exp["key"]

			expMap["operator"] = exp["operator"]

			if exp["value"] != nil {
				expMap["value"] = exp["value"]
			}

			expList = append(expList, expMap)
		}
	}

	schematemplateanpepgselectorMap["name"] = name
	schematemplateanpepgselectorMap["expressions"] = expList

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/selectors/-", template, anpName, epgName)

	schematemplateanpepgselector := models.NewSchemaTemplateAnpEpgSelector("add", path, schematemplateanpepgselectorMap)

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schematemplateanpepgselector)
	if err != nil {
		return err
	}

	index, err := fetchIndexSelector(cont, template, anpName, epgName, name)
	if err != nil {
		return err
	}
	if index == -1 {
		d.SetId("")
		return fmt.Errorf("The given selector name is not found")
	} else {
		d.SetId(name)
	}
	return resourceMSOSchemaTemplateAnpEpgSelectorRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schematemplateanpepgselectorMap := make(map[string]interface{})

	dn := d.Id()

	schemaID := d.Get("schema_id").(string)

	template := d.Get("template_name").(string)

	anpName := d.Get("anp_name").(string)

	epgName := d.Get("epg_name").(string)

	name := d.Get("name").(string)

	expList := make([]interface{}, 0, 1)
	if exp, ok := d.GetOk("expressions"); ok {
		exps := exp.([]interface{})

		for _, val := range exps {
			exp := val.(map[string]interface{})

			expMap := make(map[string]interface{})

			expMap["key"] = exp["key"]

			expMap["operator"] = exp["operator"]

			if exp["value"] != "" {
				expMap["value"] = exp["value"]
			}

			expList = append(expList, expMap)
		}
	}

	schematemplateanpepgselectorMap["name"] = name
	schematemplateanpepgselectorMap["expressions"] = expList

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/selectors/%s", template, anpName, epgName, dn)

	schematemplateanpepgselector := models.NewSchemaTemplateAnpEpgSelector("replace", path, schematemplateanpepgselectorMap)

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schematemplateanpepgselector)
	if err != nil {
		return err
	}

	index, err := fetchIndexSelector(cont, template, anpName, epgName, name)
	if err != nil {
		return err
	}

	if index == -1 {
		d.SetId("")
		return fmt.Errorf("The given selector name is not found")
	} else {
		d.SetId(name)
	}
	return resourceMSOSchemaTemplateAnpEpgSelectorRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	found := false
	msoClient := m.(*client.Client)

	dn := d.Id()
	schemaID := d.Get("schema_id").(string)
	template := d.Get("template_name").(string)
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
				d.Set("template_name", tempName)
				break
			}
		}
	}
	if found {
		d.Set("schemaID", schemaID)
	} else {
		d.SetId("")
	}
	return nil
}

func resourceMSOSchemaTemplateAnpEpgSelectorDelete(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)

	dn := d.Id()
	schemaID := d.Get("schema_id").(string)
	template := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s/selectors/%s", template, anpName, epgName, dn)

	schematemplateanpepgselector := models.NewSchemaTemplateAnpEpgSelector("remove", path, nil)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schematemplateanpepgselector)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func fetchIndexSelector(cont *container.Container, templateName, anpName, epgName, name string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return index, fmt.Errorf("No Template found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return index, err
		}

		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return index, fmt.Errorf("No Anp found")
			}

			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return index, err
				}

				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				if currentAnpName == anpName {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return index, fmt.Errorf("No Epg found")
					}

					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return index, err
						}

						currentEpgName := models.StripQuotes(epgCont.S("name").String())
						if currentEpgName == epgName {
							subnetCount, err := epgCont.ArrayCount("selectors")
							if err != nil {
								return index, fmt.Errorf("No selectors found")
							}

							for s := 0; s < subnetCount; s++ {
								subnetCont, err := epgCont.ArrayElement(s, "selectors")
								if err != nil {
									return index, err
								}

								currentName := models.StripQuotes(subnetCont.S("name").String())
								if currentName == name {
									index = s
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

		}
		if found {
			break
		}
	}
	return index, nil
}
