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

func resourceMSOSchemaSiteAnp() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteAnpCreate,
		Read:   resourceMSOSchemaSiteAnpRead,
		Delete: resourceMSOSchemaSiteAnpDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteAnpImport,
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
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaSiteAnpImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	stateSite := get_attribute[2]
	found := false
	stateAnp := get_attribute[4]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {

			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)
				if match[3] == stateAnp {
					d.SetId(match[3])
					d.Set("anp_name", match[3])
					d.Set("schema_id", match[1])
					d.Set("template_name", match[2])
					d.Set("site_id", apiSite)
					found = true
					break
				}
			}
		}
	}

	if !found {
		d.SetId("")
		return nil, fmt.Errorf("Unable to find Site Anp %s", stateAnp)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteAnpCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)

	var anp_schema_id, anp_template_name string
	anp_schema_id = schemaId
	anp_template_name = templateName

	anpRefMap := make(map[string]interface{})
	anpRefMap["schemaId"] = anp_schema_id
	anpRefMap["templateName"] = anp_template_name
	anpRefMap["anpName"] = anpName

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	if versionInt != 1 {
		path := fmt.Sprintf("/sites/%s-%s/anps/%s", siteId, templateName, anpName)
		anpStruct := models.NewSchemaSiteAnp("replace", path, anpRefMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
	}

	if versionInt == 1 || err != nil {
		path := fmt.Sprintf("/sites/%s-%s/anps/-", siteId, templateName)
		anpStruct := models.NewSchemaSiteAnp("add", path, anpRefMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)
	}

	if err != nil {
		return err
	}

	return resourceMSOSchemaSiteAnpRead(d, m)
}

func resourceMSOSchemaSiteAnpRead(d *schema.ResourceData, m interface{}) error {
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
	stateSite := d.Get("site_id").(string)
	found := false
	stateAnp := d.Get("anp_name").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get Anp list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				anpRef := models.StripQuotes(anpCont.S("anpRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/anps/(.*)")
				match := re.FindStringSubmatch(anpRef)
				if match[3] == stateAnp {
					d.SetId(match[3])
					d.Set("anp_name", match[3])
					d.Set("schema_id", match[1])
					d.Set("template_name", match[2])
					d.Set("site_id", apiSite)
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

func resourceMSOSchemaSiteAnpDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Anp: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)

	var anp_schema_id, anp_template_name string
	anp_schema_id = schemaId
	anp_template_name = templateName
	anpRefMap := make(map[string]interface{})
	anpRefMap["schemaId"] = anp_schema_id
	anpRefMap["templateName"] = anp_template_name
	anpRefMap["anpName"] = anpName

	path := fmt.Sprintf("/sites/%s-%s/anps/%s", siteId, templateName, anpName)
	anpStruct := models.NewSchemaSiteAnp("remove", path, anpRefMap)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")
	return nil
}
