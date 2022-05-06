package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchema() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"display_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"tenant_id": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
		}),
	}
}

func datasourceMSOSchemaRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	name := d.Get("name").(string)
	con, err := msoClient.GetViaURL("api/v1/schemas")
	if err != nil {
		return err
	}
	data := con.S("schemas").Data().([]interface{})
	var flag bool
	var count int
	for _, info := range data {
		val := info.(map[string]interface{})
		if val["displayName"].(string) == name {
			flag = true
			break
		}
		count = count + 1
	}
	if flag != true {
		return fmt.Errorf("Schema of specified name not found")
	}

	dataCon := con.S("schemas").Index(count)
	d.SetId(models.StripQuotes(dataCon.S("id").String()))
	d.Set("name", models.StripQuotes(dataCon.S("displayName").String()))
	countTemplate, err := dataCon.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	templates := make([]interface{}, 0)
	for i := 0; i < countTemplate; i++ {
		templatesCont, err := dataCon.ArrayElement(i, "templates")
		if err != nil {
			return fmt.Errorf("Unable to parse the templates list")
		}
		map_template := make(map[string]interface{})
		map_template["name"] = models.StripQuotes(templatesCont.S("name").String())
		map_template["display_name"] = models.StripQuotes(templatesCont.S("displayName").String())
		map_template["tenant_id"] = models.StripQuotes(templatesCont.S("tenantId").String())
		templates = append(templates, map_template)

	}
	d.Set("template", templates)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
