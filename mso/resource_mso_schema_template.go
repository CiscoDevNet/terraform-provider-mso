package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceMSOSchemaTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateCreate,
		Read:   resourceMSOSchemaTemplateRead,
		Update: resourceMSOSchemaTemplateUpdate,
		Delete: resourceMSOSchemaTemplateDelete,

		SchemaVersion: 1,

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

func resourceMSOSchemaTemplateCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	Name := d.Get("name").(string)
	tenantId := d.Get("tenant_id").(string)
	displayName := d.Get("display_name").(string)

	schematemplate := models.NewSchemaTemplate("add", "/templates/-", tenantId, Name, displayName)

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schematemplate)
	if err != nil {
		return err
	}

	id := cont.S("id")
	log.Println("Id value", id)
	d.SetId(fmt.Sprintf("%v", id))
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

		cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schematemplate)
		if err != nil {
			return err
		}

		id := cont.S("id")
		log.Println("Id value", id)
		d.SetId(fmt.Sprintf("%v", id))
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

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schematemplate)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
