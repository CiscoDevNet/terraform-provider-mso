package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteL3out() *schema.Resource {
	return &schema.Resource{
		Read:          datasourceMSOSchemaSiteL3outRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"l3out_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
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
		},
	}
}

func datasourceMSOSchemaSiteL3outRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site L3out: Beginning Read")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	l3outName := d.Get("l3out_name").(string)
	l3outMap := models.IntersiteL3outs{
		SchemaID:     schemaId,
		SiteId:       siteId,
		TemplateName: templateName,
		VRFName:      vrfName,
		L3outName:    l3outName,
	}
	l3outMapRemote, err := msoClient.ReadIntersiteL3outs(&l3outMap)
	if err != nil {
		d.SetId("")
		return err
	}
	setMSOSchemaSiteL3outAttributes(l3outMapRemote, d)
	d.SetId(L3outModelToL3outId(&l3outMap))
	log.Println("[DEBUG] Schema Site L3out: Reading Completed", d.Id())
	return nil
}
