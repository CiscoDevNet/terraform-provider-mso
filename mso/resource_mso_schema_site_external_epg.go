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
			"l3out_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"l3out_on_apic": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	for i := 0; i < count && !found; i++ {
		siteCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSiteId := models.StripQuotes(siteCont.S("siteId").String())

		if apiSiteId == stateSiteId {
			externalEpgCount, err := siteCont.ArrayCount("externalEpgs")
			if err != nil {
				return nil, fmt.Errorf("Unable to get External EPG list")
			}
			for j := 0; j < externalEpgCount; j++ {
				externalEpgCont, err := siteCont.ArrayElement(j, "externalEpgs")
				if err != nil {
					return nil, err
				}
				externalEpgRef := models.StripQuotes(externalEpgCont.S("externalEpgRef").String())
				re := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/externalEpgs/(.*)")
				match := re.FindStringSubmatch(externalEpgRef)
				log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead externalEpgRef: %s match: %s", externalEpgRef, match)
				if len(match) >= 4 {
					if match[3] == stateExternalEpg {
						d.SetId(match[3])
						d.Set("external_epg_name", match[3])
						d.Set("schema_id", match[1])
						d.Set("template_name", match[2])
						d.Set("site_id", apiSiteId)

						l3outRef := models.StripQuotes(externalEpgCont.S("l3outRef").String())
						l3outDn := models.StripQuotes(externalEpgCont.S("l3outDn").String())
						if l3outRef != "{}" && l3outRef != "" {
							reL3out := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/l3outs/(.*)")
							matchL3out := reL3out.FindStringSubmatch(l3outRef)
							log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead l3outRef: %s matchL3out: %s", l3outRef, matchL3out)
							if len(matchL3out) >= 4 {
								d.Set("l3out_name", matchL3out[3])
								d.Set("l3out_schema_id", matchL3out[1])
								d.Set("l3out_template_name", matchL3out[2])
								d.Set("l3out_on_apic", false)
							} else {
								return nil, fmt.Errorf("Error in parsing l3outRef to get L3Out name")
							}
						} else if l3outDn != "{}" && l3outDn != "" {
							reL3out := regexp.MustCompile("uni/tn-(.*?)/out-(.*)")
							matchL3out := reL3out.FindStringSubmatch(l3outDn)
							log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead l3outDn: %s matchL3out: %s", l3outDn, matchL3out)
							if len(matchL3out) >= 2 {
								d.Set("l3out_name", matchL3out[1])
								d.Set("l3out_on_apic", true)
							} else {
								return fmt.Errorf("Error in parsing l3outDn to get L3Out name")
							}
						}

						found = true
						break
					}
				} else {
					return nil, fmt.Errorf("Error in parsing externalEpgRef to get External EPG name")
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
	l3outTemplate := d.Get("l3out_template_name").(string)
	l3outSchema := d.Get("l3out_schema_id").(string)
	l3outOnApic := d.Get("l3out_on_apic").(bool)

	siteEpgMap := make(map[string]interface{})

	if l3outName != "" {
		// Get tenant name
		tenantName, err := GetTenantNameViaTemplateName(msoClient, schemaId, templateName)
		if err != nil {
			return err
		}
		
		l3outRefMap := make(map[string]interface{})

		if l3outOnApic {
			siteEpgMap["l3outRef"] = ""
		} else {		
			l3outRefMap["schemaId"] = l3outSchema
			l3outRefMap["templateName"] = l3outTemplate
			l3outRefMap["l3outName"] = l3outName

			siteEpgMap["l3outRef"] = l3outRefMap
		} 
		
		siteEpgMap["l3outDn"] = fmt.Sprintf("uni/tn-%s/out-%s", tenantName, l3outName)
	} else {
		siteEpgMap["l3outDn"] = ""
	}

	externalEpgRefMap := make(map[string]interface{})
	externalEpgRefMap["schemaId"] = schemaId
	externalEpgRefMap["templateName"] = templateName
	externalEpgRefMap["externalEpgName"] = externalEpgName

	siteEpgMap["externalEpgRef"] = externalEpgRefMap

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	if versionInt != 1 {
		path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s", siteId, templateName, externalEpgName)
		siteExternalEpgStruct := models.NewSchemaSiteExternalEpg("replace", path, siteEpgMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), siteExternalEpgStruct)
	}

	if versionInt == 1 || err != nil {
		path := fmt.Sprintf("/sites/%s-%s/externalEpgs/-", siteId, templateName)
		siteExternalEpgStruct := models.NewSchemaSiteExternalEpg("add", path, siteEpgMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), siteExternalEpgStruct)
	}

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
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}
	stateSiteId := d.Get("site_id").(string)
	found := false
	stateExternalEpg := d.Get("external_epg_name").(string)
	for i := 0; i < count && !found; i++ {
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
				re := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/externalEpgs/(.*)")
				match := re.FindStringSubmatch(externalEpgRef)
				log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead externalEpgRef: %s match: %s", externalEpgRef, match)
				if len(match) >= 4 {
					if match[3] == stateExternalEpg {
						d.SetId(match[3])
						d.Set("external_epg_name", match[3])
						d.Set("schema_id", match[1])
						d.Set("template_name", match[2])
						d.Set("site_id", apiSiteId)

						l3outRef := models.StripQuotes(externalEpgCont.S("l3outRef").String())
						l3outDn := models.StripQuotes(externalEpgCont.S("l3outDn").String())
						if l3outRef != "{}" && l3outRef != "" {
							reL3out := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/l3outs/(.*)")
							matchL3out := reL3out.FindStringSubmatch(l3outRef)
							log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead l3outRef: %s matchL3out: %s", l3outRef, matchL3out)
							if len(matchL3out) >= 4 {
								d.Set("l3out_name", matchL3out[3])
								d.Set("l3out_schema_id", matchL3out[1])
								d.Set("l3out_template_name", matchL3out[2])
								d.Set("l3out_on_apic", false)
							} else {
								return fmt.Errorf("Error in parsing l3outRef to get L3Out name")
							}
						} else if l3outDn != "{}" && l3outDn != "" {
							reL3out := regexp.MustCompile("uni/tn-(.*?)/out-(.*)")
							matchL3out := reL3out.FindStringSubmatch(l3outDn)
							log.Printf("[TRACE] resourceMSOSchemaSiteExternalEpgRead l3outDn: %s matchL3out: %s", l3outDn, matchL3out)
							if len(matchL3out) >= 2 {
								d.Set("l3out_name", matchL3out[1])
								d.Set("l3out_on_apic", true)
							} else {
								return fmt.Errorf("Error in parsing l3outDn to get L3Out name")
							}
						}

						found = true
						break
					}
				} else {
					return fmt.Errorf("Error in parsing externalEpgRef to get External EPG name")
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
	l3outTemplate := d.Get("l3out_template_name").(string)
	l3outSchema := d.Get("l3out_schema_id").(string)
	l3outOnApic := d.Get("l3out_on_apic").(bool)

	siteEpgMap := make(map[string]interface{})

	if l3outName != "" {
		// Get tenant name
		tenantName, err := GetTenantNameViaTemplateName(msoClient, schemaId, templateName)
		if err != nil {
			return err
		}

		l3outRefMap := make(map[string]interface{})

		if l3outOnApic {
			siteEpgMap["l3outRef"] = ""
		} else {

			l3outRefMap["schemaId"] = l3outSchema
			l3outRefMap["templateName"] = l3outTemplate
			l3outRefMap["l3outName"] = l3outName

			siteEpgMap["l3outRef"] = l3outRefMap
		} 
		siteEpgMap["l3outDn"] = fmt.Sprintf("uni/tn-%s/out-%s", tenantName, l3outName)
	} else {
		siteEpgMap["l3outDn"] = ""
	}

	externalEpgRefMap := make(map[string]interface{})
	externalEpgRefMap["schemaId"] = schemaId
	externalEpgRefMap["templateName"] = templateName
	externalEpgRefMap["externalEpgName"] = externalEpgName

	siteEpgMap["externalEpgRef"] = externalEpgRefMap

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s", siteId, templateName, externalEpgName)
	siteExternalEpgStruct := models.NewSchemaSiteExternalEpg("replace", path, siteEpgMap)

	_, patchErr := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), siteExternalEpgStruct)
	if patchErr != nil {
		return patchErr
	}

	return resourceMSOSchemaSiteExternalEpgRead(d, m)
}

func resourceMSOSchemaSiteExternalEpgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template External EPG: Beginning Delete")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)

	externalEpgRefMap := make(map[string]interface{})
	externalEpgRefMap["schemaId"] = schemaId
	externalEpgRefMap["templateName"] = templateName
	externalEpgRefMap["externalEpgName"] = externalEpgName

	path := fmt.Sprintf("/sites/%s-%s/externalEpgs/%s", siteId, templateName, externalEpgName)
	siteExternalEpgStruct := models.NewSchemaSiteExternalEpg("remove", path, externalEpgRefMap)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), siteExternalEpgStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	d.SetId("")
	return nil
}

// Gets tenant name by doing the following
// GET and loop through all the schemas and check if the schema is present ("api/v1/schemas/list-identity")
// GET and loop through all the templates in the schema and check if the template is present
// If template present then get tenantId from template contents
// GET tenant_name from tenantId "api/v1/tenants/{id}"
func GetTenantNameViaTemplateName(msoClient *client.Client, id string, tempName string) (string, error) {
	cont, err := msoClient.GetViaURL("api/v1/schemas/list-identity")
	if err != nil {
		return "", err
	}
	schemaCount, err := cont.ArrayCount("schemas")
	if err != nil {
		return "", err
	}

	for i := 0; i < schemaCount; i++ {
		schemaCont, err := cont.ArrayElement(i, "schemas")
		if err != nil {
			return "", err
		}
		schemaId := models.StripQuotes(schemaCont.S("id").String())

		if schemaId == id {
			allTemplates := schemaCont.S("templates").Data().([]interface{})

			for _, info := range allTemplates {
				template := info.(map[string]interface{})
				if tempName == template["name"] {
					tenantId := template["tenantId"]
					tenantCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/tenants/%v", tenantId))

					if err != nil {
						return "", err
					}

					tenantMap := tenantCont.Data().(map[string]interface{})
					tenantName := tenantMap["name"].(string)
					return tenantName, nil
				}

			}
		}

	}
	return "", fmt.Errorf(tempName)
}
