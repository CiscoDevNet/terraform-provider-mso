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

func resourceMSOSchemaTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateCreate,
		Read:   resourceMSOSchemaTemplateRead,
		Update: resourceMSOSchemaTemplateUpdate,
		Delete: resourceMSOSchemaTemplateDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateImport,
		},

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"tenant_id": &schema.Schema{
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
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func resourceMSOSchemaTemplateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	get_attribute := strings.Split(d.Id(), "/")
	msoClient := m.(*client.Client)
	name := get_attribute[2]
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	d.Set("schema_id", schemaId)
	data := cont.S("templates").Data().([]interface{})
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
		return nil, fmt.Errorf("Template of specified name not found")
	}

	dataCon := cont.S("templates").Index(count)
	d.SetId(models.StripQuotes(dataCon.S("name").String()))
	d.Set("name", models.StripQuotes(dataCon.S("name").String()))
	d.Set("display_name", models.StripQuotes(dataCon.S("displayName").String()))
	d.Set("tenant_id", models.StripQuotes(dataCon.S("tenantId").String()))

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	Name := d.Get("name").(string)
	tenantId := d.Get("tenant_id").(string)
	displayName := d.Get("display_name").(string)

	schematemplate := models.NewSchemaTemplate("add", "/templates/-", tenantId, Name, displayName)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schematemplate)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v", Name))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateRead(d, m)
}

func resourceMSOSchemaTemplateRead(d *schema.ResourceData, m interface{}) error {
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
	stateTenantId := d.Get("tenant_id").(string)
	stateTemplateName := d.Get("name").(string)
	stateTemplateDisplayName := d.Get("display_name").(string)

	found := false

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTenantId := models.StripQuotes(tempCont.S("tenantId").String())
		apiTemplateName := models.StripQuotes(tempCont.S("name").String())
		apiTemplateDisplayName := models.StripQuotes(tempCont.S("displayName").String())

		if apiTenantId == stateTenantId && apiTemplateName == stateTemplateName && apiTemplateDisplayName == stateTemplateDisplayName {
			d.SetId(apiTemplateName)
			d.Set("tenant_id", apiTenantId)
			d.Set("name", apiTemplateName)
			d.Set("display_name", apiTemplateDisplayName)
			found = true
		}

	}

	if !found {
		d.SetId("")
		d.Set("tenant_id", "")
		d.Set("name", "")
		d.Set("display_name", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	Name := d.Get("name").(string)

	if d.HasChange("display_name") {
		tenantId := d.Get("tenant_id").(string)
		displayName := d.Get("display_name").(string)

		schematemplate := models.NewSchemaTemplate("replace", fmt.Sprintf("/templates/%s", Name), tenantId, Name, displayName)

		_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schematemplate)
		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("%v", Name))
		log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	}

	return resourceMSOSchemaTemplateRead(d, m)
}

func resourceMSOSchemaTemplateDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	tenantId := d.Get("tenant_id").(string)
	templateName := d.Get("name").(string)
	templateDisplayName := d.Get("display_name").(string)

	schematemplate := models.NewSchemaTemplate("remove", fmt.Sprintf("/templates/%s", templateName), tenantId, templateName, templateDisplayName)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schematemplate)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
