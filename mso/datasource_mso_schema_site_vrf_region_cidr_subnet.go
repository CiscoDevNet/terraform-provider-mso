package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteVrfRegionCidrSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteVrfRegionCidrSubnetRead,

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
			"cidr_ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteVrfRegionCidrSubnetRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrf := d.Get("vrf_name").(string)
	region := d.Get("region_name").(string)
	cidrIp := d.Get("cidr_ip").(string)
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

	cidrCont, err := getSiteVrfRegionCIDR(cidrIp, regionCont)
	if err != nil {
		return err
	} else {
		d.Set("cidr_ip", cidrIp)
	}

	subnetCount, err := cidrCont.ArrayCount("subnets")
	if err != nil {
		return fmt.Errorf("Unable to get Subnet list")
	}

	found := false
	for m := 0; m < subnetCount; m++ {
		subnetCont, err := cidrCont.ArrayElement(m, "subnets")
		if err != nil {
			return err
		}
		currentIp := models.StripQuotes(subnetCont.S("ip").String())
		if currentIp == ip {
			found = true
			d.SetId(fmt.Sprintf("%s/sites/%s-%s/vrfs/%s/regions/%s/cidr/%s/subnet/%s", schemaId, siteId, templateName, vrf, region, cidrIp, ip))
			d.Set("ip", ip)
			if subnetCont.Exists("zone") {
				d.Set("zone", models.StripQuotes(subnetCont.S("zone").String()))
			}
			if subnetCont.Exists("usage") {
				d.Set("usage", models.StripQuotes(subnetCont.S("usage").String()))
			}
			if subnetCont.Exists("subnetGroup") {
				d.Set("subnet_group", models.StripQuotes(subnetCont.S("subnetGroup").String()))
			}
			if subnetCont.Exists("name") {
				d.Set("name", models.StripQuotes(subnetCont.S("name").String()))
			}
			break
		}
	}

	if !found {
		d.SetId("")
		return fmt.Errorf("Unable to find VRF Region Cidr Subnet %s", ip)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
