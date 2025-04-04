package mso

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceMSOSchemaTemplateVrf() *schema.Resource {
	return &schema.Resource{

		Read: datasourceMSOSchemaTemplateVrfRead,

		SchemaVersion: version,
		Schema: (map[string]*schema.Schema{
			"schema_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"layer3_multicast": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vzany": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ip_data_plane_learning": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_group": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_aware_policy_enforcement": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rendezvous_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mutlicast_route_map_policy_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}
func datasourceMSOSchemaTemplateVrfRead(d *schema.ResourceData, m any) error {
	schemaId := d.Get("schema_id").(string)
	msoClient := m.(*client.Client)
	cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
	if err != nil {
		return err
	}

	count, err := cont.ArrayCount("templates")
	if err != nil {
		return fmt.Errorf("no template found")
	}

	templateName := d.Get("template").(string)
	vrfName := d.Get("name").(string)

	found := false
	for i := 0; i < count && !found; i++ {
		tempCont, err := cont.ArrayElement(i, "templates")
		if err != nil {
			return err
		}
		currentTemplateName := models.StripQuotes(tempCont.S("name").String())
		if currentTemplateName == templateName {
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
				if currentVrfName == vrfName {
					d.SetId(fmt.Sprintf("%s/templates/%s/vrfs/%s", schemaId, templateName, vrfName))
					d.Set("name", currentVrfName)
					d.Set("template", currentTemplateName)
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
	}

	if !found {
		d.SetId("")
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}
