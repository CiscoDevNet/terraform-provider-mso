package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteVrfRegionCidr() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteVrfRegionCidrRead,

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
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"region_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"primary": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteVrfRegionCidrRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrf := d.Get("vrf_name").(string)
	region := d.Get("region_name").(string)
	ip := d.Get("ip").(string)

	siteCont, err := getSiteFromSiteIdAndTemplate(schemaId, siteId, templateName, msoClient)
	if err != nil {
		return err
	} else {
		d.Set("schema_id", schemaId)
		d.Set("site_id", siteId)
		d.Set("template_name", templateName)
	}

	vrfCont, err := getSiteVrf(vrf, siteCont)
	if err != nil {
		return err
	} else {
		d.Set("vrf_name", vrf)
	}

	regionCont, err := getSiteVrfRegion(region, vrfCont)
	if err != nil {
		return err
	} else {
		d.Set("region_name", region)
	}

	cidrCont, err := getSiteVrfRegionCIDR(ip, regionCont)
	if err != nil {
		return err
	} else {
		d.SetId(fmt.Sprintf("%s/sites/%s-%s/vrfs/%s/regions/%s/cidr/%s", schemaId, siteId, templateName, vrf, region, ip))
		d.Set("ip", ip)
		d.Set("primary", cidrCont.S("primary").Data().(bool))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
