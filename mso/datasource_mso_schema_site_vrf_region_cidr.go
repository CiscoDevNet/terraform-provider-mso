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

func dataSourceMSOSchemaSiteVrfRegionCidr() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceMSOSchemaSiteVrfRegionCidrRead,

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
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"primary": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		}),
	}
}

func dataSourceMSOSchemaSiteVrfRegionCidrRead(d *schema.ResourceData, m interface{}) error {
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
	stateIp := d.Get("ip").(string)

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
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
					d.Set("site_id", apiSite)
					d.Set("schema_id", split[2])
					d.Set("template_name", split[4])
					d.Set("vrf_name", split[6])
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
								apiIp := models.StripQuotes(cidrCont.S("ip").String())
								if apiIp == stateIp {
									d.SetId(apiIp)
									d.Set("ip", apiIp)
									d.Set("primary", cidrCont.S("primary").Data().(bool))
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

	if !found {
		return fmt.Errorf("Unable to find the given Vrf Region Cidr")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
