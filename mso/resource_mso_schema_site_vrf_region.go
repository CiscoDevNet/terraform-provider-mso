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

func resourceMSOSchemaSiteVrfRegion() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteVrfRegionCreate,
		Update: resourceMSOSchemaSiteVrfRegionUpdate,
		Read:   resourceMSOSchemaSiteVrfRegionRead,
		Delete: resourceMSOSchemaSiteVrfRegionDelete,

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
				Required:     true,
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
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"tenant_name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
					},
				},
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_ip": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1000),
						},
						"primary": &schema.Schema{
							Type:     schema.TypeBool,
							Required: true,
						},
						"subnet": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 1000),
									},
									"zone": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
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

func resourceMSOSchemaSiteVrfRegionCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf Region: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	regionName := d.Get("region_name").(string)

	var vpnGateway bool
	if vpn, ok := d.GetOk("vpn_gateway"); ok {
		vpnGateway = vpn.(bool)
	}

	var hubEnable bool
	if hub, ok := d.GetOk("hub_network_enable"); ok {
		hubEnable = hub.(bool)
	}

	hubNetworkMap := make(map[string]interface{})
	if hubEnable {
		if tp, ok := d.GetOk("hub_network"); ok {
			hubNetwork := tp.(map[string]interface{})

			log.Println("check  hub Network ... :", hubNetwork["name"], hubNetwork["tenant_name"])

			if hubNetwork["name"] != nil && hubNetwork["tenant_name"] != nil {
				hubNetworkMap["name"] = hubNetwork["name"]
				hubNetworkMap["tenantName"] = hubNetwork["tenant_name"]
			} else {
				return fmt.Errorf("missing attribute in hub_network. Either name or tenant_name missing")
			}

		} else {
			return fmt.Errorf("hub_network field is missing.")
		}
	}

	cidrs := d.Get("cidr").([]interface{})
	cidrsList := make([]interface{}, 0, 1)
	for _, tempCidr := range cidrs {
		cidr := tempCidr.(map[string]interface{})

		cidrMap := make(map[string]interface{})
		cidrMap["ip"] = cidr["cidr_ip"].(string)
		cidrMap["primary"] = cidr["primary"].(bool)

		subnets := cidr["subnet"].([]interface{})
		subnetList := make([]interface{}, 0, 10)
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

		cidrMap["subnets"] = subnetList

		cidrsList = append(cidrsList, cidrMap)
	}

	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/-", siteId, templateName, vrfName)
	vrfRegionStruct := models.NewSchemaSiteVrfRegion("add", path, regionName, vpnGateway, hubEnable, hubNetworkMap, cidrsList)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), vrfRegionStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteVrfRegionRead(d, m)
}

func resourceMSOSchemaSiteVrfRegionUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf Region: Beginning Updation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	regionName := d.Get("region_name").(string)

	var vpnGateway bool
	if vpn, ok := d.GetOk("vpn_gateway"); ok {
		vpnGateway = vpn.(bool)
	}

	var hubEnable bool
	if hub, ok := d.GetOk("hub_network_enable"); ok {
		hubEnable = hub.(bool)
	}

	hubNetworkMap := make(map[string]interface{})
	if hubEnable {
		if tp, ok := d.GetOk("hub_network"); ok {
			hubNetwork := tp.(map[string]interface{})

			log.Println("check  hub Network ... :", hubNetwork["name"], hubNetwork["tenant_name"])

			if hubNetwork["name"] != nil && hubNetwork["tenant_name"] != nil {
				hubNetworkMap["name"] = hubNetwork["name"]
				hubNetworkMap["tenantName"] = hubNetwork["tenant_name"]
			} else {
				return fmt.Errorf("missing attribute in hub_network. Either name or tenant_name missing")
			}

		} else {
			return fmt.Errorf("hub_network field is missing.")
		}
	}

	cidrs := d.Get("cidr").([]interface{})
	cidrsList := make([]interface{}, 0, 1)
	for _, tempCidr := range cidrs {
		cidr := tempCidr.(map[string]interface{})

		cidrMap := make(map[string]interface{})
		cidrMap["ip"] = cidr["cidr_ip"].(string)
		cidrMap["primary"] = cidr["primary"].(bool)

		subnets := cidr["subnet"].([]interface{})
		subnetList := make([]interface{}, 0, 10)
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

		cidrMap["subnets"] = subnetList

		cidrsList = append(cidrsList, cidrMap)
	}

	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s", siteId, templateName, vrfName, regionName)
	vrfRegionStruct := models.NewSchemaSiteVrfRegion("replace", path, regionName, vpnGateway, hubEnable, hubNetworkMap, cidrsList)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), vrfRegionStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteVrfRegionRead(d, m)
}

func resourceMSOSchemaSiteVrfRegionRead(d *schema.ResourceData, m interface{}) error {
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
		d.SetId("")
		d.Set("schema_id", "")
		d.Set("site_id", "")
		d.Set("template_name", "")
		d.Set("region_name", "")
		d.Set("vrf_name", "")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteVrfRegionDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf Region: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	regionName := d.Get("region_name").(string)

	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s", siteId, templateName, vrfName, regionName)
	vrfRegionStruct := models.NewSchemaSiteVrfRegion("remove", path, regionName, false, false, nil, nil)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), vrfRegionStruct)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
