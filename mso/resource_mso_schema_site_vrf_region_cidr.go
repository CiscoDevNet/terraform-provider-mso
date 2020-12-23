package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteVrfRegionCidr() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteVrfRegionCidrCreate,
		Update: resourceMSOSchemaSiteVrfRegionCidrUpdate,
		Read:   resourceMSOSchemaSiteVrfRegionCidrRead,
		Delete: resourceMSOSchemaSiteVrfRegionCidrDelete,

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
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"primary": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		}),
	}
}

func resourceMSOSchemaSiteVrfRegionCidrCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf Region Cidr: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	regionName := d.Get("region_name").(string)
	ip := d.Get("ip").(string)
	primary := d.Get("primary").(bool)

	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s/cidrs/-", siteId, templateName, vrfName, regionName)
	VrfRegionCidrStruct := models.NewSchemaSiteVrfRegionCidr("add", path, ip, primary)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), VrfRegionCidrStruct)
	if err != nil {
		return err
	}
	return resourceMSOSchemaSiteVrfRegionCidrRead(d, m)
}

func resourceMSOSchemaSiteVrfRegionCidrRead(d *schema.ResourceData, m interface{}) error {
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
		d.SetId("")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaSiteVrfRegionCidrUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf Region Cidr: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	regionName := d.Get("region_name").(string)
	ip := d.Get("ip").(string)
	primary := d.Get("primary").(bool)

	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	index, err := indexCont(cont, siteId, vrfName, regionName, id)
	if err != nil {
		return err
	}
	if index == -1 {
		return fmt.Errorf("The given Vrf Region Cidr is not found")
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s/cidrs/%s", siteId, templateName, vrfName, regionName, indexs)
	VrfRegionCidrStruct := models.NewSchemaSiteVrfRegionCidr("replace", path, ip, primary)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), VrfRegionCidrStruct)
	if errs != nil {
		return errs
	}
	return resourceMSOSchemaSiteVrfRegionCidrRead(d, m)
}

func resourceMSOSchemaSiteVrfRegionCidrDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Site Vrf Region Cidr: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	siteId := d.Get("site_id").(string)
	templateName := d.Get("template_name").(string)
	vrfName := d.Get("vrf_name").(string)
	regionName := d.Get("region_name").(string)
	ip := d.Get("ip").(string)
	primary := d.Get("primary").(bool)

	id := d.Id()
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	index, err := indexCont(cont, siteId, vrfName, regionName, id)
	if err != nil {
		return err
	}
	if index == -1 {
		d.SetId("")
		return nil
	}
	indexs := strconv.Itoa(index)

	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s/regions/%s/cidrs/%s", siteId, templateName, vrfName, regionName, indexs)
	VrfRegionCidrStruct := models.NewSchemaSiteVrfRegionCidr("remove", path, ip, primary)

	_, errs := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), VrfRegionCidrStruct)
	if errs != nil {
		return errs
	}
	d.SetId("")
	return nil
}

func indexCont(cont *container.Container, stateSite, stateVrf, stateRegion, stateIp string) (int, error) {
	found := false
	index := -1
	count, err := cont.ArrayCount("sites")
	if err != nil {
		return index, fmt.Errorf("No Sites found")
	}

	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "sites")
		if err != nil {
			return index, err
		}
		apiSite := models.StripQuotes(tempCont.S("siteId").String())

		if apiSite == stateSite {
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return index, fmt.Errorf("Unable to get Vrf list")
			}
			for j := 0; j < vrfCount; j++ {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return index, err
				}
				apiVrfRef := models.StripQuotes(vrfCont.S("vrfRef").String())
				split := strings.Split(apiVrfRef, "/")
				apiVrf := split[6]
				if apiVrf == stateVrf {
					regionCount, err := vrfCont.ArrayCount("regions")
					if err != nil {
						return index, fmt.Errorf("Unable to get Regions list")
					}
					for k := 0; k < regionCount; k++ {
						regionCont, err := vrfCont.ArrayElement(k, "regions")
						if err != nil {
							return index, err
						}
						apiRegion := models.StripQuotes(regionCont.S("name").String())
						if apiRegion == stateRegion {
							cidrCount, err := regionCont.ArrayCount("cidrs")
							if err != nil {
								return index, fmt.Errorf("Unable to get Cidr list")
							}
							for l := 0; l < cidrCount; l++ {
								cidrCont, err := regionCont.ArrayElement(l, "cidrs")
								if err != nil {
									return index, err
								}
								apiIp := models.StripQuotes(cidrCont.S("ip").String())
								if apiIp == stateIp {
									log.Println("found correct Cidr")
									index = l
									found = true
									break
								}
							}
						}
						if found {
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
		if found {
			break
		}
	}
	return index, nil

}
