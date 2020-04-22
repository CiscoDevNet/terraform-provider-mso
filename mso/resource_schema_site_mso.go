package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"

	// "github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMSOSchemaSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteCreate,
		Update: resourceMSOSchemaSiteUpdate,
		Read:   resourceMSOSchemaSiteRead,
		Delete: resourceMSOSchemaSiteDelete,

		// Importer: &schema.ResourceImporter{
		//     State: resourceMSOSchemaSiteImport,
		// },

		SchemaVersion: 1,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"site_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteCreate(d *schema.ResourceData, m interface{}) error {
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	siteId := d.Get("site_id").(string)

	schemasite := models.NewSchemaSite("add", siteId, templateName)

	cont, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemasite)
	if err != nil {
		return err
	}

	id := cont.S("id")
	log.Println("Id value", id)
	d.SetId(fmt.Sprintf("%v", id))
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
	d.SetId(fmt.Sprintf("%v", cont.S("id")))
	d.Set("schema", cont.S("schema").String())
	d.Set("templates", cont.S("templates").String())
	d.Set("sites", cont.S("sites").String())

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaSiteUpdate(d *schema.ResourceData, m interface{}) error {
	// log.Printf("[DEBUG] CloudApplicationcontainer: Beginning Update")

	// msoClient := m.(*client.Client)

	// schemasiteAttr := models.SchemaSiteAttributes{}

	// if d.HasChange("schema") {
	// 	if schema, ok := d.GetOk("schema"); ok {
	// 		schemasiteAttr.Template = schema.(string)
	// 	}
	// }

	// if d.HasChange("templates") {
	// 	if templates, ok := d.GetOk("templates"); ok {
	// 		schemasiteAttr.Template = templates.(string)
	// 	}
	// }

	// if d.HasChange("sites") {
	// 	if site, ok := d.GetOk("site"); ok {
	// 		schemasiteAttr.Site = site.(string)
	// 	}
	// }
	// schemasite := models.NewSchemaSite(schemasiteAttr)
	// cont, err := msoClient.PatchbyID("api/v1/schemas/sites/"+d.Id(), schemasite)

	// if err != nil {
	// 	return err
	// }

	// id := cont.S("id")
	// log.Println("Id value", id)
	// d.SetId(fmt.Sprintf("%v", id))
	// log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	// return resourceMSOSchemaSiteRead(d, m)
	return nil

}

func resourceMSOSchemaSiteDelete(d *schema.ResourceData, m interface{}) error {
	// log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	// msoClient := m.(*client.Client)
	// dn := d.Id()
	// err := msoClient.DeletebyId("api/v1/schemas/sites/" + dn)
	// if err != nil {
	// 	return err
	// }

	// log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	// d.SetId("")
	// return err
	return nil
}
