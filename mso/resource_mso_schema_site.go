package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteCreate,
		Read:   resourceMSOSchemaSiteRead,
		Delete: resourceMSOSchemaSiteDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteImport,
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
		}),
	}
}

func resourceMSOSchemaSiteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	name := get_attribute[2]
	con, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/sites"))
	if err != nil {
		return nil, err
	}

	data := con.S("sites").Data().([]interface{})
	var flag bool
	var count int
	for _, info := range data {
		val := info.(map[string]interface{})
		if val["name"].(string) == name {
			flag = true
			break
		}
		count = count + 1
	}
	if flag != true {
		return nil, fmt.Errorf("Site of specified name not found")
	}

	dataCon := con.S("sites").Index(count)
	stateSiteId := models.StripQuotes(dataCon.S("id").String())

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}

	countSites, err := cont.ArrayCount("sites")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}

	found := false

	for i := 0; i < countSites; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return nil, err
		}
		apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSiteId == stateSiteId {
			d.SetId(apiSiteId)
			d.Set("schema_id", schemaId)
			d.Set("site_id", apiSiteId)
			d.Set("template_name", apiTemplate)
			found = true
		}

	}

	if !found {
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)

	schemasite := models.NewSchemaSite("add", "/sites/-", siteId, templateName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemasite)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", siteId))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaSiteRead(d, m)
}

func resourceMSOSchemaSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	stateSiteId := d.Get("site_id").(string)
	stateTemplate := d.Get("template_name").(string)
	found := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSiteId == stateSiteId && apiTemplate == stateTemplate {
			d.SetId(apiSiteId)
			d.Set("schema_id", schemaId)
			d.Set("site_id", apiSiteId)
			d.Set("template_name", apiTemplate)
			found = true
		}

	}

	if !found {
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaSiteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)

	schemasite := models.NewSchemaSite("remove", fmt.Sprintf("/sites/%s-%s", siteId, templateName), siteId, templateName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemasite)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
