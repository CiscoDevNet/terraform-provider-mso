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

func dataSourceMSOSchemaSiteVrfRegionCidrSubnet() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteVrfRegionCidrSubnetRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"site_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"region_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"cidr_ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"usage": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
		}),
	}
}

func dataSourceMSOSchemaSiteVrfRegionCidrSubnetRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return fmt.Errorf("No Sites found")
	}

	stateSite := d.Get("site_id").(string)
	found := false
	stateVrf := d.Get("vrf_name").(string)
	stateRegion := d.Get("region_name").(string)
	stateCidr := d.Get("cidr_ip").(string)
	stateIp := d.Get("ip").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			apiTemplate := models.StripQuotes(tempCont.S("templateName").String())
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return fmt.Errorf("Unable to get Vrf list")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return err
				}
				apiVrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
				split := strings.Split(apiVrfRef, "/")
				apiVrf := split[6]
				if apiVrf == stateVrf {
					regionCount, err := vrfCont.ArrayCount("regions")
					if err != nil {
						return fmt.Errorf("Unable to get Regions list")
					}
					for k := 0; k < regionCount; k++ {
						regionCont, err := vrfCont.ArrayElement(k, "regions")
						if err != nil {
							return err
						}
						apiRegion := models.StripQuotes(regionCont.S("name").String())
						if apiRegion == stateRegion {
							cidrCount, err := regionCont.ArrayCount("cidrs")
							if err != nil {
								return fmt.Errorf("Unable to get Cidr list")
							}
							for l := 0; l < cidrCount; l++ {
								cidrCont, err := regionCont.ArrayElement(l, "cidrs")
								if err != nil {
									return err
								}
								apiCidr := models.StripQuotes(cidrCont.S("ip").String())
								log.Println("Current Cidr Ip", apiCidr)
								if apiCidr == stateCidr {
									subnetCount, err := cidrCont.ArrayCount("subnets")
									if err != nil {
										return fmt.Errorf("Unable to get Subnet list")
									}
									for m := 0; m < subnetCount; m++ {
										subnetCont, err := cidrCont.ArrayElement(m, "subnets")
										if err != nil {
											return err
										}
										apiIp := models.StripQuotes(subnetCont.S("ip").String())
										if apiIp == stateIp {
											d.SetId(apiIp)
											d.Set("ip", apiIp)
											d.Set("site_id", apiSite)
											d.Set("template_name", apiTemplate)
											d.Set("cidr_name", apiCidr)
											d.Set("vrf_name", apiVrf)
											d.Set("region_name", apiRegion)
											if subnetCont.Exists("zone") {
												d.Set("zone", models.StripQuotes(subnetCont.S("zone").String()))
											}
											if subnetCont.Exists("usage") {
												d.Set("usage", models.StripQuotes(subnetCont.S("usage").String()))
											}
											found = true
											break
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if !found {
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
		d.Set("region_name", "")
		d.Set("vrf_name", "")
		return fmt.Errorf("Unable to find VRF Region Cidr Subnet %s", stateIp)

	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
