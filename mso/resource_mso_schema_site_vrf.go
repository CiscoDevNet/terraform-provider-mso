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

func resourceMSOSchemaSiteVrf() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteVrfCreate,
		Read:   resourceMSOSchemaSiteVrfRead,
		Delete: resourceMSOSchemaSiteVrfDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteVrfImport,
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
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaSiteVrfImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
	stateVrf := get_attribute[4]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return nil, fmt.Errorf("Unable to get Vrf list")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return nil, err
				}
				vrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
				match := re.FindStringSubmatch(vrfRef)
				if match[3] == stateVrf {
					d.SetId(match[3])
					d.Set("vrf_name", match[3])
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
		return nil, fmt.Errorf("Unable to find Site Vrf %s", stateVrf)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaSiteVrfCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var vrf_schema_id, vrf_template_name string
	vrf_schema_id = schemaId
	vrf_template_name = templateName

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	versionInt, err := msoClient.CompareVersion("4.0.0.0")
	if err != nil {
		return err
	}

	if versionInt != 1 {
		path := fmt.Sprintf("/sites/%s-%s/vrfs/%s", siteId, templateName, vrfName)
		vrfStruct := models.NewSchemaSiteVrf("replace", path, vrfRefMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), vrfStruct)
	}

	if versionInt == 1 || err != nil {
		path := fmt.Sprintf("/sites/%s-%s/vrfs/-", siteId, templateName)
		vrfStruct := models.NewSchemaSiteVrf("add", path, vrfRefMap)
		_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), vrfStruct)
	}

	if err != nil {
		return err
	}

	return resourceMSOSchemaSiteVrfRead(d, m)
}

func resourceMSOSchemaSiteVrfRead(d *schema.ResourceData, m interface{}) error {
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
	stateSite := d.Get("site_id").(string)
	found := false
	stateVrf := d.Get("vrf_name").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return fmt.Errorf("Unable to get Vrf list")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return err
				}
				vrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
				re := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
				match := re.FindStringSubmatch(vrfRef)
				if match[3] == stateVrf {
					d.SetId(match[3])
					d.Set("vrf_name", match[3])
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

func resourceMSOSchemaSiteVrfDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf: Beginning Deletion")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)

	var vrf_schema_id, vrf_template_name string
	vrf_schema_id = schemaId
	vrf_template_name = templateName
	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s", siteId, templateName, vrfName)
	vrfStruct := models.NewSchemaSiteVrf("remove", path, vrfRefMap)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), vrfStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")
	return nil
}
