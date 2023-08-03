package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteBd() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteBdRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"host_route": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"svi_mac": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteBdRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bd := d.Get("bd_name").(string)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("schema_id", schemaId)
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	bdCont, err := getSiteBd(bd, siteCont)
	if err != nil {
		return err
	} else {
		d.SetId(fmt.Sprintf("%s/sites/%s-%s/bds/%s", schemaId, siteId, templateName, bd))
		d.Set("bd_name", bd)
	}

	if bdCont.Exists("hostBasedRouting") {
		d.Set("host_route", bdCont.S("hostBasedRouting").Data().(bool))
	}
	if bdCont.Exists("mac") {
		d.Set("svi_mac", models.StripQuotes(bdCont.S("mac").String()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
