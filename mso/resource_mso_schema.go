package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"

	// "github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMSOSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaCreate,
		Update: resourceMSOSchemaUpdate,
		Read:   resourceMSOSchemaRead,
		Delete: resourceMSOSchemaDelete,

		// Importer: &schema.ResourceImporter{
		// 	State: resourceMSOSchemaImport,
		// },

		SchemaVersion: 1,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceMSOSchemaCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	templateName := d.Get("template_name").(string)
	tenandId := d.Get("tenant_id").(string)

	schemaApp := models.NewSchema("", name, templateName, tenandId)

	cont, err := msoClient.Save("api/v1/schemas", schemaApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Schema Creation finished successfully", d.Id())

	return resourceMSOSchemaRead(d, m)
}

func resourceMSOSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	templateName := d.Get("template_name").(string)
	tenandId := d.Get("tenant_id").(string)

	schemaApp := models.NewSchema(d.Id(), name, templateName, tenandId)

	cont, err := msoClient.Put(fmt.Sprintf("api/v1/schemas/%s", d.Id()), schemaApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Schema Creation finished successfully", d.Id())

	return resourceMSOSchemaRead(d, m)
	return nil

}

func resourceMSOSchemaRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	dn := d.Id()

	con, err := msoClient.GetViaURL("api/v1/schemas/" + dn)
	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("name", models.StripQuotes(con.S("displayName").String()))
	count, err := con.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	stateTenant := d.Get("tenant_id").(string)
	found := false
	for i := 0; i < count; i++ {
		tempCont, err := con.ArrayElement(i, "templates")

		if err != nil {
			return fmt.Errorf("Unable to parse the template list")
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())
		apiTenant := models.StripQuotes(tempCont.S("tenantId").String())
		log.Printf("apitemp %s apiten %s statetemp %s stateten %s", apiTemplate, apiTenant, stateTemplate, stateTenant)
		if apiTemplate == stateTemplate && apiTenant == stateTenant {
			d.Set("template_name", apiTemplate)
			d.Set("tenant_id", apiTenant)
			found = true
			break
		}
	}
	if !found {
		d.Set("template_name", "")
		d.Set("tenant_id", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId("api/v1/schemas/" + dn)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

func toStringList(configured interface{}) []string {
	vs := make([]string, 0, 1)
	val, ok := configured.(string)
	if ok && val != "" {
		vs = append(vs, val)
	}
	return vs
}
