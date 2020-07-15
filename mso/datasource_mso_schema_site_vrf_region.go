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

func dataSourceMSOSchemaSiteVrfRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSchemaSiteVrfRegionRead,

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
			"vpn_gateway": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"hub_network_enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"hub_network": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"tenant_name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_ip": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"primary": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"subnet": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
									},
									"zone": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
									},
									"usage": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
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
							d.SetId(apiRegion)
							d.Set("region_name", apiRegion)
							if regionCont.Exists("isVpnGatewayRouter") {
								d.Set("vpn_gateway", regionCont.S("isVpnGatewayRouter").Data().(bool))
							}
							if regionCont.Exists("isTGWAttachment") {
								d.Set("hub_network_enable", regionCont.S("isTGWAttachment").Data().(bool))
							}

							hubMap := make(map[string]interface{})
							if regionCont.Exists("cloudRsCtxProfileToGatewayRouterP") {
								temp := regionCont.S("cloudRsCtxProfileToGatewayRouterP").Data().(map[string]interface{})

								hubMap["name"] = temp["name"]
								hubMap["tenant_name"] = temp["tenantName"]

								d.Set("hub_network", hubMap)
							} else {
								d.Set("hub_network", hubMap)
							}

							cidrList := make([]interface{}, 0, 1)
							cidrs := regionCont.S("cidrs").Data().([]interface{})
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
									subnetMap["zone"] = subnet["zone"]
									if subnet["usage"] != nil {
										subnetMap["usage"] = subnet["usage"]
									}

									subnetList = append(subnetList, subnetMap)
								}
								cidrMap["subnet"] = subnetList

								cidrList = append(cidrList, cidrMap)
							}
							d.Set("cidr", cidrList)
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("Unable to find the Site Vrf Region %s", stateRegion)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}
