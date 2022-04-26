package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteAnpEpgUsegAttr() *schema.Resource {
	return &schema.Resource{
		Read:          datasourceMSOSchemaSiteAnpEpgUsegAttrRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"schema_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"anp_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"epg_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"useg_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"useg_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fv_subnet": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func datasourceMSOSchemaSiteAnpEpgUsegAttrRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Datasource Schema Site Anp Epg Useg Attr: Beginning Read")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	epgName := d.Get("epg_name").(string)
	usegName := d.Get("useg_name").(string)

	usegMap := models.SiteUsegAttr{
		SiteID:       siteId,
		TemplateName: templateName,
		SchemaID:     schemaId,
		AnpName:      anpName,
		EpgName:      epgName,
		UsegName:     usegName,
	}
	usegMapRemote, _, err := msoClient.ReadAnpEpgUsegAttr(&usegMap)
	if err != nil {
		return err
	}
	setMSOSchemaSiteAnpEpgUsegAttributes(usegMapRemote, d)
	id := UsegAttrModelToUsegId(usegMapRemote)
	d.SetId(id)
	log.Println("[DEBUG] Datasource Schema Site Anp Epg Useg Attr: Reading Completed", d.Id())
	return nil
}
