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
			"template_name": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "see template",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "see template",
			},
			"template": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_type": {
							Type:     schema.TypeString,
							Computed: true,
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
	d.Set("description", models.StripQuotes(con.S("description").String()))

	// Commented block so it can be used when templates container is initialized as empty list
	// Currently in NDO 4.x the templates container is initialized as null instead of empty list
	//  so when no templates are provided during create or import it is impossible to PATCH add a template
	// NDO 4.x allows us to specify schema without templates thus skipping error of no templates provided and version >=4.x
	// versionInt, err := msoClient.CompareVersion("4.0.0.0")
	// if err != nil {
	// 	return err
	// }
	// countTemplate, err := dataCon.ArrayCount("templates")
	// if err != nil && versionInt == 1 {
	countTemplate, err := dataCon.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	templates := make([]interface{}, 0)
	for i := 0; i < countTemplate; i++ {
		tempCont, err := dataCon.ArrayElement(i, "templates")

		if err != nil {
			return fmt.Errorf("Unable to parse the template list")
		}
		if i == 0 {
			d.Set("template_name", models.StripQuotes(tempCont.S("name").String()))
			d.Set("tenant_id", models.StripQuotes(tempCont.S("tenantId").String()))
		}
		map_template := make(map[string]interface{})
		map_template["name"] = models.StripQuotes(tempCont.S("name").String())
		map_template["display_name"] = models.StripQuotes(tempCont.S("displayName").String())
		map_template["tenant_id"] = models.StripQuotes(tempCont.S("tenantId").String())
		map_template["description"] = models.StripQuotes(tempCont.S("description").String())
		map_template["template_type"] = getSchemaTemplateType(tempCont)
		templates = append(templates, map_template)

	}
	d.Set("template", templates)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
