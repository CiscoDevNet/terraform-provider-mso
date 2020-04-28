package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceMSOLabel() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOLabelRead,

		SchemaVersion: 1,

		Schema: (map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func datasourceMSOLabelRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	name := d.Get("label").(string)
	con, err := msoClient.GetViaURL("api/v1/labels")
	if err != nil {
		return err
	}
	data := con.S("labels").Data().([]interface{})
	var flag bool
	var cnt int
	for _, info := range data {
		val := info.(map[string]interface{})
		if val["displayName"].(string) == name {
			flag = true
			break
		}
		cnt = cnt + 1
	}
	if flag != true {
		return fmt.Errorf("Label of specified name not found")
	}

	dataCon := con.S("labels").Index(cnt)
	d.SetId(models.StripQuotes(dataCon.S("id").String()))
	if dataCon.Exists("displayName") {
		d.Set("label", models.StripQuotes(dataCon.S("displayName").String()))
	}
	if dataCon.Exists("type") {
		d.Set("type", models.StripQuotes(dataCon.S("type").String()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
