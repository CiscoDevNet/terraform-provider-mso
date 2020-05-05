package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func datasourceMSOSchemaTemplateAnp() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaTemplateAnpRead,

		SchemaVersion: version,
		Schema: (map[string]*schema.Schema{

			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
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
func datasourceMSOSchemaTemplateAnpRead(d *schema.ResourceData, m interface{}) error {
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

	templateName := d.Get("template").(string)
	anpName := d.Get("name").(string)
	found := false

	for i := 0; i < count; i++ {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
			anpCount, err := tempCont.ArrayCount("anps")

			if err != nil {
				return fmt.Errorf("No Anp found")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")

				if err != nil {
					return err
				}
				currentAnpName := models.StripQuotes(anpCont.S("name").String())
				log.Println("currentanpname", currentAnpName)
				if currentAnpName == anpName {
					log.Println("found correct anpname")
					d.SetId(currentAnpName)
					d.Set("name", currentAnpName)
					if anpCont.Exists("displayName") {
						d.Set("display_name", models.StripQuotes(anpCont.S("displayName").String()))
					}
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	if !found {
		d.SetId("")
	    return fmt.Errorf("The ANP is not found")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
