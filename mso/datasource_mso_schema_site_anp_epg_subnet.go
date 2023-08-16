package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteAnpEpgSubnet() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaSiteAnpEpgSubnetRead,

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
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"no_default_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"querier": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"primary": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		}),
	}
}

func datasourceMSOSchemaSiteAnpEpgSubnetRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anp := d.Get("anp_name").(string)
	epg := d.Get("epg_name").(string)
	ip := d.Get("ip").(string)

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

	subnetCount, err := epgCont.ArrayCount("subnets")
	if err != nil {
		return fmt.Errorf("Unable to get Subnet list")
	}

	found := false
	for l := 0; l < subnetCount; l++ {
		subnetCont, err := epgCont.ArrayElement(l, "subnets")
		if err != nil {
			return err
		}
		currentIp := models.StripQuotes(subnetCont.S("ip").String())
		if ip == currentIp {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/anps/%s/epgs/%s/subnets/%s", schemaId, siteId, templateName, anp, epg, ip))
			if subnetCont.Exists("ip") {
				d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
			}
			if subnetCont.Exists("description") {
				d.Set("description", models.StripQuotes(subnetCont.S("description").String()))
			}
			if subnetCont.Exists("scope") {
				d.Set("scope", models.StripQuotes(subnetCont.S("scope").String()))
			}
			if subnetCont.Exists("shared") {
				d.Set("shared", subnetCont.S("shared").Data().(bool))
			}
			if subnetCont.Exists("noDefaultGateway") {
				d.Set("no_default_gateway", subnetCont.S("noDefaultGateway").Data().(bool))
			}
			if subnetCont.Exists("primary") {
				d.Set("primary", subnetCont.S("primary").Data().(bool))
			}
			if subnetCont.Exists("querier") {
				d.Set("querier", subnetCont.S("querier").Data().(bool))
			}
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find subnet entry with ip: %s", ip)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
