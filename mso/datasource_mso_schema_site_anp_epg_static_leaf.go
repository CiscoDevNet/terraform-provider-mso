package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteAnpEpgStaticleaf() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteAnpEpgStaticleafRead,

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
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"port_encap_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteAnpEpgStaticleafRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	path := d.Get("path").(string)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("schema_id", schemaId)
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	anpCont, err := getSiteAnp(anp, siteCont)
	if err != nil {
		return err
	} else {
		d.Set("anp_name", anp)
	}

	epgCont, err := getSiteEpg(epg, anpCont)
	if err != nil {
		return err
	} else {
		d.Set("epg_name", epg)
	}

	staticLeafCount, err := epgCont.ArrayCount("staticLeafs")
	if err != nil {
		return fmt.Errorf("Unable to get Static Leaf list")
	}

	found := false
	for s := 0; s < staticLeafCount; s++ {
		staticLeafCont, err := epgCont.ArrayElement(s, "staticLeafs")
		if err != nil {
			return err
		}
		currentPath := models.StripQuotes(staticLeafCont.S("path").String())
		if currentPath == path {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/anps/%s/epgs/%s/staticLeafs/%s", schemaId, siteId, templateName, anp, epg, path))
			d.Set("path", path)
			port, _ := strconv.Atoi(staticLeafCont.S("portEncapVlan").String())
			d.Set("port_encap_vlan", port)
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find the Site ANP EPG Static Leaf: %s", path)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
