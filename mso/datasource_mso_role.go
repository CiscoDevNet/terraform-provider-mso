package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSORole() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSORoleRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"read_permissions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"write_permissions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func datasourceMSORoleRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	con, err := msoClient.GetViaURL("api/v1/roles")
	if err != nil {
		return err
	}

	data := con.S("roles").Data().([]interface{})
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
		return fmt.Errorf("Role of specified name not found")
	}

	dataCon := con.S("roles").Index(count)
	d.SetId(models.StripQuotes(dataCon.S("id").String()))
	d.Set("name", models.StripQuotes(dataCon.S("name").String()))
	d.Set("display_name", models.StripQuotes(dataCon.S("displayName").String()))
	d.Set("description", models.StripQuotes(dataCon.S("description").String()))
	d.Set("read_permissions", dataCon.S("readPermissions").Data().([]interface{}))
	d.Set("write_permissions", dataCon.S("writePermissions").Data().([]interface{}))

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
