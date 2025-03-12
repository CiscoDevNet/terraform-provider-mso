package mso

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOSchemaTemplateVrf() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOSchemaTemplateVrfCreate,
		Update: resourceMSOSchemaTemplateVrfUpdate,
		Read:   resourceMSOSchemaTemplateVrfRead,
		Delete: resourceMSOSchemaTemplateVrfDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMSOSchemaTemplateVrfImport,
		},

		Schema: (map[string]*schema.Schema{

			"schema_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"template": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"layer3_multicast": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"vzany": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"ip_data_plane_learning": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"preferred_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"site_aware_policy_enforcement": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"rendezvous_points": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"static",
								"fabric",
								"unknown",
							}, false),
							Required: true,
						},
						"mutlicast_route_map_policy_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func resourceMSOSchemaTemplateVrfImport(d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Schema Template Vrf: Beginning Import")
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
		return nil, fmt.Errorf("no template found")
	}
	templateName := get_attribute[2]
	vrfName := get_attribute[4]
	found := false
	for i := range count {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return nil, err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())
		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
			vrfCount, err := tempCont.ArrayCount("vrfs")
			if err != nil {
				return nil, fmt.Errorf("no vrf found")
			}
			for j := range vrfCount {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")
				if err != nil {
					return nil, err
				}
				currentVrfName := models.StripQuotes(vrfCont.S("name").String())
				if currentVrfName == vrfName {
					d.SetId(currentVrfName)
					d.Set("name", currentVrfName)
					d.Set("display_name", models.StripQuotes(vrfCont.S("displayName").String()))
					d.Set("uuid", models.StripQuotes(vrfCont.S("uuid").String()))
					if vrfCont.Exists("l3MCast") {
						l3Mcast, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("l3MCast").String()))
						d.Set("layer3_multicast", l3Mcast)
					}
					if vrfCont.Exists("vzAnyEnabled") {
						vzAnyEnabled, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("vzAnyEnabled").String()))
						d.Set("vzany", vzAnyEnabled)
					}
					if vrfCont.Exists("ipDataPlaneLearning") {
						d.Set("ip_data_plane_learning", models.StripQuotes(vrfCont.S("ipDataPlaneLearning").String()))
					}
					if vrfCont.Exists("preferredGroup") {
						preferredGroup, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("preferredGroup").String()))
						d.Set("preferred_group", preferredGroup)
					}
					if vrfCont.Exists("description") {
						d.Set("description", models.StripQuotes(vrfCont.S("description").String()))
					}
					if vrfCont.Exists("siteAwarePolicyEnforcementMode") {
						siteAwarePolicyEnforcementMode, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("siteAwarePolicyEnforcementMode").String()))
						d.Set("site_aware_policy_enforcement", siteAwarePolicyEnforcementMode)
					}
					if vrfCont.Exists("rpConfigs") {
						rpCount, err := vrfCont.ArrayCount("rpConfigs")
						if err != nil {
							return nil, fmt.Errorf("no rendezvous points found")
						}
						rendezvousPoints := make([]any, 0)
						for k := range rpCount {
							rpCont, err := vrfCont.ArrayElement(k, "rpConfigs")
							if err != nil {
								return nil, fmt.Errorf("unable to parse the rendezvous points list")
							}
							rpConfig := make(map[string]any)
							rpConfig["ip_address"] = models.StripQuotes(rpCont.S("ipAddress").String())
							rpConfig["type"] = models.StripQuotes(rpCont.S("rpType").String())
							rpConfig["mutlicast_route_map_policy_uuid"] = models.StripQuotes(rpCont.S("mcastRtMapPolicyRef").String())
							rendezvousPoints = append(rendezvousPoints, rpConfig)
						}
						d.Set("rendezvous_points", rendezvousPoints)
					}
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	if !found {
		d.SetId("")
		d.Set("name", "")
		d.Set("display_name", "")
		d.Set("uuid", "")
	}

	log.Printf("[DEBUG] %s: Schema Template Vrf Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOSchemaTemplateVrfCreate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] Schema Template Vrf: Beginning Creation")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template"); ok {
		templateName = template.(string)
	}

	var Name string
	if name, ok := d.GetOk("name"); ok {
		Name = name.(string)
	}

	var displayName string
	if display_name, ok := d.GetOk("display_name"); ok {
		displayName = display_name.(string)
	}

	var l3m bool
	if L3M, ok := d.GetOk("layer3_multicast"); ok {
		l3m = L3M.(bool)
	}

	var vzany bool
	if vzAny, ok := d.GetOk("vzany"); ok {
		vzany = vzAny.(bool)
	}

	var ipDataPlaneLearning string
	if ip_data_plane_learning, ok := d.GetOk("ip_data_plane_learning"); ok {
		ipDataPlaneLearning = ip_data_plane_learning.(string)
	}

	var preferredGroup bool
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}

	var description string
	if tempVar, ok := d.GetOk("description"); ok {
		description = tempVar.(string)
	}

	var siteAwarePolicyEnforcementMode bool
	if site_aware_policy_enforcement, ok := d.GetOk("site_aware_policy_enforcement"); ok {
		siteAwarePolicyEnforcementMode = site_aware_policy_enforcement.(bool)
	}

	rendezvousPoints := make([]any, 0, 1)
	if val, ok := d.GetOk("rendezvous_points"); ok {
		rp_list := val.(*schema.Set).List()
		for _, val := range rp_list {

			rpConfig := make(map[string]any)
			rendezvousPoint := val.(map[string]any)
			if rendezvousPoint["ip_address"] != "" {
				rpConfig["ipAddress"] = fmt.Sprintf("%v", rendezvousPoint["ip_address"])
			}
			if rendezvousPoint["type"] != "" {
				rpConfig["rpType"] = fmt.Sprintf("%v", rendezvousPoint["type"])
			}
			if rendezvousPoint["mutlicast_route_map_policy_uuid"] != "" {
				rpConfig["mcastRtMapPolicyRef"] = fmt.Sprintf("%v", rendezvousPoint["mutlicast_route_map_policy_uuid"])
			}
			rendezvousPoints = append(rendezvousPoints, rpConfig)
		}
	}

	schemaTemplateVrfApp := models.NewSchemaTemplateVrf("add", fmt.Sprintf("/templates/%s/vrfs/-", templateName), Name, displayName, ipDataPlaneLearning, description, l3m, vzany, preferredGroup, siteAwarePolicyEnforcementMode, rendezvousPoints)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateVrfApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", Name))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateVrfRead(d, m)
}

func resourceMSOSchemaTemplateVrfUpdate(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] Schema Template Vrf: Beginning Creation")
	msoClient := m.(*client.Client)

	var schemaId string
	if schema_id, ok := d.GetOk("schema_id"); ok {
		schemaId = schema_id.(string)
	}

	var templateName string
	if template, ok := d.GetOk("template"); ok {
		templateName = template.(string)
	}

	var Name string
	if name, ok := d.GetOk("name"); ok {
		Name = name.(string)
	}

	var displayName string
	if display_name, ok := d.GetOk("display_name"); ok {
		displayName = display_name.(string)
	}

	var l3m bool
	if L3M, ok := d.GetOk("layer3_multicast"); ok {
		l3m = L3M.(bool)
	}

	var vzany bool
	if vzAny, ok := d.GetOk("vzany"); ok {
		vzany = vzAny.(bool)
	}

	var ipDataPlaneLearning string
	if ip_data_plane_learning, ok := d.GetOk("ip_data_plane_learning"); ok {
		ipDataPlaneLearning = ip_data_plane_learning.(string)
	}

	var preferredGroup bool
	if preferred_group, ok := d.GetOk("preferred_group"); ok {
		preferredGroup = preferred_group.(bool)
	}

	var description string
	if tempVar, ok := d.GetOk("description"); ok {
		description = tempVar.(string)
	}

	var siteAwarePolicyEnforcementMode bool
	if site_aware_policy_enforcement, ok := d.GetOk("site_aware_policy_enforcement"); ok {
		siteAwarePolicyEnforcementMode = site_aware_policy_enforcement.(bool)
	}

	rendezvousPoints := make([]any, 0, 1)
	if val, ok := d.GetOk("rendezvous_points"); ok {
		rp_list := val.(*schema.Set).List()
		for _, val := range rp_list {

			rpConfig := make(map[string]any)
			rendezvousPoint := val.(map[string]any)
			if rendezvousPoint["ip_address"] != "" {
				rpConfig["ipAddress"] = fmt.Sprintf("%v", rendezvousPoint["ip_address"])
			}
			if rendezvousPoint["type"] != "" {
				rpConfig["rpType"] = fmt.Sprintf("%v", rendezvousPoint["type"])
			}
			if rendezvousPoint["mutlicast_route_map_policy_uuid"] != "" {
				rpConfig["mcastRtMapPolicyRef"] = fmt.Sprintf("%v", rendezvousPoint["mutlicast_route_map_policy_uuid"])
			}
			rendezvousPoints = append(rendezvousPoints, rpConfig)
		}
	}

	schemaTemplateVrfApp := models.NewSchemaTemplateVrf("replace", fmt.Sprintf("/templates/%s/vrfs/%s", templateName, Name), Name, displayName, ipDataPlaneLearning, description, l3m, vzany, preferredGroup, siteAwarePolicyEnforcementMode, rendezvousPoints)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), schemaTemplateVrfApp)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(fmt.Sprintf("%v", Name))
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceMSOSchemaTemplateVrfRead(d, m)
}

func resourceMSOSchemaTemplateVrfRead(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), cont, d)
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("no template found")
	}

	templateName := d.Get("template").(string)
	vrfName := d.Get("name").(string)
	found := false

	for i := range count {

		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())

		if currentTemplateName == templateName {
			d.Set("template", currentTemplateName)
			vrfCount, err := tempCont.ArrayCount("vrfs")

			if err != nil {
				return fmt.Errorf("no vrf found")
			}
			for j := range vrfCount {
				vrfCont, err := tempCont.ArrayElement(j, "vrfs")

				if err != nil {
					return err
				}
				currentVrfName := models.StripQuotes(vrfCont.S("name").String())
				log.Println("currentvrfname", currentVrfName)
				if currentVrfName == vrfName {
					log.Println("found correct vrfname")
					d.SetId(currentVrfName)
					d.Set("name", currentVrfName)
					if vrfCont.Exists("displayName") {
						d.Set("display_name", models.StripQuotes(vrfCont.S("displayName").String()))
					}
					if vrfCont.Exists("uuid") {
						d.Set("uuid", models.StripQuotes(vrfCont.S("uuid").String()))
					}
					if vrfCont.Exists("l3MCast") {
						l3Mcast, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("l3MCast").String()))
						d.Set("layer3_multicast", l3Mcast)
					}
					if vrfCont.Exists("vzAnyEnabled") {
						vzAnyEnabled, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("vzAnyEnabled").String()))
						d.Set("vzany", vzAnyEnabled)
					}
					if vrfCont.Exists("ipDataPlaneLearning") {
						d.Set("ip_data_plane_learning", models.StripQuotes(vrfCont.S("ipDataPlaneLearning").String()))
					}
					if vrfCont.Exists("preferredGroup") {
						preferredGroup, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("preferredGroup").String()))
						d.Set("preferred_group", preferredGroup)
					}
					if vrfCont.Exists("description") {
						d.Set("description", models.StripQuotes(vrfCont.S("description").String()))
					}
					if vrfCont.Exists("siteAwarePolicyEnforcementMode") {
						siteAwarePolicyEnforcementMode, _ := strconv.ParseBool(models.StripQuotes(vrfCont.S("siteAwarePolicyEnforcementMode").String()))
						d.Set("site_aware_policy_enforcement", siteAwarePolicyEnforcementMode)
					}
					if vrfCont.Exists("rpConfigs") {
						rpCount, err := vrfCont.ArrayCount("rpConfigs")
						if err != nil {
							return err
						}
						rendezvousPoints := make([]any, 0)
						for k := range rpCount {
							rpCont, err := vrfCont.ArrayElement(k, "rpConfigs")
							if err != nil {
								return err
							}
							rpConfig := make(map[string]any)
							rpConfig["ip_address"] = models.StripQuotes(rpCont.S("ipAddress").String())
							rpConfig["type"] = models.StripQuotes(rpCont.S("rpType").String())
							rpConfig["mutlicast_route_map_policy_uuid"] = models.StripQuotes(rpCont.S("mcastRtMapPolicyRef").String())
							rendezvousPoints = append(rendezvousPoints, rpConfig)
						}
						d.Set("rendezvous_points", rendezvousPoints)
					}
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	if !found {
		d.SetId("")
		d.Set("name", "")
		d.Set("display_name", "")
		d.Set("uuid", "")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOSchemaTemplateVrfDelete(d *schema.ResourceData, m any) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	msoClient := m.(*client.Client)
	schemaId := d.Get("schema_id").(string)
	template := d.Get("template").(string)
	name := d.Get("name").(string)

	vrfRemovePatchPayload := models.GetRemovePatchPayload(fmt.Sprintf("/templates/%s/vrfs/%s", template, name))
	response, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", schemaId), vrfRemovePatchPayload)

	// Ignoring Error with code 141: Resource Not Found when deleting
	if err != nil && !(response.Exists("code") && response.S("code").String() == "141") {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return nil
}
