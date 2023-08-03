package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteBdL3out() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteBdL3outRead,

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
			"l3out_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func dataSourceMSOSchemaSiteBdL3outRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	bd := d.Get("bd_name").(string)
	l3out := d.Get("l3out_name").(string)

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
		d.Set("bd_name", bd)
	}

	l3outCount, err := bdCont.ArrayCount("l3Outs")
	if err != nil {
		return fmt.Errorf("Unable to get l3Outs list")
	}

	found := false
	for k := 0; k < l3outCount; k++ {
		l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
		if err != nil {
			return err
		}
		currentL3out := strings.Trim(l3outCont.String(), "\"")
		if currentL3out == l3out {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/bds/%s/l3outs/%s", schemaId, siteId, templateName, bd, currentL3out))
			d.Set("l3out_name", currentL3out)
			break
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the Site BD L3out: %s", l3out)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
