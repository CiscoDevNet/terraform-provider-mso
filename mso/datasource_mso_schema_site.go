package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
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
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func datasourceMSOSchemaSiteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	name := d.Get("name").(string)
	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	con, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/sites"))
	if err != nil {
		return err
	}

	var flag bool
	var siteId string
	for _, info := range con.S("sites").Data().([]interface{}) {
		val := info.(map[string]interface{})
		if val["name"].(string) == name {
			flag = true
			siteId = val["id"].(string)
			break
		}
	}
	if flag != true {
		return fmt.Errorf("Site of specified name not found")
	}
	_, err = getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/sites/%s-%s", schemaId, siteId, templateName))
	d.Set("schema_id", schemaId)
	d.Set("site_id", siteId)
	d.Set("template_name", templateName)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
