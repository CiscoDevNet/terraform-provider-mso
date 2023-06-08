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

func resourceMSOSchemaSiteAnpEpgSelector() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchemaSiteApnEpgSelectorCreate,
		Update: resourceSchemaSiteApnEpgSelectorUpdate,
		Read:   resourceSchemaSiteApnEpgSelectorRead,
		Delete: resourceSchemaSiteApnEpgSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSchemaSiteApnEpgSelectorImport,
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

func resourceSchemaSiteApnEpgSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	found := false
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaID := get_attribute[0]
	siteID := get_attribute[2]
	template := get_attribute[4]
	anpName := get_attribute[6]
	epgName := get_attribute[8]
	name := get_attribute[10]

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return nil, err
	}

	siteCount, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}

	for i := 0; i < siteCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}

		currentSite := models.StripQuotes(siteCont.S("siteId").String())
		currentTemp := models.StripQuotes(siteCont.S("templateName").String())

		if currentTemp == template && currentSite == siteID {
			anpCount, err := siteCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("No Anp found")
			}

			for j := 0; j < anpCount; j++ {
				anpCont, err := siteCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}

				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				tokens := strings.Split(anpRef, "/")
				currentAnpName := tokens[len(tokens)-1]
				if currentAnpName == anpName {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, fmt.Errorf("No Epg found")
					}

					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
						}

						epgRef := models.StripQuotes(epgCont.S("epgRef").String())
						tokensEpg := strings.Split(epgRef, "/")
						currentEpgName := tokensEpg[len(tokensEpg)-1]
						if currentEpgName == epgName {
							selectorCount, err := epgCont.ArrayCount("selectors")
							if err != nil {
								return nil, fmt.Errorf("No selectors found")
							}

							for s := 0; s < selectorCount; s++ {
								selectorCont, err := epgCont.ArrayElement(s, "selectors")
								if err != nil {
									return nil, err
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
		return nil, fmt.Errorf("No Site ANP EPG selector found for given name")
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceSchemaSiteApnEpgSelectorCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)

	schemasiteanpepgselectorMap := make(map[string]interface{})

	schemaID := d.Get("schema_id").(string)

	siteID := d.Get("site_id").(string)

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

	schemasiteanpepgselectorMap["name"] = name
	schemasiteanpepgselectorMap["expressions"] = expList

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/selectors/-", siteID, template, anpName, epgName)

	schemasiteanpepgselector := models.NewSchemaTemplateAnpEpgSelector("add", path, schemasiteanpepgselectorMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemasiteanpepgselector)
	if err != nil {
		return err
	}

	d.SetId(name)
	return resourceSchemaSiteApnEpgSelectorRead(d, m)
}

func resourceSchemaSiteApnEpgSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	dn := d.Id()

	schemasiteanpepgselectorMap := make(map[string]interface{})

	schemaID := d.Get("schema_id").(string)

	siteID := d.Get("site_id").(string)

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

	schemasiteanpepgselectorMap["name"] = name
	schemasiteanpepgselectorMap["expressions"] = expList

	contGet, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	indexGet, err := checkSchemaSiteApnEpgSelector(contGet, siteID, template, anpName, epgName, dn)
	if err != nil {
		return err
	}
	if indexGet == -1 {
		return fmt.Errorf("No selectors found")
	}

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/selectors/%v", siteID, template, anpName, epgName, indexGet)

	schemasiteanpepgselector := models.NewSchemaTemplateAnpEpgSelector("replace", path, schemasiteanpepgselectorMap)

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemasiteanpepgselector)
	if err != nil {
		return err
	}

	d.SetId(name)
	return resourceSchemaSiteApnEpgSelectorRead(d, m)
}

func resourceSchemaSiteApnEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	found := false
	msoClient := m.(*client.Client)

	dn := d.Id()
	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	template := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return errorForObjectNotFound(err, dn, cont, d)
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
								if currentName == dn {
									found = true
									d.SetId(dn)
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
	}
	return nil
}

func resourceSchemaSiteApnEpgSelectorDelete(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	dn := d.Id()

	schemasiteanpepgselectorMap := make(map[string]interface{})

	schemaID := d.Get("schema_id").(string)

	siteID := d.Get("site_id").(string)

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

	schemasiteanpepgselectorMap["name"] = name
	schemasiteanpepgselectorMap["expressions"] = expList

	contGet, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
	if err != nil {
		return err
	}

	indexGet, err := checkSchemaSiteApnEpgSelector(contGet, siteID, template, anpName, epgName, dn)
	if err != nil {
		return err
	}
	if indexGet == -1 {
		d.SetId("")
		return nil
	}

	path := fmt.Sprintf("/sites/%s-%s/anps/%s/epgs/%s/selectors/%v", siteID, template, anpName, epgName, indexGet)

	schemasiteanpepgselector := models.NewSchemaTemplateAnpEpgSelector("remove", path, schemasiteanpepgselectorMap)

	response, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), schemasiteanpepgselector)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err1 != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err1
	}

	d.SetId("")
	return nil
}

func checkSchemaSiteApnEpgSelector(cont *container.Container, siteID, templateName, anpName, epgName, name string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return index, fmt.Errorf("No Sites found")
	}

	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return index, err
		}

		currentSite := models.StripQuotes(siteCont.S("siteId").String())
		currentTemp := models.StripQuotes(siteCont.S("templateName").String())

		if currentTemp == templateName && currentSite == siteID {
			anpCount, err := siteCont.ArrayCount("anps")
			if err != nil {
				return index, fmt.Errorf("No Anp found")
			}

			for j := 0; j < anpCount; j++ {
				anpCont, err := siteCont.ArrayElement(j, "anps")
				if err != nil {
					return index, err
				}

				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				tokens := strings.Split(anpRef, "/")
				currentAnpName := tokens[len(tokens)-1]
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

						epgRef := models.StripQuotes(epgCont.S("epgRef").String())
						tokensEpg := strings.Split(epgRef, "/")
						currentEpgName := tokensEpg[len(tokensEpg)-1]
						if currentEpgName == epgName {
							selectorCount, err := epgCont.ArrayCount("selectors")
							if err != nil {
								return index, fmt.Errorf("No selectors found")
							}

							for s := 0; s < selectorCount; s++ {
								selectorCont, err := epgCont.ArrayElement(s, "selectors")
								if err != nil {
									return index, err
								}

								currentName := models.StripQuotes(selectorCont.S("name").String())
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
