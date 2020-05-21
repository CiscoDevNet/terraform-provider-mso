package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSite() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaSiteRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},

			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func datasourceMSOSchemaSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	schemaId := d.Get("schema_id").(string)
	con, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/sites"))
	if err != nil {
		return err
	}

	data := con.S("sites").Data().([]interface{})
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
		return fmt.Errorf("Site of specified name not found")
	}

	dataCon := con.S("sites").Index(count)
	stateSiteId := models.StripQuotes(dataCon.S("id").String())

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	countSites, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Template found")
	}

	found := false

	for i := 0; i < countSites; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSiteId := models.StripQuotes(tempCont.S("siteId").String())
		apiTemplate := models.StripQuotes(tempCont.S("templateName").String())

		if apiSiteId == stateSiteId {
			d.SetId(apiSiteId)
			d.Set("schema_id", schemaId)
			d.Set("site_id", apiSiteId)
			d.Set("template_name", apiTemplate)
			found = true
		}

	}

	if !found {
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
