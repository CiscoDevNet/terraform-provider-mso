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

func resourceMSOSchemaSiteExternalEpg() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteExternalEpgCreate,
		Read:   resourceMSOSchemaSiteExternalEpgRead,
		Update: resourceMSOSchemaSiteExternalEpgUpdate,
		Delete: resourceMSOSchemaSiteExternalEpgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteExternalEpgImport,
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
			"site_id": &schema.Schema{
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
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaSiteExternalEpgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", get_attribute[0]))
	if err != nil {
		return nil, err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Sites found")
	}
	stateSiteId := get_attribute[2]
	found := false
	stateExternalEpg := get_attribute[4]
	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		if apiSiteId == stateSiteId {
			externalEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Externalepg list")
			}
			for j := 0; j < externalEpgCount; j++ {
				externalEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return nil, err
				}
				externalEpgRef := models.StripQuotes(externalEpgCont.S("externalEpgRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/externalEpgs/(.*)")
				match := re.FindStringSubmatch(externalEpgRef)
				if match[3] == stateExternalEpg {
					d.SetId(match[3])
					d.Set("external_epg_name", match[3])
					d.Set("schema_id", match[1])
					d.Set("template_name", match[2])
					d.Set("site_id", apiSiteId)

					l3outRef := models.StripQuotes(externalEpgCont.S("l3outRef").String())
					reL3out := regexp.MustCompile("/schemas/(.*)/templates/(.*)/l3outs/(.*)")
					matchL3out := reL3out.FindStringSubmatch(l3outRef)
					d.Set("l3out_name", matchL3out[3])

					found = true
					break
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return nil, fmt.Errorf("Unable to find the given Schema Site external EPG")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteExternalEpgCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site External EPG: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	templateName := d.Get("template_name").(string)
	l3outName := d.Get("l3out_name").(string)

	siteEpgMap := make(map[string]interface{})

	l3outRefMap := make(map[string]interface{})

	l3outRefMap["schemaId"] = schemaId
	l3outRefMap["templateName"] = templateName
	l3outRefMap["l3outName"] = l3outName

	siteEpgMap["l3outRef"] = l3outRefMap

	siteEpgMap["l3outDn"] = fmt.Sprintf("uni/tn-test_tenant/out-%s", l3outName)

	var ext_epg_schema_id, ext_epg_template_name string
	ext_epg_schema_id = schemaId
	ext_epg_template_name = templateName

	externalEpgRefMap := make(map[string]interface{})
	externalEpgRefMap["schemaId"] = ext_epg_schema_id
	externalEpgRefMap["templateName"] = ext_epg_template_name
	externalEpgRefMap["externalEpgName"] = externalEpgName

	siteEpgMap["externalEpgRef"] = externalEpgRefMap

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/-", siteId, templateName)
	siteExternalEpgStruct := models.NewSchemaSiteExternalEpg("add", path, siteEpgMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), siteExternalEpgStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteExternalEpgRead(d, m)
}

func resourceMSOSchemaSiteExternalEpgRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	stateSiteId := d.Get("site_id").(string)
	found := false
	stateExternalEpg := d.Get("external_epg_name").(string)
	for i := 0; i < count; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		if apiSiteId == stateSiteId {
			externalEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return fmt.Errorf("Unable to get External EPG list")
			}
			for j := 0; j < externalEpgCount; j++ {
				externalEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return err
				}
				externalEpgRef := models.StripQuotes(externalEpgCont.S("externalEpgRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/externalEpgs/(.*)")
				match := re.FindStringSubmatch(externalEpgRef)
				if match[3] == stateExternalEpg {
					d.SetId(match[3])
					d.Set("external_epg_name", match[3])
					d.Set("schema_id", match[1])
					d.Set("template_name", match[2])
					d.Set("site_id", apiSiteId)

					l3outRef := models.StripQuotes(externalEpgCont.S("l3outRef").String())
					reL3out := regexp.MustCompile("/schemas/(.*)/templates/(.*)/l3outs/(.*)")
					matchL3out := reL3out.FindStringSubmatch(l3outRef)
					d.Set("l3out_name", matchL3out[3])

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

func resourceMSOSchemaSiteExternalEpgUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Externalepg: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	l3outName := d.Get("l3out_name").(string)

	siteEpgMap := make(map[string]interface{})

	l3outRefMap := make(map[string]interface{})

	l3outRefMap["schemaId"] = schemaId
	l3outRefMap["templateName"] = templateName
	l3outRefMap["l3outName"] = l3outName

	siteEpgMap["l3outDn"] = fmt.Sprintf("uni/tn-test_tenant/out-%s", l3outName)

	var ext_epg_schema_id, ext_epg_template_name string
	ext_epg_schema_id = schemaId
	ext_epg_template_name = templateName

	externalEpgRefMap := make(map[string]interface{})
	externalEpgRefMap["schemaId"] = ext_epg_schema_id
	externalEpgRefMap["templateName"] = ext_epg_template_name
	externalEpgRefMap["externalEpgName"] = externalEpgName

	siteEpgMap["externalEpgRef"] = externalEpgRefMap
	siteEpgMap["l3outRef"] = l3outRefMap

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s", siteId, templateName, externalEpgName)
	siteExternalEpgStruct := models.NewSchemaSiteExternalEpg("replace", path, siteEpgMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), siteExternalEpgStruct)

	if err != nil {
		return err
	}

	return resourceMSOSchemaSiteExternalEpgRead(d, m)
}

func resourceMSOSchemaSiteExternalEpgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template External EPG: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)

	var l3outRefMap map[string]interface{}
	if tempVar, ok := d.GetOk("l3out_name"); ok {
		l3outName := tempVar.(string)
		var l3outSchemaID, l3outTemplate string
		if tmpVar, oki := d.GetOk("l3out_schema_id"); oki {
			l3outSchemaID = tmpVar.(string)
		} else {
			l3outSchemaID = schemaId
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

	var ext_epg_schema_id, ext_epg_template_name string
	ext_epg_schema_id = schemaId
	ext_epg_template_name = templateName

	externalEpgRefMap := make(map[string]interface{})
	externalEpgRefMap["schemaId"] = ext_epg_schema_id
	externalEpgRefMap["templateName"] = ext_epg_template_name
	externalEpgRefMap["externalEpgName"] = externalEpgName

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s", siteId, templateName, externalEpgName)
	siteExternalEpgStruct := models.NewSchemaSiteExternalEpg("remove", path, externalEpgRefMap)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), siteExternalEpgStruct)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
