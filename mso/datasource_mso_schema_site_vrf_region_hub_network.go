package mso

import (
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaSiteVRFRegionHubNetwork() *schema.Resource {
	return &schema.Resource{
		Read:          datasourceMSOSchemaSiteVRFRegionHubNetworkRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"tenant_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"region_name": {
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

func datasourceMSOSchemaSiteVRFRegionHubNetworkRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Beginning Read")
	msoClient := m.(*client.Client)
	schemaID := d.Get("schema_id").(string)
	siteID := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	regionName := d.Get("region_name").(string)
	name := d.Get("name").(string)
	tenantName := d.Get("tenant_name").(string)
	hubNetworkMap := models.InterSchemaSiteVrfRegionHubNetork{
		SchemaID:     schemaID,
		SiteID:       siteID,
		TemplateName: templateName,
		VrfName:      vrfName,
		Region:       regionName,
		Name:         name,
		TenantName:   tenantName,
	}
	hubNetworkMapRemote, err := msoClient.ReadInterSchemaSiteVrfRegionHubNetwork(&hubNetworkMap)
	if err != nil {
		d.SetId("")
		return err
	}
	setMSOSchemaSiteVrfRegionHubNetworkAttributes(hubNetworkMapRemote, d)
	d.SetId(hubNetworkModeltohubNetworkID(&hubNetworkMap))
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Reading Completed", d.Id())
	return nil
}
