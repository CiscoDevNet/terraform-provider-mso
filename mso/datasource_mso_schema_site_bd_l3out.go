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
			"l3out_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"l3out_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
	l3outName := d.Get("l3out_name").(string)
	l3outSchemaId := d.Get("l3out_schema_id").(string)
	if l3outSchemaId == "" {
		l3outSchemaId = schemaId
	}
	l3outTemplateName := d.Get("l3out_template_name").(string)
	if l3outTemplateName == "" {
		l3outTemplateName = templateName
	}

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

	found := false
	l3outCount, err := bdCont.ArrayCount("l3Outs")
	if err == nil {
		for k := range l3outCount {
			l3outCont, err := bdCont.ArrayElement(k, "l3Outs")
			if err != nil {
				continue
			}
			if strings.Trim(l3outCont.String(), "\"") == l3outName {
				found = true
				d.SetId(fmt.Sprintf("%s/sites/%s-%s/bds/%s/l3outs/%s-%s-%s", schemaId, siteId, templateName, bd, l3outSchemaId, l3outTemplateName, l3outName))
				d.Set("l3out_name", l3outName)
				d.Set("l3out_template_name", l3outTemplateName)
				d.Set("l3out_schema_id", l3outSchemaId)
				break
			}
		}
	}

	// Fallback for older NDO versions that used full reference paths in l3OutRefs.
	if !found && bdCont.Exists("l3OutRefs") {
		l3OutRefs := bdCont.S("l3OutRefs").Data()
		for _, l3OutRef := range l3OutRefs.([]interface{}) {
			splitl3OutRef := strings.Split(l3OutRef.(string), "/")
			if splitl3OutRef[6] == l3outName && splitl3OutRef[4] == l3outTemplateName && splitl3OutRef[2] == l3outSchemaId {
				found = true
				d.SetId(fmt.Sprintf("%s/sites/%s-%s/bds/%s/l3outs/%s-%s-%s", schemaId, siteId, templateName, bd, l3outSchemaId, l3outTemplateName, l3outName))
				d.Set("l3out_name", splitl3OutRef[6])
				d.Set("l3out_template_name", splitl3OutRef[4])
				d.Set("l3out_schema_id", splitl3OutRef[2])
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the Site BD L3out %s in Template %s of Schema Id %s ", l3outName, l3outTemplateName, l3outSchemaId)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
