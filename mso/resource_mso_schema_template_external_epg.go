package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOTemplateExtenalepg() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateExtenalepgCreate,
		Read:   resourceMSOTemplateExtenalepgRead,
		Update: resourceMSOTemplateExtenalepgUpdate,
		Delete: resourceMSOTemplateExtenalepgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOTemplateExtenalepgImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
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
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"external_epg_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"on-premise",
					"cloud",
				}, false),
			},
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"l3out_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"l3out_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"include_in_preferred_group": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"site_id": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"selector_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"selector_ip": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOTemplateExtenalepgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := get_attribute[2]
	found := false
	stateExternalepg := get_attribute[4]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			externalepgCount, err := tempCont.ArrayCount("externalEpgs")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Externalepg list")
			}
			for j := 0; j < externalepgCount; j++ {
				externalepgCont, err := tempCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return nil, err
				}
				apiExternalepg := models.StripQuotes(externalepgCont.S("name").String())
				if apiExternalepg == stateExternalepg {
					id := fmt.Sprintf("/schemas/%s/templates/%s/externalEpgs/%s", schemaId, apiTemplate, apiExternalepg)
					d.SetId(id)
					d.Set("external_epg_name", apiExternalepg)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(externalepgCont.S("displayName").String()))
					d.Set("external_epg_type", models.StripQuotes(externalepgCont.S("extEpgType").String()))
					if externalepgCont.Exists("preferredGroup") {
						d.Set("include_in_preferred_group", externalepgCont.S("preferredGroup").Data().(bool))
					} else {
						d.Set("include_in_preferred_group", false)
					}

					vrfRef := models.StripQuotes(externalepgCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])
					l3outRef := models.StripQuotes(externalepgCont.S("l3outRef").String())
					if l3outRef != "{}" && l3outRef != "" {
						reL3out := regexp.MustCompile("/schemas/(.*)/templates/(.*)/l3outs/(.*)")
						matchL3out := reL3out.FindStringSubmatch(l3outRef)
						d.Set("l3out_name", matchL3out[3])
						d.Set("l3out_schema_id", matchL3out[1])
						d.Set("l3out_template_name", matchL3out[2])
					} else {
						d.Set("l3out_name", "")
						d.Set("l3out_schema_id", "")
						d.Set("l3out_template_name", "")
					}

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

					epgType := d.Get("external_epg_type").(string)
					if epgType == "cloud" {
						selList := externalepgCont.S("selectors").Data().([]interface{})

						selector := selList[0].(map[string]interface{})
						d.Set("selector_name", selector["name"])
						expList := selector["expressions"].([]interface{})
						exp := expList[0].(map[string]interface{})
						d.Set("selector_ip", exp["value"])
					} else {
						d.Set("site_id", make([]interface{}, 0, 1))
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
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOTemplateExtenalepgCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var extEpgType string
	if tempVar, ok := d.GetOk("external_epg_type"); ok {
		extEpgType = tempVar.(string)
	} else {
		extEpgType = "on-premise"
	}
	var vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	preferredGroup := d.Get("include_in_preferred_group").(bool)

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	var l3outRefMap map[string]interface{}
	if tempVar, ok := d.GetOk("l3out_name"); ok {
		l3outName := tempVar.(string)
		var l3outSchemaID, l3outTemplate string
		if tmpVar, oki := d.GetOk("l3out_schema_id"); oki {
			l3outSchemaID = tmpVar.(string)
		} else {
			l3outSchemaID = schemaID
		}

		if tpVar, okj := d.GetOk("l3out_template_name"); okj {
			l3outTemplate = tpVar.(string)
		} else {
			l3outTemplate = templateName
		}

		l3outRefMap = make(map[string]interface{})

		l3outRefMap["schemaId"] = l3outSchemaID
		l3outRefMap["templateName"] = l3outTemplate
		l3outRefMap["l3outName"] = l3outName

	}

	anpRefMap := make(map[string]interface{})
	if aName, ok := d.GetOk("anp_name"); ok {
		anpName := aName.(string)

		var anpSchemaID, anpTemplateName string
		if schID, ok := d.GetOk("anp_schema_id"); ok {
			anpSchemaID = schID.(string)
		} else {
			anpSchemaID = schemaID
		}

		if tmpName, ok := d.GetOk("anp_template_name"); ok {
			anpTemplateName = tmpName.(string)
		} else {
			anpTemplateName = templateName
		}

		anpRefMap["schemaId"] = anpSchemaID
		anpRefMap["templateName"] = anpTemplateName
		anpRefMap["anpName"] = anpName
	} else {
		anpRefMap = nil
	}

	platform := msoClient.GetPlatform()

	if extEpgType == "cloud" {
		var selectorName string
		if selName, ok := d.GetOk("selector_name"); ok {
			selectorName = selName.(string)
		} else {
			if platform == "mso" {
				return fmt.Errorf("selector_name attribute is required for cloud configuration")
			}
			selectorName = ""
		}

		var selectorIP string
		if selIP, ok := d.GetOk("selector_ip"); ok {
			selectorIP = selIP.(string)
		} else {
			if platform == "mso" {
				return fmt.Errorf("selector_ip attribute is required for cloud configuration")
			}
			selectorIP = ""
		}

		selectorList := make([]interface{}, 0, 1)
		if selectorName != "" && selectorIP != "" {
			expressionList := make([]interface{}, 0, 1)
			selectionMap := make(map[string]interface{})
			expMap := make(map[string]interface{})

			expMap["key"] = "ipAddress"
			expMap["operator"] = "equals"
			expMap["value"] = selectorIP
			expressionList = append(expressionList, expMap)

			selectionMap["name"] = selectorName
			selectionMap["expressions"] = expressionList
			selectorList = append(selectorList, selectionMap)
		}

		pathTemp := fmt.Sprintf("/templates/%s/externalEpgs/-", templateName)
		externalepgStruct := models.NewTemplateExternalepg("add", pathTemp, externalEpgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap, selectorList)

		structList := make([]models.Model, 0, 1)
		structList = append(structList, externalepgStruct)

		var sites []interface{}
		if site, ok := d.GetOk("site_id"); ok {
			sites = site.([]interface{})
			if platform == "nd" {
				return fmt.Errorf("site_id attribute is not supported when running on NDO. Use mso_schema_site_external_epg to add the external EPG to a site.")
			}
		} else {
			if platform == "mso" {
				return fmt.Errorf("site_id attribute is required for cloud configuration")
			}
		}
		for _, site := range sites {
			siteEpgMap := make(map[string]interface{})
			epgRefMap := make(map[string]interface{})

			epgRefMap["schemaId"] = schemaID
			epgRefMap["templateName"] = templateName
			epgRefMap["externalEpgName"] = externalEpgName

			siteEpgMap["externalEpgRef"] = epgRefMap
			siteEpgMap["l3outDn"] = "l3out"

			pathSite := fmt.Sprintf("/sites/%s-%s/externalEpgs/-", site.(string), templateName)
			siteExternalepgStruct := models.NewSchemaSiteExternalEpg("add", pathSite, siteEpgMap)
			structList = append(structList, siteExternalepgStruct)
		}

		d.Partial(true)
		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), structList...)
		if err != nil {
			return err
		}
		d.Partial(false)

	} else {
		path := fmt.Sprintf("/templates/%s/externalEpgs/-", templateName)
		externalepgStruct := models.NewTemplateExternalepg("add", path, externalEpgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap, nil)

		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)
		if err != nil {
			return err
		}
	}
	return resourceMSOTemplateExtenalepgRead(d, m)
}

func resourceMSOTemplateExtenalepgRead(d *schema.ResourceData, m interface{}) error {
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
	stateExternalepg := d.Get("external_epg_name")
	for i := 0; i < count; i++ {
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
					id := fmt.Sprintf("/schemas/%s/templates/%s/externalEpgs/%s", schemaId, apiTemplate, apiExternalepg)
					d.SetId(id)
					d.Set("external_epg_name", apiExternalepg)
					d.Set("schema_id", schemaId)
					d.Set("template_name", apiTemplate)
					d.Set("display_name", models.StripQuotes(externalepgCont.S("displayName").String()))
					d.Set("external_epg_type", models.StripQuotes(externalepgCont.S("extEpgType").String()))
					if externalepgCont.Exists("preferredGroup") {
						d.Set("include_in_preferred_group", externalepgCont.S("preferredGroup").Data().(bool))
					} else {
						d.Set("include_in_preferred_group", false)
					}

					vrfRef := models.StripQuotes(externalepgCont.S("vrfRef").String())
					re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
					match := re.FindStringSubmatch(vrfRef)
					d.Set("vrf_name", match[3])
					d.Set("vrf_schema_id", match[1])
					d.Set("vrf_template_name", match[2])
					l3outRef := models.StripQuotes(externalepgCont.S("l3outRef").String())
					if l3outRef != "{}" && l3outRef != "" {
						reL3out := regexp.MustCompile("/schemas/(.*)/templates/(.*)/l3outs/(.*)")
						matchL3out := reL3out.FindStringSubmatch(l3outRef)
						d.Set("l3out_name", matchL3out[3])
						d.Set("l3out_schema_id", matchL3out[1])
						d.Set("l3out_template_name", matchL3out[2])
					} else {
						d.Set("l3out_name", "")
						d.Set("l3out_schema_id", "")
						d.Set("l3out_template_name", "")
					}

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

					epgType := d.Get("external_epg_type").(string)
					if epgType == "cloud" {
						selList := externalepgCont.S("selectors").Data().([]interface{})

						selector := selList[0].(map[string]interface{})
						d.Set("selector_name", selector["name"])
						expList := selector["expressions"].([]interface{})
						exp := expList[0].(map[string]interface{})
						d.Set("selector_ip", exp["value"])
					} else {
						d.Set("site_id", make([]interface{}, 0, 1))
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
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOTemplateExtenalepgUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	preferredGroup := d.Get("include_in_preferred_group").(bool)

	var extEpgType string
	if tempVar, ok := d.GetOk("external_epg_type"); ok {
		extEpgType = tempVar.(string)
	} else {
		extEpgType = "on-premise"
	}

	var vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	var l3outRefMap map[string]interface{}
	if tempVar, ok := d.GetOk("l3out_name"); ok {
		l3outName := tempVar.(string)
		var l3outSchemaID, l3outTemplate string
		if tmpVar, oki := d.GetOk("l3out_schema_id"); oki {
			l3outSchemaID = tmpVar.(string)
		} else {
			l3outSchemaID = schemaID
		}

		if tpVar, okj := d.GetOk("l3out_template_name"); okj {
			l3outTemplate = tpVar.(string)
		} else {
			l3outTemplate = templateName
		}

		l3outRefMap = make(map[string]interface{})

		l3outRefMap["schemaId"] = l3outSchemaID
		l3outRefMap["templateName"] = l3outTemplate
		l3outRefMap["l3outName"] = l3outName

	}

	anpRefMap := make(map[string]interface{})
	if aName, ok := d.GetOk("anp_name"); ok {
		anpName := aName.(string)

		var anpSchemaID, anpTemplateName string
		if schID, ok := d.GetOk("anp_schema_id"); ok {
			anpSchemaID = schID.(string)
		} else {
			anpSchemaID = schemaID
		}

		if tmpName, ok := d.GetOk("anp_template_name"); ok {
			anpTemplateName = tmpName.(string)
		} else {
			anpTemplateName = templateName
		}

		anpRefMap["schemaId"] = anpSchemaID
		anpRefMap["templateName"] = anpTemplateName
		anpRefMap["anpName"] = anpName
	} else {
		anpRefMap = nil
	}

	platform := msoClient.GetPlatform()

	if extEpgType == "cloud" {
		var selectorName string
		if selName, ok := d.GetOk("selector_name"); ok {
			selectorName = selName.(string)
		} else {
			if platform == "mso" {
				return fmt.Errorf("selector_name attribute is required for cloud configuration")
			}
			selectorName = ""
		}

		var selectorIP string
		if selIP, ok := d.GetOk("selector_ip"); ok {
			selectorIP = selIP.(string)
		} else {
			if platform == "mso" {
				return fmt.Errorf("selector_ip attribute is required for cloud configuration")
			}
			selectorIP = ""
		}

		selectorList := make([]interface{}, 0, 1)
		if selectorName != "" && selectorIP != "" {
			expressionList := make([]interface{}, 0, 1)
			selectionMap := make(map[string]interface{})
			expMap := make(map[string]interface{})

			expMap["key"] = "ipAddress"
			expMap["operator"] = "equals"
			expMap["value"] = selectorIP
			expressionList = append(expressionList, expMap)

			selectionMap["name"] = selectorName
			selectionMap["expressions"] = expressionList
			selectorList = append(selectorList, selectionMap)
		}

		pathTemp := fmt.Sprintf("/templates/%s/externalEpgs/%s", templateName, externalEpgName)
		externalepgStruct := models.NewTemplateExternalepg("replace", pathTemp, externalEpgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap, selectorList)

		structList := make([]models.Model, 0, 1)
		structList = append(structList, externalepgStruct)

		var sites []interface{}
		if site, ok := d.GetOk("site_id"); ok {
			sites = site.([]interface{})
			if platform == "nd" {
				return fmt.Errorf("site_id attribute is not supported when running on NDO. Use mso_schema_site_external_epg to add the external EPG to a site.")
			}
		} else {
			if platform == "mso" {
				return fmt.Errorf("site_id attribute is required for cloud configuration")
			}
		}
		for _, site := range sites {
			siteEpgMap := make(map[string]interface{})
			epgRefMap := make(map[string]interface{})

			epgRefMap["schemaId"] = schemaID
			epgRefMap["templateName"] = templateName
			epgRefMap["externalEpgName"] = externalEpgName

			siteEpgMap["externalEpgRef"] = epgRefMap
			siteEpgMap["l3outDn"] = "l3out"

			cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return err
			}
			flag, err := checkSiteExternalEpg(cont, site.(string), templateName, externalEpgName)
			if err != nil {
				return err
			}
			var op string
			var pathSite string
			if !flag {
				op = "add"
				pathSite = fmt.Sprintf("/sites/%s-%s/externalEpgs/-", site.(string), templateName)
			} else {
				op = "replace"
				pathSite = fmt.Sprintf("/sites/%s-%s/externalEpgs/%s", site.(string), templateName, externalEpgName)
			}

			siteExternalepgStruct := models.NewSchemaSiteExternalEpg(op, pathSite, siteEpgMap)
			structList = append(structList, siteExternalepgStruct)
		}

		d.Partial(true)
		_, err1 := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), structList...)
		if err1 != nil {
			return err1
		}
		d.Partial(false)

	} else {
		path := fmt.Sprintf("/templates/%s/externalEpgs/%s", templateName, externalEpgName)
		externalepgStruct := models.NewTemplateExternalepg("replace", path, externalEpgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap, nil)

		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)
		if err != nil {
			return err
		}
	}
	return resourceMSOTemplateExtenalepgRead(d, m)
}

func resourceMSOTemplateExtenalepgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	displayName := d.Get("display_name").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	preferredGroup := d.Get("include_in_preferred_group").(bool)

	var extEpgType string
	if tempVar, ok := d.GetOk("external_epg_type"); ok {
		extEpgType = tempVar.(string)
	} else {
		extEpgType = "on-premise"
	}
	var vrf_schema_id, vrf_template_name string

	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaID
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	var l3outRefMap map[string]interface{}
	if tempVar, ok := d.GetOk("l3out_name"); ok {
		l3outName := tempVar.(string)
		var l3outSchemaID, l3outTemplate string
		if tmpVar, oki := d.GetOk("l3out_schema_id"); oki {
			l3outSchemaID = tmpVar.(string)
		} else {
			l3outSchemaID = schemaID
		}

		if tpVar, okj := d.GetOk("l3out_template_name"); okj {
			l3outTemplate = tpVar.(string)
		} else {
			l3outTemplate = templateName
		}

		l3outRefMap = make(map[string]interface{})

		l3outRefMap["schemaId"] = l3outSchemaID
		l3outRefMap["templateName"] = l3outTemplate
		l3outRefMap["l3outName"] = l3outName

	}

	anpRefMap := make(map[string]interface{})
	if aName, ok := d.GetOk("anp_name"); ok {
		anpName := aName.(string)

		var anpSchemaID, anpTemplateName string
		if schID, ok := d.GetOk("anp_schema_id"); ok {
			anpSchemaID = schID.(string)
		} else {
			anpSchemaID = schemaID
		}

		if tmpName, ok := d.GetOk("anp_template_name"); ok {
			anpTemplateName = tmpName.(string)
		} else {
			anpTemplateName = templateName
		}

		anpRefMap["schemaId"] = anpSchemaID
		anpRefMap["templateName"] = anpTemplateName
		anpRefMap["anpName"] = anpName
	} else {
		anpRefMap = nil
	}

	platform := msoClient.GetPlatform()

	if extEpgType == "cloud" {
		var selectorName string
		if selName, ok := d.GetOk("selector_name"); ok {
			selectorName = selName.(string)
		} else {
			if platform == "mso" {
				return fmt.Errorf("selector_name attribute is required for cloud configuration")
			}
			selectorName = ""
		}

		var selectorIP string
		if selIP, ok := d.GetOk("selector_ip"); ok {
			selectorIP = selIP.(string)
		} else {
			if platform == "mso" {
				return fmt.Errorf("selector_ip attribute is required for cloud configuration")
			}
			selectorIP = ""
		}

		selectorList := make([]interface{}, 0, 1)
		if selectorName != "" && selectorIP != "" {
			expressionList := make([]interface{}, 0, 1)
			selectionMap := make(map[string]interface{})
			expMap := make(map[string]interface{})

			expMap["key"] = "ipAddress"
			expMap["operator"] = "equals"
			expMap["value"] = selectorIP
			expressionList = append(expressionList, expMap)

			selectionMap["name"] = selectorName
			selectionMap["expressions"] = expressionList
			selectorList = append(selectorList, selectionMap)
		}

		pathTemp := fmt.Sprintf("/templates/%s/externalEpgs/%s", templateName, externalEpgName)
		externalepgStruct := models.NewTemplateExternalepg("remove", pathTemp, externalEpgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap, selectorList)

		structList := make([]models.Model, 0, 1)
		structList = append(structList, externalepgStruct)

		var sites []interface{}
		if site, ok := d.GetOk("site_id"); ok {
			sites = site.([]interface{})
			if platform == "nd" {
				return fmt.Errorf("site_id attribute is not supported when running on NDO. Use mso_schema_site_external_epg to add the external EPG to a site.")
			}
		} else {
			if platform == "mso" {
				return fmt.Errorf("site_id attribute is required for cloud configuration")
			}
		}
		for _, site := range sites {
			siteEpgMap := make(map[string]interface{})
			epgRefMap := make(map[string]interface{})

			epgRefMap["schemaId"] = schemaID
			epgRefMap["templateName"] = templateName
			epgRefMap["externalEpgName"] = externalEpgName

			siteEpgMap["externalEpgRef"] = epgRefMap
			siteEpgMap["l3outDn"] = "l3out"

			pathSite := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s", site.(string), templateName, externalEpgName)
			siteExternalepgStruct := models.NewSchemaSiteExternalEpg("remove", pathSite, siteEpgMap)
			structList = append(structList, siteExternalepgStruct)
		}

		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), structList...)
		if err != nil {
			return err
		}

	} else {
		path := fmt.Sprintf("/templates/%s/externalEpgs/%s", templateName, externalEpgName)
		externalepgStruct := models.NewTemplateExternalepg("remove", path, externalEpgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap, nil)

		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}

func checkSiteExternalEpg(cont *container.Container, site, template, epgName string) (bool, error) {
	flag := false

	siteCount, err := cont.ArrayCount("sites")
	if err != nil {
		return flag, err
	}

	for i := 0; i < siteCount; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return flag, err
		}

		dn := models.StripQuotes(siteCont.S("siteId").String())
		temp := models.StripQuotes(siteCont.S("templateName").String())
		if dn == site && temp == template {
			epgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return flag, err
			}

			for j := 0; j < epgCount; j++ {
				epgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return flag, err
				}

				epgRef := models.StripQuotes(epgCont.S("externalEpgRef").String())
				tokens := strings.Split(epgRef, "/")
				if epgName == tokens[len(tokens)-1] {
					flag = true
					break
				}
			}
		}
		if flag {
			break
		}
	}
	return flag, nil
}
