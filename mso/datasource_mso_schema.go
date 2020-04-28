package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func datasourceMSOSchema() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaRead,

		SchemaVersion: 1,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"tenant_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
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
		return fmt.Errorf("Schema of specified name not found")
	}

	dataCon := con.S("schemas").Index(cnt)
	d.SetId(models.StripQuotes(dataCon.S("id").String()))
	d.Set("name", models.StripQuotes(dataCon.S("displayName").String()))
	count, err := dataCon.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	found := false
	for i := 0; i < count; i++ {
		tempCont, err := dataCon.ArrayElement(i, "templates")

		if err != nil {
			return fmt.Errorf("Unable to parse the template list")
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())
		apiTenant := models.StripQuotes(tempCont.S("tenantId").String())

		d.Set("template_name", apiTemplate)
		d.Set("tenant_id", apiTenant)
		found = true
		break
	}

	if !found {
		d.Set("template_name", "")
		d.Set("tenant_id", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
