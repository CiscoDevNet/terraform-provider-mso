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

func resourceMSOTemplateExtenalepg() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOTemplateExtenalepgCreate,
		Read:   resourceMSOTemplateExtenalepgRead,
		Update: resourceMSOTemplateExtenalepgUpdate,
		Delete: resourceMSOTemplateExtenalepgDelete,

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
		}),
	}
}

func resourceMSOTemplateExtenalepgCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	extenalepgName := d.Get("external_epg_name").(string)
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
	}

	path := fmt.Sprintf("/templates/%s/externalEpgs/-", templateName)
	externalepgStruct := models.NewTemplateExternalepg("add", path, extenalepgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)

	if err != nil {
		return err
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
					d.SetId(apiExternalepg)
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
					if l3outRef != "{}" {
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
					if anpRef != "{}" {
						tokens := strings.Split(anpRef, "/")
						d.Set("anp_name", tokens[len(tokens)-1])
						d.Set("anp_schema_id", tokens[len(tokens)-5])
						d.Set("anp_template_name", tokens[len(tokens)-3])
					} else {
						d.Set("anp_name", "")
						d.Set("anp_schema_id", "")
						d.Set("anp_template_name", "")
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
	extenalepgName := d.Get("external_epg_name").(string)
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
	}

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s", templateName, extenalepgName)
	externalepgStruct := models.NewTemplateExternalepg("replace", path, extenalepgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)

	if err != nil {
		return err
	}
	return resourceMSOTemplateExtenalepgRead(d, m)
}

func resourceMSOTemplateExtenalepgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg: Beginning Update")
	msoClient := m.(*client.Client)

	schemaID := d.Get("schema_id").(string)
	extenalepgName := d.Get("external_epg_name").(string)
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
	}

	path := fmt.Sprintf("/templates/%s/externalEpgs/%s", templateName, extenalepgName)
	externalepgStruct := models.NewTemplateExternalepg("remove", path, extenalepgName, displayName, extEpgType, preferredGroup, vrfRefMap, l3outRefMap, anpRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaID), externalepgStruct)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
