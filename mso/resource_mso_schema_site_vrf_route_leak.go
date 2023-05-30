package mso

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaSiteVrfRouteLeak() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaSiteVrfRouteLeakCreate,
		Update: resourceMSOSchemaSiteVrfRouteLeakUpdate,
		Read:   resourceMSOSchemaSiteVrfRouteLeakRead,
		Delete: resourceMSOSchemaSiteVrfRouteLeakDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaSiteVrfRouteLeakImport,
		},

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
			"target_vrf_schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"target_vrf_template_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"target_vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"tenant_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "leak_all",
				ValidateFunc: validation.StringInSlice([]string{
					"leak_all",
					"subnet_ip",
					"all_subnet_ips",
				}, false),
			},
			"subnet_ips": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func getVrfRef(schemaId, TemplateName, VrfName string) string {
	return fmt.Sprintf("/schemas/%s/templates/%s/vrfs/%s", schemaId, TemplateName, VrfName)
}

func getTargetVrfRef(d *schema.ResourceData) string {
	targetVrfSchemaId := d.Get("target_vrf_schema_id").(string)
	if targetVrfSchemaId == "" {
		targetVrfSchemaId = d.Get("schema_id").(string)
	}
	targetVrfTemplateName := d.Get("target_vrf_template_name").(string)
	if targetVrfTemplateName == "" {
		targetVrfTemplateName = d.Get("template_name").(string)
	}
	targetVrfName := d.Get("target_vrf_name").(string)
	return getVrfRef(targetVrfSchemaId, targetVrfTemplateName, targetVrfName)
}

func getSubnetDetails(d *schema.ResourceData) ([]map[string]string, bool) {
	var prefixSubnets []map[string]string
	var includeAllSubnets bool
	routeLeakType := d.Get("type").(string)
	if routeLeakType == "all_subnet_ips" {
		includeAllSubnets = true
	} else if routeLeakType == "leak_all" {
		prefixSubnet := map[string]string{"ip": "0.0.0.0/0"}
		prefixSubnets = append(prefixSubnets, prefixSubnet)
	} else {
		subnetIps := d.Get("subnet_ips").(*schema.Set).List()
		for _, subnetIp := range subnetIps {
			prefixSubnet := map[string]string{"ip": subnetIp.(string)}
			prefixSubnets = append(prefixSubnets, prefixSubnet)
		}
	}
	return prefixSubnets, includeAllSubnets
}

func setRouteLeakFromSchema(d *schema.ResourceData, schemaCont *container.Container, schemaId, siteId, templateName, vrfName, targetVrfRef string) error {
	sites := schemaCont.Search("sites").Data()
	if sites == nil || len(sites.([]interface{})) == 0 {
		return fmt.Errorf("no sites found")
	}
	for _, site := range sites.([]interface{}) {
		siteDetails := site.(map[string]interface{})
		if siteDetails["siteId"].(string) == siteId && siteDetails["templateName"].(string) == templateName {
			for _, vrf := range siteDetails["vrfs"].(interface{}).([]interface{}) {
				vrfDetails := vrf.(map[string]interface{})
				if strings.Split(vrfDetails["vrfRef"].(string), "/")[6] == vrfName {
					for index, routeLeak := range vrfDetails["routeLeak"].(interface{}).([]interface{}) {
						routeLeakDetails := routeLeak.(map[string]interface{})
						if routeLeakDetails["vrfRef"].(string) == targetVrfRef {
							d.SetId(fmt.Sprintf("/sites/%s-%s/vrfs/%s/routeLeak/%v", siteId, templateName, vrfName, index))
							d.Set("schema_id", schemaId)
							d.Set("site_id", siteId)
							d.Set("template_name", templateName)
							d.Set("vrf_name", vrfName)
							d.Set("target_vrf_schema_id", strings.Split(targetVrfRef, "/")[2])
							d.Set("target_vrf_template_name", strings.Split(targetVrfRef, "/")[4])
							d.Set("target_vrf_name", strings.Split(targetVrfRef, "/")[6])
							d.Set("tenant_name", routeLeakDetails["tenantName"].(string))
							subnetIps := []string{}
							for _, subnetIp := range routeLeakDetails["prefixsubnet"].(interface{}).([]interface{}) {
								subnetIp := subnetIp.(map[string]interface{})["ip"].(string)
								subnetIps = append(subnetIps, subnetIp)
							}
							if routeLeakDetails["includeAllSubnets"].(bool) == true {
								d.Set("type", "all_subnet_ips")
							} else if len(subnetIps) == 1 && subnetIps[0] == "0.0.0.0/0" {
								d.Set("type", "leak_all")
							} else {
								d.Set("type", "subnet_ip")
							}
							d.Set("subnet_ips", subnetIps)
							return nil
						}
					}
				}
			}
		}
	}
	d.SetId("")
	return fmt.Errorf("Unable to find route leak information for the Site Vrf %s", vrfName)
}

func resourceMSOSchemaSiteVrfRouteLeakImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	msoClient := m.(*client.Client)
	splitImport := strings.Split(d.Id(), "/")
	schemaId := splitImport[0]
	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	err = setRouteLeakFromSchema(d, schemaCont, schemaId, splitImport[2], splitImport[4], splitImport[6], getVrfRef(splitImport[8], splitImport[9], splitImport[10]))
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaSiteVrfRouteLeakCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Create", d.Id())
	msoClient := m.(*client.Client)
	siteId := d.Get("site_id").(string)
	prefixSubnets, includeAllSubnets := getSubnetDetails(d)
	path := fmt.Sprintf("/sites/%s-%s/vrfs/%s/routeLeak/-", d.Get("site_id").(string), d.Get("template_name").(string), d.Get("vrf_name").(string))
	vrfRouteLeakStruct := models.NewSchemaSiteVrfRouteLeak(
		"add", path, d.Get("tenant_name").(string), getTargetVrfRef(d), includeAllSubnets, prefixSubnets, []string{siteId},
	)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", d.Get("schema_id").(string)), vrfRouteLeakStruct)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Create finished successfully", d.Id())
	return resourceMSOSchemaSiteVrfRouteLeakRead(d, m)
}

func resourceMSOSchemaSiteVrfRouteLeakUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Update", d.Id())
	msoClient := m.(*client.Client)
	prefixSubnets, includeAllSubnets := getSubnetDetails(d)
	vrfRegionStruct := models.NewSchemaSiteVrfRouteLeak(
		"replace", d.Id(), d.Get("tenant_name").(string), getTargetVrfRef(d), includeAllSubnets, prefixSubnets, []string{d.Get("site_id").(string)},
	)
	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", d.Get("schema_id").(string)), vrfRegionStruct)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceMSOSchemaSiteVrfRouteLeakRead(d, m)
}

func resourceMSOSchemaSiteVrfRouteLeakRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	schemaCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	setRouteLeakFromSchema(d, schemaCont, schemaId, d.Get("site_id").(string), d.Get("template_name").(string), d.Get("vrf_name").(string), getTargetVrfRef(d))
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaSiteVrfRouteLeakDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Delete", d.Id())
	msoClient := m.(*client.Client)
	if d.Id() != "" {
		response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", d.Get("schema_id").(string)), models.GetRemovePatchPayload(d.Id()))
		if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
			return err
		}
		d.SetId("")
	}
	log.Printf("[DEBUG] Delete finished successfully")
	return nil
}
