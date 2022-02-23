package mso

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateAnpEpg() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateAnpEpgCreate,
		Read:   resourceMSOSchemaTemplateAnpEpgRead,
		Update: resourceMSOSchemaTemplateAnpEpgUpdate,
		Delete: resourceMSOSchemaTemplateAnpEpgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateAnpEpgImport,
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
			"anp_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"bd_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"bd_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bd_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"useg_epg": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"intra_epg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"intersite_multicast_source": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"proxy_arp": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"preferred_group": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"epg_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"application",
					"service",
				}, false),
			},
			"cloud_service_epg_config": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"private",
								"public",
								"public_and_private",
							}, false),
						},
						"deployment_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"cloud_native",
								"cloud_native_managed",
								"third_party",
							}, false),
						},
						"service_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"azure_api_management_services",
								"azure_cosmos_db",
								"azure_databricks",
								"azure_sql",
								"azure_storage",
								"azure_storage_blob",
								"azure_storage_file",
								"azure_storage_queue",
								"azure_storage_table",
								"azure_kubernetes_services",
								"azure_ad_domain_services",
								"azure_contain_registry",
								"azure_key_vault",
								"redis_cache",
								"custom",
							}, false),
						},
						"custom_service_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func resourceMSOSchemaTemplateAnpEpgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	msoClient := m.(*client.Client)

	get_attribute := strings.Split(d.Id(), "/")
	schemaId := get_attribute[0]
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return nil, err
	}
	d.Set("schema_id", schemaId)
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return nil, fmt.Errorf("No Template found")
	}
	stateTemplate := get_attribute[2]
	found := false
	stateANP := get_attribute[4]
	stateEPG := get_attribute[6]
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			d.Set("template_name", apiTemplate)
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return nil, fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return nil, err
				}
				apiANP := models.StripQuotes(anpCont.S("name").String())
				if apiANP == stateANP {
					d.Set("anp_name", apiANP)
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return nil, fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return nil, err
						}
						apiEPG := models.StripQuotes(epgCont.S("name").String())
						if apiEPG == stateEPG {
							d.SetId(apiEPG)
							d.Set("name", apiEPG)
							d.Set("display_name", models.StripQuotes(epgCont.S("displayName").String()))
							d.Set("intra_epg", models.StripQuotes(epgCont.S("intraEpg").String()))
							d.Set("useg_epg", epgCont.S("uSegEpg").Data().(bool))
							if epgCont.Exists("mCastSource") {
								d.Set("intersite_multicast_source", epgCont.S("mCastSource").Data().(bool))
							}
							if epgCont.Exists("proxyArp") {
								d.Set("proxy_arp", epgCont.S("proxyArp").Data().(bool))
							}
							d.Set("preferred_group", epgCont.S("preferredGroup").Data().(bool))

							vrfRef := models.StripQuotes(epgCont.S("vrfRef").String())
							re_vrf := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
							match_vrf := re_vrf.FindStringSubmatch(vrfRef)
							if len(match_vrf) == 3 {
								d.Set("vrf_name", match_vrf[3])
								d.Set("vrf_schema_id", match_vrf[1])
								d.Set("vrf_template_name", match_vrf[2])
							}

							bdRef := models.StripQuotes(epgCont.S("bdRef").String())
							re_bd := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
							match_bd := re_bd.FindStringSubmatch(bdRef)
							if len(match_bd) == 3 {
								d.Set("bd_name", match_bd[3])
								d.Set("bd_schema_id", match_bd[1])
								d.Set("bd_template_name", match_bd[2])
							}
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("Unable to find the EPG %s", stateEPG)
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil

}

func resourceMSOSchemaTemplateAnpEpgCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template Anp Epg: Beginning Creation")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	Name := d.Get("name").(string)
	bdName := d.Get("bd_name").(string)
	vrfName := d.Get("vrf_name").(string)
	displayName := d.Get("display_name").(string)

	var intraEpg, vrf_schema_id, vrf_template_name, bd_schema_id, bd_template_name, epgType string
	var uSegEpg, intersiteMulticasteSource, preferredGroup, proxyArp bool

	if intra_epg, ok := d.GetOk("intra_epg"); ok {
		intraEpg = intra_epg.(string)
	}
	if useg_epg, ok := d.GetOk("useg_epg"); ok {
		uSegEpg = useg_epg.(bool)
	}
	if intersite_multicast_source, ok := d.GetOk("intersite_multicast_source"); ok {
		intersiteMulticasteSource = intersite_multicast_source.(bool)
	}
	if proxy_arp, ok := d.GetOk("proxy_arp"); ok {
		proxyArp = proxy_arp.(bool)
	}
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}

	if tempVar, ok := d.GetOk("bd_schema_id"); ok {
		bd_schema_id = tempVar.(string)
	} else {
		bd_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("bd_template_name"); ok {
		bd_template_name = tempVar.(string)
	} else {
		bd_template_name = templateName
	}
	if epg_type, ok := d.GetOk("epg_type"); ok {
		epgType = epg_type.(string)
	}

	service := make(map[string]interface{})
	if val, ok := d.GetOk("cloud_service_epg_config"); ok {
		service_list := val.(*schema.Set).List()
		for _, val := range service_list {

			map_service := make(map[string]interface{})
			inner_service := val.(map[string]interface{})

			if inner_service["access_type"] == "private" {
				map_service["accessType"] = "Private"
			} else if inner_service["access_type"] == "public" {
				map_service["accessType"] = "Public"
			} else if inner_service["access_type"] == "public_and_private" {
				map_service["accessType"] = "PublicAndPrivate"
			}

			if inner_service["deployment_type"] == "cloud_native" {
				map_service["deploymentType"] = "CloudNative"
			} else if inner_service["deployment_type"] == "cloud_native_managed" {
				map_service["deploymentType"] = "CloudNativeManaged"
			} else if inner_service["deployment_type"] == "third_party" {
				map_service["deploymentType"] = "Third-party"
			}

			if inner_service["service_type"] == "azure_api_management_services" {
				map_service["serviceType"] = "Azure-ApiManagement"
			} else if inner_service["service_type"] == "azure_cosmos_db" {
				map_service["serviceType"] = "Azure-CosmosDB"
			} else if inner_service["service_type"] == "azure_databricks" {
				map_service["serviceType"] = "Azure-Databricks"
			} else if inner_service["service_type"] == "azure_sql" {
				map_service["serviceType"] = "Azure-SqlServer"
			} else if inner_service["service_type"] == "azure_storage" {
				map_service["serviceType"] = "Azure-Storage"
			} else if inner_service["service_type"] == "azure_storage_blob" {
				map_service["serviceType"] = "Azure-StorageBlob"
			} else if inner_service["service_type"] == "azure_storage_file" {
				map_service["serviceType"] = "Azure-StorageFile"
			} else if inner_service["service_type"] == "azure_storage_queue" {
				map_service["serviceType"] = "Azure-StorageQueue"
			} else if inner_service["service_type"] == "azure_storage_table" {
				map_service["serviceType"] = "Azure-StorageTable"
			} else if inner_service["service_type"] == "azure_kubernetes_services" {
				map_service["serviceType"] = "Azure-AksCluster"
			} else if inner_service["service_type"] == "azure_ad_domain_services" {
				map_service["serviceType"] = "Azure-ADDS"
			} else if inner_service["service_type"] == "azure_contain_registry" {
				map_service["serviceType"] = "Azure-ContainerRegistry"
			} else if inner_service["service_type"] == "azure_key_vault" {
				map_service["serviceType"] = "Azure-KeyVault"
			} else if inner_service["service_type"] == "redis_cache" {
				map_service["serviceType"] = "Azure-Redis"
			} else if inner_service["service_type"] == "custom" {
				map_service["serviceType"] = "Custom"
				map_service["customSvcType"] = inner_service["custom_service_type"]
			}
			service = map_service
		}
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/-", templateName, anpName)
	anpEpgStruct := models.NewTemplateAnpEpg("add", path, Name, displayName, intraEpg, epgType, uSegEpg, intersiteMulticasteSource, preferredGroup, proxyArp, vrfRefMap, bdRefMap, service)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaTemplateAnpEpgRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)

	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}
	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("No Template found")
	}
	stateTemplate := d.Get("template_name").(string)
	found := false
	stateANP := d.Get("anp_name").(string)
	stateEPG := d.Get("name").(string)
	for i := 0; i < count; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		apiTemplate := models.StripQuotes(tempCont.S("name").String())

		if apiTemplate == stateTemplate {
			anpCount, err := tempCont.ArrayCount("anps")
			if err != nil {
				return fmt.Errorf("Unable to get ANP list")
			}
			for j := 0; j < anpCount; j++ {
				anpCont, err := tempCont.ArrayElement(j, "anps")
				if err != nil {
					return err
				}
				apiANP := models.StripQuotes(anpCont.S("name").String())
				if apiANP == stateANP {
					epgCount, err := anpCont.ArrayCount("epgs")
					if err != nil {
						return fmt.Errorf("Unable to get EPG list")
					}
					for k := 0; k < epgCount; k++ {
						epgCont, err := anpCont.ArrayElement(k, "epgs")
						if err != nil {
							return err
						}
						apiEPG := models.StripQuotes(epgCont.S("name").String())
						if apiEPG == stateEPG {
							d.SetId(apiEPG)
							d.Set("schema_id", schemaId)
							d.Set("name", apiEPG)
							d.Set("template_name", apiTemplate)
							d.Set("display_name", models.StripQuotes(epgCont.S("displayName").String()))
							d.Set("intra_epg", models.StripQuotes(epgCont.S("intraEpg").String()))
							d.Set("useg_epg", epgCont.S("uSegEpg").Data().(bool))
							if epgCont.Exists("mCastSource") {
								d.Set("intersite_multicast_source", epgCont.S("mCastSource").Data().(bool))
							}
							if epgCont.Exists("proxyArp") {
								d.Set("proxy_arp", epgCont.S("proxyArp").Data().(bool))
							}
							d.Set("preferred_group", epgCont.S("preferredGroup").Data().(bool))
							d.Set("epg_type", models.StripQuotes(epgCont.S("epgType").String()))
							if d.Get("epg_type") == "service" {
								services := make([]interface{}, 0)
								map_service := make(map[string]interface{})

								servicesCont := epgCont.S("cloudServiceEpgConfig")
								log.Printf("SERVSERV %s", servicesCont)
								if err != nil {
									return fmt.Errorf("Unable to parse the Cloud Services")
								}

								map_service["deployment_type"] = models.StripQuotes(servicesCont.S("deploymentType").String())
								map_service["access_type"] = models.StripQuotes(servicesCont.S("accessType").String())
								map_service["service_type"] = models.StripQuotes(servicesCont.S("serviceType").String())
								log.Printf("SERVSERV %s", map_service)
								services = append(services, map_service)

								d.Set("cloud_service_epg_config", services)
							}

							vrfRef := models.StripQuotes(epgCont.S("vrfRef").String())
							re_vrf := regexp.MustCompile("/schemas/(.*)/templates/(.*)/vrfs/(.*)")
							match_vrf := re_vrf.FindStringSubmatch(vrfRef)
							if len(match_vrf) == 3 {
								d.Set("vrf_name", match_vrf[3])
								d.Set("vrf_schema_id", match_vrf[1])
								d.Set("vrf_template_name", match_vrf[2])
							}
							bdRef := models.StripQuotes(epgCont.S("bdRef").String())
							re_bd := regexp.MustCompile("/schemas/(.*)/templates/(.*)/bds/(.*)")
							match_bd := re_bd.FindStringSubmatch(bdRef)
							if len(match_bd) == 3 {
								d.Set("bd_name", match_bd[3])
								d.Set("bd_schema_id", match_bd[1])
								d.Set("bd_template_name", match_bd[2])
							}
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
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil

}

func resourceMSOSchemaTemplateAnpEpgUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	Name := d.Get("name").(string)
	bdName := d.Get("bd_name").(string)
	vrfName := d.Get("vrf_name").(string)
	displayName := d.Get("display_name").(string)

	var intraEpg, vrf_schema_id, vrf_template_name, bd_schema_id, bd_template_name, epgType string
	var uSegEpg, intersiteMulticasteSource, preferredGroup, proxyArp bool

	if intra_epg, ok := d.GetOk("intra_epg"); ok {
		intraEpg = intra_epg.(string)
	}
	if useg_epg, ok := d.GetOk("useg_epg"); ok {
		uSegEpg = useg_epg.(bool)
	}
	if intersite_multicast_source, ok := d.GetOk("intersite_multicast_source"); ok {
		intersiteMulticasteSource = intersite_multicast_source.(bool)
	}
	if proxy_arp, ok := d.GetOk("proxy_arp"); ok {
		proxyArp = proxy_arp.(bool)
	}
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	if tempVar, ok := d.GetOk("bd_schema_id"); ok {
		bd_schema_id = tempVar.(string)
	} else {
		bd_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("bd_template_name"); ok {
		bd_template_name = tempVar.(string)
	} else {
		bd_template_name = templateName
	}
	if epg_type, ok := d.GetOk("epg_type"); ok {
		epgType = epg_type.(string)
	}

	service := make(map[string]interface{})
	if val, ok := d.GetOk("cloud_service_epg_config"); ok {
		service_list := val.(*schema.Set).List()
		for _, val := range service_list {

			map_service := make(map[string]interface{})
			inner_service := val.(map[string]interface{})

			if inner_service["access_type"] == "private" {
				map_service["accessType"] = "Private"
			} else if inner_service["access_type"] == "public" {
				map_service["accessType"] = "Public"
			} else if inner_service["access_type"] == "public_and_private" {
				map_service["accessType"] = "PublicAndPrivate"
			}

			if inner_service["deployment_type"] == "cloud_native" {
				map_service["deploymentType"] = "CloudNative"
			} else if inner_service["deployment_type"] == "cloud_native_managed" {
				map_service["deploymentType"] = "CloudNativeManaged"
			} else if inner_service["deployment_type"] == "third_party" {
				map_service["deploymentType"] = "Third-party"
			}

			if inner_service["service_type"] == "azure_api_management_services" {
				map_service["serviceType"] = "Azure-ApiManagement"
			} else if inner_service["service_type"] == "azure_cosmos_db" {
				map_service["serviceType"] = "Azure-CosmosDB"
			} else if inner_service["service_type"] == "azure_databricks" {
				map_service["serviceType"] = "Azure-Databricks"
			} else if inner_service["service_type"] == "azure_sql" {
				map_service["serviceType"] = "Azure-SqlServer"
			} else if inner_service["service_type"] == "azure_storage" {
				map_service["serviceType"] = "Azure-Storage"
			} else if inner_service["service_type"] == "azure_storage_blob" {
				map_service["serviceType"] = "Azure-StorageBlob"
			} else if inner_service["service_type"] == "azure_storage_file" {
				map_service["serviceType"] = "Azure-StorageFile"
			} else if inner_service["service_type"] == "azure_storage_queue" {
				map_service["serviceType"] = "Azure-StorageQueue"
			} else if inner_service["service_type"] == "azure_storage_table" {
				map_service["serviceType"] = "Azure-StorageTable"
			} else if inner_service["service_type"] == "azure_kubernetes_services" {
				map_service["serviceType"] = "Azure-AksCluster"
			} else if inner_service["service_type"] == "azure_ad_domain_services" {
				map_service["serviceType"] = "Azure-ADDS"
			} else if inner_service["service_type"] == "azure_contain_registry" {
				map_service["serviceType"] = "Azure-ContainerRegistry"
			} else if inner_service["service_type"] == "azure_key_vault" {
				map_service["serviceType"] = "Azure-KeyVault"
			} else if inner_service["service_type"] == "redis_cache" {
				map_service["serviceType"] = "Azure-Redis"
			} else if inner_service["service_type"] == "custom" {
				map_service["serviceType"] = "Custom"
				map_service["customSvcType"] = inner_service["custom_service_type"]
			}
			service = map_service
		}
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s", templateName, anpName, d.Id())
	anpEpgStruct := models.NewTemplateAnpEpg("replace", path, Name, displayName, intraEpg, epgType, uSegEpg, intersiteMulticasteSource, preferredGroup, proxyArp, vrfRefMap, bdRefMap, service)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)

	if err != nil {
		return err
	}
	return resourceMSOSchemaTemplateAnpEpgRead(d, m)
}

func resourceMSOSchemaTemplateAnpEpgDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Template BD: Beginning Update")
	msoClient := m.(*client.Client)

	schemaId := d.Get("schema_id").(string)
	templateName := d.Get("template_name").(string)
	anpName := d.Get("anp_name").(string)
	Name := d.Get("name").(string)
	bdName := d.Get("bd_name").(string)
	vrfName := d.Get("vrf_name").(string)
	displayName := d.Get("display_name").(string)

	var intraEpg, vrf_schema_id, vrf_template_name, bd_schema_id, bd_template_name, epgType string
	var uSegEpg, intersiteMulticasteSource, preferredGroup, proxyArp bool

	if intra_epg, ok := d.GetOk("intra_epg"); ok {
		intraEpg = intra_epg.(string)
	}
	if useg_epg, ok := d.GetOk("useg_epg"); ok {
		uSegEpg = useg_epg.(bool)
	}
	if intersite_multicast_source, ok := d.GetOk("intersite_multicast_source"); ok {
		intersiteMulticasteSource = intersite_multicast_source.(bool)
	}
	if proxy_arp, ok := d.GetOk("proxy_arp"); ok {
		proxyArp = proxy_arp.(bool)
	}
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}
	if tempVar, ok := d.GetOk("vrf_schema_id"); ok {
		vrf_schema_id = tempVar.(string)
	} else {
		vrf_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("vrf_template_name"); ok {
		vrf_template_name = tempVar.(string)
	} else {
		vrf_template_name = templateName
	}
	if tempVar, ok := d.GetOk("bd_schema_id"); ok {
		bd_schema_id = tempVar.(string)
	} else {
		bd_schema_id = schemaId
	}
	if tempVar, ok := d.GetOk("bd_template_name"); ok {
		bd_template_name = tempVar.(string)
	} else {
		bd_template_name = templateName
	}
	if epg_type, ok := d.GetOk("epg_type"); ok {
		epgType = epg_type.(string)
	}

	service := make(map[string]interface{})
	if val, ok := d.GetOk("cloud_service_epg_config"); ok {
		service_list := val.(*schema.Set).List()
		for _, val := range service_list {

			map_service := make(map[string]interface{})
			inner_service := val.(map[string]interface{})

			if inner_service["access_type"] == "private" {
				map_service["accessType"] = "Private"
			} else if inner_service["access_type"] == "public" {
				map_service["accessType"] = "Public"
			} else if inner_service["access_type"] == "public_and_private" {
				map_service["accessType"] = "PublicAndPrivate"
			}

			if inner_service["deployment_type"] == "cloud_native" {
				map_service["deploymentType"] = "CloudNative"
			} else if inner_service["deployment_type"] == "cloud_native_managed" {
				map_service["deploymentType"] = "CloudNativeManaged"
			} else if inner_service["deployment_type"] == "third_party" {
				map_service["deploymentType"] = "Third-party"
			}

			if inner_service["service_type"] == "azure_api_management_services" {
				map_service["serviceType"] = "Azure-ApiManagement"
			} else if inner_service["service_type"] == "azure_cosmos_db" {
				map_service["serviceType"] = "Azure-CosmosDB"
			} else if inner_service["service_type"] == "azure_databricks" {
				map_service["serviceType"] = "Azure-Databricks"
			} else if inner_service["service_type"] == "azure_sql" {
				map_service["serviceType"] = "Azure-SqlServer"
			} else if inner_service["service_type"] == "azure_storage" {
				map_service["serviceType"] = "Azure-Storage"
			} else if inner_service["service_type"] == "azure_storage_blob" {
				map_service["serviceType"] = "Azure-StorageBlob"
			} else if inner_service["service_type"] == "azure_storage_file" {
				map_service["serviceType"] = "Azure-StorageFile"
			} else if inner_service["service_type"] == "azure_storage_queue" {
				map_service["serviceType"] = "Azure-StorageQueue"
			} else if inner_service["service_type"] == "azure_storage_table" {
				map_service["serviceType"] = "Azure-StorageTable"
			} else if inner_service["service_type"] == "azure_kubernetes_services" {
				map_service["serviceType"] = "Azure-AksCluster"
			} else if inner_service["service_type"] == "azure_ad_domain_services" {
				map_service["serviceType"] = "Azure-ADDS"
			} else if inner_service["service_type"] == "azure_contain_registry" {
				map_service["serviceType"] = "Azure-ContainerRegistry"
			} else if inner_service["service_type"] == "azure_key_vault" {
				map_service["serviceType"] = "Azure-KeyVault"
			} else if inner_service["service_type"] == "redis_cache" {
				map_service["serviceType"] = "Azure-Redis"
			} else if inner_service["service_type"] == "custom" {
				map_service["serviceType"] = "Custom"
				map_service["customSvcType"] = inner_service["custom_service_type"]
			}
			service = map_service
		}
	}

	vrfRefMap := make(map[string]interface{})
	vrfRefMap["schemaId"] = vrf_schema_id
	vrfRefMap["templateName"] = vrf_template_name
	vrfRefMap["vrfName"] = vrfName

	bdRefMap := make(map[string]interface{})
	bdRefMap["schemaId"] = bd_schema_id
	bdRefMap["templateName"] = bd_template_name
	bdRefMap["bdName"] = bdName

	path := fmt.Sprintf("/templates/%s/anps/%s/epgs/%s", templateName, anpName, d.Id())
	anpEpgStruct := models.NewTemplateAnpEpg("remove", path, Name, displayName, intraEpg, epgType, uSegEpg, intersiteMulticasteSource, preferredGroup, proxyArp, vrfRefMap, bdRefMap, service)

	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), anpEpgStruct)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}
	d.SetId("")
	return nil
}
