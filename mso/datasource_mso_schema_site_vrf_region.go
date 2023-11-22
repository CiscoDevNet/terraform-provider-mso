package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOSchemaSiteVrfRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteVrfRegionRead,

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
			"vpn_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"hub_network_enable": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"hub_network": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"subnet": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
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
								},
							},
						},
					},
				},
			},
		}),
	}
}

func dataSourceMSOSchemaSiteVrfRegionRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrf := d.Get("vrf_name").(string)
	region := d.Get("region_name").(string)

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
		d.SetId(fmt.Sprintf("%s/sites/%s-%s/vrfs/%s/regions/%s", schemaId, siteId, templateName, vrf, region))
		d.Set("region_name", region)
	}

	var cidrs []interface{}
	var regionOrVpcsContainer map[string]interface{}
	if regionCont.Exists("vpcs") {
		vpcsInterface := regionCont.S("vpcs").Data()
		vpcs := vpcsInterface.([]interface{})
		if len(vpcs) > 0 {
			regionOrVpcsContainer = vpcs[0].(map[string]interface{})
			if cidrsInterface, exists := regionOrVpcsContainer["cidrs"]; exists {
				cidrs = cidrsInterface.([]interface{})
			}
		}
	} else {
		regionOrVpcsContainer = regionCont.Data().(map[string]interface{})
		if regionCont.Exists("cidrs") {
			cidrsData := regionCont.S("cidrs").Data()
			if cidrsData != nil {
				cidrs = cidrsData.([]interface{})
			}
		}
	}

	if isVpnGatewayRouter, exists := regionOrVpcsContainer["isVpnGatewayRouter"]; exists {
		d.Set("vpn_gateway", isVpnGatewayRouter)
	}
	if isTGWAttachment, exists := regionOrVpcsContainer["isTGWAttachment"]; exists {
		d.Set("hub_network_enable", isTGWAttachment)
	}
	hubMap := make(map[string]interface{})
	if cloudRsCtxProfileToGatewayRouterP, exists := regionOrVpcsContainer["cloudRsCtxProfileToGatewayRouterP"]; exists {
		temp := cloudRsCtxProfileToGatewayRouterP.(map[string]interface{})
		hubMap["name"] = temp["name"]
		hubMap["tenant_name"] = temp["tenantName"]
		d.Set("hub_network", hubMap)
	}

	cidrList := make([]interface{}, 0, 1)
	for _, tempCidr := range cidrs {
		cidr := tempCidr.(map[string]interface{})

		cidrMap := make(map[string]interface{})
		cidrMap["cidr_ip"] = cidr["ip"]
		cidrMap["primary"] = cidr["primary"]

		subnets := cidr["subnets"].([]interface{})
		subnetList := make([]interface{}, 0, 1)
		for _, tempSubnet := range subnets {
			subnet := tempSubnet.(map[string]interface{})

			subnetMap := make(map[string]interface{})
			subnetMap["ip"] = subnet["ip"]
			if subnet["name"] != nil {
				subnetMap["name"] = subnet["name"]
			}
			if subnet["zone"] != nil {
				subnetMap["zone"] = subnet["zone"]
			}
			if subnet["usage"] != nil {
				subnetMap["usage"] = subnet["usage"]
			}
			if subnet["subnetGroup"] != nil {
				subnetMap["subnet_group"] = subnet["subnetGroup"]
			}

			subnetList = append(subnetList, subnetMap)
		}
		cidrMap["subnet"] = subnetList

		cidrList = append(cidrList, cidrMap)
	}
	d.Set("cidr", cidrList)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
