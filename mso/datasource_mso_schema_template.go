package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func datasourceMSOSchemaTemplate() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaTemplateRead,

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
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func datasourceMSOSchemaTemplateRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

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
		return fmt.Errorf("Template of specified name not found")
	}

	dataCon := cont.S("templates").Index(count)
	d.SetId(models.StripQuotes(dataCon.S("name").String()))
	d.Set("name", models.StripQuotes(dataCon.S("name").String()))
	d.Set("display_name", models.StripQuotes(dataCon.S("displayName").String()))
	d.Set("tenant_id", models.StripQuotes(dataCon.S("tenantId").String()))

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
