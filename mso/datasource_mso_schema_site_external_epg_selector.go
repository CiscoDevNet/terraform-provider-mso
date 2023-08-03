package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteExternalEpgSelector() *schema.Resource {
	return &schema.Resource{
		Read: datasourceMSOSchemaSiteExternalEpgSelectorRead,

		Schema: map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"external_epg_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceMSOSchemaSiteExternalEpgSelectorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Data source Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	externalEpgName := d.Get("external_epg_name").(string)
	subnetName := d.Get("name").(string)

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
		d.Set("external_epg_name", externalEpgName)
	}

	subnetCount, err := externalEpgCont.ArrayCount("subnets")
	if err != nil {
		return err
	}

	found := false
	for k := 0; k < subnetCount; k++ {
		subnetCont, err := externalEpgCont.ArrayElement(k, "subnets")
		if err != nil {
			return err
		}
		currentSubnetName := models.StripQuotes(subnetCont.S("name").String())
		if subnetName == currentSubnetName {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/externalEpgs/%s/subnets/%s", schemaId, siteId, templateName, externalEpgName, subnetName))
			d.Set("name", subnetName)
			d.Set("ip", models.StripQuotes(subnetCont.S("ip").String()))
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find External EPG Subnet: %s", subnetName)
	}

	log.Printf("[DEBUG] %s: Datasource Read finished successfully", d.Id())
	return nil
}
