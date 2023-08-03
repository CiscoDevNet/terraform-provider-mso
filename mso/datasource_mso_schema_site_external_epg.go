package mso

import (
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteExternalEpg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteExternalEpgRead,

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
			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"l3out_dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l3out_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l3out_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l3out_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteExternalEpgRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("schema_id", schemaId)
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	externalEpgCont, err := getSiteExternalEpg(externalEpgName, siteCont)
	if err != nil {
		return err
	} else {
		d.SetId(fmt.Sprintf("%s/sites/%s-%s/externalEpgs/%s", schemaId, siteId, templateName, externalEpgName))
		d.Set("external_epg_name", externalEpgName)
	}

	l3outRef := models.StripQuotes(externalEpgCont.S("l3outRef").String())
	if l3outRef != "{}" && l3outRef != "" {
		re := regexp.MustCompile("/schemas/(.*?)/templates/(.*?)/l3outs/(.*)")
		currentL3out := re.FindStringSubmatch(l3outRef)
		if len(currentL3out) >= 4 {
			d.Set("l3out_name", currentL3out[3])
			d.Set("l3out_schema_id", currentL3out[1])
			d.Set("l3out_template_name", currentL3out[2])
		} else {
			return fmt.Errorf("Error in parsing l3outRef to get L3Out name")
		}
	}

	if externalEpgCont.Exists("l3outDn") {
		d.Set("l3out_dn", models.StripQuotes(externalEpgCont.S("l3outDn").String()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
