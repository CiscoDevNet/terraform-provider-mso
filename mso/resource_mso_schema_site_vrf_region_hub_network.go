package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteVRFRegionHubNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteVRFRegionHubNetworkCreate,
		Delete: resourceMSOSchemaSiteVRFRegionHubNetworkDelete,
		Read:   resourceMSOSchemaSiteVRFRegionHubNetworkRead,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteVRFRegionHubNetworkImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"tenant_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"region_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"schema_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		},
	}
}

func setMSOSchemaSiteVrfRegionHubNetworkAttributes(hubNetwork *models.InterSchemaSiteVrfRegionHubNetork, d *schema.ResourceData) {
	d.Set("name", hubNetwork.Name)
	d.Set("tenant_name", hubNetwork.TenantName)
	d.Set("site_id", hubNetwork.SiteID)
	d.Set("template_name", hubNetwork.TemplateName)
	d.Set("vrf_name", hubNetwork.VrfName)
	d.Set("region_name", hubNetwork.Region)
	d.Set("schema_id", hubNetwork.SchemaID)
}

func resourceMSOSchemaSiteVRFRegionHubNetworkImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Beginning Import", d.Id())
	msoClient := m.(client.Client)
	id := d.Id()
	hubNetwork, err := hubNetworkIDtohubNetwork(id)
	if err != nil {
		return nil, err
	}
	hubNetworkMapRemote, err := msoClient.ReadInterSchemaSiteVrfRegionHubNetwork(hubNetwork)
	if err != nil {
		return nil, err
	}
	setMSOSchemaSiteVrfRegionHubNetworkAttributes(hubNetworkMapRemote, d)
	d.SetId(id)
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Import Completed", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteVRFRegionHubNetworkCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Beginning Creation")
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	name := d.Get("name").(string)
	regionName := d.Get("region_name").(string)
	tenantName := d.Get("tenant_name").(string)
	hubNetworkMap := models.InterSchemaSiteVrfRegionHubNetork{
		Name:         name,
		TenantName:   tenantName,
		Region:       regionName,
		VrfName:      vrfName,
		TemplateName: templateName,
		SiteID:       siteId,
		SchemaID:     schemaId,
	}
	err := msoClient.CreateInterSchemaSiteVrfRegionHubNetwork(&hubNetworkMap)
	if err != nil {
		return err
	}
	hubNetworkID := hubNetworkModeltohubNetworkID(&hubNetworkMap)
	d.SetId(hubNetworkID)
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Creation Complete")
	return resourceMSOSchemaSiteVRFRegionHubNetworkRead(d, m)
}

func resourceMSOSchemaSiteVRFRegionHubNetworkRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Beginning Read")
	msoClient := m.(*client.Client)
	id := d.Id()
	hubNetwork, err := hubNetworkIDtohubNetwork(id)
	if err != nil {
		return err
	}
	hubNetworkMapRemote, err := msoClient.ReadInterSchemaSiteVrfRegionHubNetwork(hubNetwork)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("%v", hubNetworkMapRemote)
	setMSOSchemaSiteVrfRegionHubNetworkAttributes(hubNetworkMapRemote, d)
	d.SetId(hubNetworkModeltohubNetworkID(hubNetwork))
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Read Complete", d.Id())
	return nil
}

func resourceMSOSchemaSiteVRFRegionHubNetworkDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Beginning Destroy")
	msoClient := m.(*client.Client)
	id := d.Id()
	hubNetwork, err := hubNetworkIDtohubNetwork(id)
	if err != nil {
		return err
	}
	err = msoClient.DeleteInterSchemaSiteVrfRegionHubNetwork(hubNetwork)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Schema Site VRF Region Hub Network: Destroy Completed", d.Id())
	d.SetId("")
	return err
}

func hubNetworkModeltohubNetworkID(m *models.InterSchemaSiteVrfRegionHubNetork) string {
	return fmt.Sprintf("%s/site/%s/template/%s/vrf/%s/region/%s/tenant/%s/%s", m.SchemaID, m.SiteID, m.TemplateName, m.VrfName, m.Region, m.TenantName, m.Name)
}

func hubNetworkIDtohubNetwork(ID string) (*models.InterSchemaSiteVrfRegionHubNetork, error) {
	getAttributes := strings.Split(ID, "/")
	if len(getAttributes) != 12 || getAttributes[1] != "site" || getAttributes[3] != "template" || getAttributes[5] != "vrf" || getAttributes[7] != "region" || getAttributes[9] != "tenant" {
		return nil, fmt.Errorf("invalid mso_schema_site_vrf_region_hub_network ID format")
	}
	hubNetwork := models.InterSchemaSiteVrfRegionHubNetork{
		SchemaID:     getAttributes[0],
		SiteID:       getAttributes[2],
		TemplateName: getAttributes[4],
		VrfName:      getAttributes[6],
		Region:       getAttributes[8],
		TenantName:   getAttributes[10],
		Name:         getAttributes[11],
	}
	return &hubNetwork, nil
}
