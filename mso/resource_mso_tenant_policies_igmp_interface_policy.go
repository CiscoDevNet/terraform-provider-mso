package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMSOIGMPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOIGMPInterfacePolicyCreate,
		Read:   resourceMSOIGMPInterfacePolicyRead,
		Update: resourceMSOIGMPInterfacePolicyUpdate,
		Delete: resourceMSOIGMPInterfacePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMSOIGMPInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The name of the IGMP Interface Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the IGMP Interface Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the IGMP Interface Policy.",
			},
			"version3_asm": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable IGMP version 3 ASM. Default: false (disabled) when unset during creation.",
			},
			"fast_leave": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable fast leave. Default: false (disabled) when unset during creation.",
			},
			"report_link_local_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable reporting link-local groups. Default: false (disabled) when unset during creation.",
			},
			"igmp_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v2", "v3",
				}, false),
				Description: "The IGMP version. Default: v2 when unset during creation.",
			},
			"group_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(3, 65535),
				Description:  "The group timeout value in seconds. Default: 260 when unset during creation. Range: 3-65535.",
			},
			"query_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 18000),
				Description:  "The query interval value in seconds. Default: 125 when unset during creation. Range: 1-18000.",
			},
			"query_response_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 25),
				Description:  "The query response interval value in seconds. Default: 10 when unset during creation. Range: 1-25.",
			},
			"last_member_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 5),
				Description:  "The last member query count value. Default: 2 when unset during creation. Range: 1-5.",
			},
			"last_member_response_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 25),
				Description:  "The last member query response time value in seconds. Default: 1 when unset during creation. Range: 1-25.",
			},
			"startup_query_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "The startup query count value. Default: 2 when unset during creation. Range: 1-10.",
			},
			"startup_query_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 18000),
				Description:  "The startup query interval value in seconds. Default: 31 when unset during creation. Range: 1-18000.",
			},
			"querier_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The querier timeout value in seconds. Default: 255 when unset during creation. Range: 1-65535.",
			},
			"robustness_variable": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 7),
				Description:  "The robustness variable value. Default: 2 when unset during creation. Range: 1-7.",
			},
			"state_limit_route_map_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the state limit route map policy for multicast.",
			},
			"report_policy_route_map_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the report policy route map for multicast.",
			},
			"static_report_route_map_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The UUID of the static report route map for multicast.",
			},
			"maximum_multicast_entries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateUint32Range(1, 4294967295),
				Description:  "The maximum multicast entries value. Default: 4294967295 when unset during creation. Range: 1-4294967295. Only applicable when state_limit_route_map is configured.",
			},
			"reserved_multicast_entries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateUint32Range(0, 4294967295),
				Description:  "The reserved multicast entries value. Default: 0 when unset during creation. Range: 0-4294967295.",
			},
		},
	}
}

func setIGMPInterfacePolicyData(d *schema.ResourceData, response *container.Container, templateId string) error {
	d.SetId(fmt.Sprintf("templateId/%s/IGMPInterfacePolicy/%s", templateId, models.StripQuotes(response.S("name").String())))
	d.Set("template_id", templateId)
	d.Set("name", models.StripQuotes(response.S("name").String()))
	d.Set("description", models.StripQuotes(response.S("description").String()))
	d.Set("uuid", models.StripQuotes(response.S("uuid").String()))

	if response.Exists("enableV3Asm") {
		if enableV3Asm, ok := response.S("enableV3Asm").Data().(bool); ok {
			d.Set("version3_asm", enableV3Asm)
		}
	}

	if response.Exists("enableFastLeaveControl") {
		if fastLeave, ok := response.S("enableFastLeaveControl").Data().(bool); ok {
			d.Set("fast_leave", fastLeave)
		}
	}

	if response.Exists("enableReportLinkLocalGroups") {
		if reportLinkLocal, ok := response.S("enableReportLinkLocalGroups").Data().(bool); ok {
			d.Set("report_link_local_groups", reportLinkLocal)
		}
	}

	if response.Exists("igmpQuerierVersion") {
		d.Set("igmp_version", models.StripQuotes(response.S("igmpQuerierVersion").String()))
	}
	if response.Exists("groupTimeout") {
		d.Set("group_timeout", int(response.S("groupTimeout").Data().(float64)))
	}
	if response.Exists("queryInterval") {
		d.Set("query_interval", int(response.S("queryInterval").Data().(float64)))
	}
	if response.Exists("queryResponseInterval") {
		d.Set("query_response_interval", int(response.S("queryResponseInterval").Data().(float64)))
	}
	if response.Exists("lastMemberCount") {
		d.Set("last_member_count", int(response.S("lastMemberCount").Data().(float64)))
	}
	if response.Exists("lastMemberResponseInterval") {
		d.Set("last_member_response_time", int(response.S("lastMemberResponseInterval").Data().(float64)))
	}
	if response.Exists("startQueryCount") {
		d.Set("startup_query_count", int(response.S("startQueryCount").Data().(float64)))
	}
	if response.Exists("startQueryInterval") {
		d.Set("startup_query_interval", int(response.S("startQueryInterval").Data().(float64)))
	}
	if response.Exists("querierTimeout") {
		d.Set("querier_timeout", int(response.S("querierTimeout").Data().(float64)))
	}
	if response.Exists("robustnessFactor") {
		d.Set("robustness_variable", int(response.S("robustnessFactor").Data().(float64)))
	}
	if response.Exists("maximumMulticastEntries") {
		d.Set("maximum_multicast_entries", int(response.S("maximumMulticastEntries").Data().(float64)))
	}
	if response.Exists("reservedMulticastEntries") {
		d.Set("reserved_multicast_entries", int(response.S("reservedMulticastEntries").Data().(float64)))
	}

	d.Set("state_limit_route_map_uuid", models.StripQuotes(response.S("stateLimitRouteMapRef").String()))
	d.Set("report_policy_route_map_uuid", models.StripQuotes(response.S("reportPolicyRouteMapRef").String()))
	d.Set("static_report_route_map_uuid", models.StripQuotes(response.S("staticReportRouteMapRef").String()))

	return nil
}

func resourceMSOIGMPInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Beginning Import: %v", d.Id())
	resourceMSOIGMPInterfacePolicyRead(d, m)
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Import Complete: %v", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceMSOIGMPInterfacePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Beginning Create: %v", d.Id())
	msoClient := m.(*client.Client)

	payload := map[string]interface{}{}

	payload["name"] = d.Get("name").(string)

	if description, ok := d.GetOk("description"); ok {
		payload["description"] = description.(string)
	}

	if version3Asm, ok := d.GetOk("version3_asm"); ok {
		payload["enableV3Asm"] = version3Asm.(bool)
	}

	if fastLeave, ok := d.GetOk("fast_leave"); ok {
		payload["enableFastLeaveControl"] = fastLeave.(bool)
	}

	if reportLinkLocal, ok := d.GetOk("report_link_local_groups"); ok {
		payload["enableReportLinkLocalGroups"] = reportLinkLocal.(bool)
	}

	if igmpVersion, ok := d.GetOk("igmp_version"); ok {
		payload["igmpQuerierVersion"] = igmpVersion.(string)
	}

	if groupTimeout, ok := d.GetOk("group_timeout"); ok {
		payload["groupTimeout"] = groupTimeout.(int)
	}

	if queryInterval, ok := d.GetOk("query_interval"); ok {
		payload["queryInterval"] = queryInterval.(int)
	}

	if queryResponseInterval, ok := d.GetOk("query_response_interval"); ok {
		payload["queryResponseInterval"] = queryResponseInterval.(int)
	}

	if lastMemberCount, ok := d.GetOk("last_member_count"); ok {
		payload["lastMemberCount"] = lastMemberCount.(int)
	}

	if lastMemberResponseTime, ok := d.GetOk("last_member_response_time"); ok {
		payload["lastMemberResponseInterval"] = lastMemberResponseTime.(int)
	}

	if startupQueryCount, ok := d.GetOk("startup_query_count"); ok {
		payload["startQueryCount"] = startupQueryCount.(int)
	}

	if startupQueryInterval, ok := d.GetOk("startup_query_interval"); ok {
		payload["startQueryInterval"] = startupQueryInterval.(int)
	}

	if querierTimeout, ok := d.GetOk("querier_timeout"); ok {
		payload["querierTimeout"] = querierTimeout.(int)
	}

	if robustnessVariable, ok := d.GetOk("robustness_variable"); ok {
		payload["robustnessFactor"] = robustnessVariable.(int)
	}

	if maxMulticastEntries, ok := d.GetOk("maximum_multicast_entries"); ok {
		payload["maximumMulticastEntries"] = maxMulticastEntries.(int)
	}

	if reservedMulticastEntries, ok := d.GetOk("reserved_multicast_entries"); ok {
		payload["reservedMulticastEntries"] = reservedMulticastEntries.(int)
	}

	if stateLimitUUID, ok := d.GetOk("state_limit_route_map_uuid"); ok {
		payload["stateLimitRouteMapRef"] = stateLimitUUID.(string)
	}

	if reportPolicyUUID, ok := d.GetOk("report_policy_route_map_uuid"); ok {
		payload["reportPolicyRouteMapRef"] = reportPolicyUUID.(string)
	}

	if staticReportUUID, ok := d.GetOk("static_report_route_map_uuid"); ok {
		payload["staticReportRouteMapRef"] = staticReportUUID.(string)
	}

	payloadModel := models.GetPatchPayload("add", "/tenantPolicyTemplate/template/igmpInterfacePolicies/-", payload)
	templateId := d.Get("template_id").(string)

	_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/IGMPInterfacePolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Create Complete: %v", d.Id())
	return resourceMSOIGMPInterfacePolicyRead(d, m)
}

func resourceMSOIGMPInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Beginning Read: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId, err := GetTemplateIdFromResourceId(d.Id())
	if err != nil {
		return err
	}

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyName, err := GetPolicyNameFromResourceId(d.Id(), "IGMPInterfacePolicy")
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "igmpInterfacePolicies")
	if err != nil {
		return err
	}

	setIGMPInterfacePolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Read Complete: %v", d.Id())
	return nil
}

func resourceMSOIGMPInterfacePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Beginning Update: %v", d.Id())
	msoClient := m.(*client.Client)
	templateId := d.Get("template_id").(string)

	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "igmpInterfacePolicies")
	if err != nil {
		return err
	}

	updatePath := fmt.Sprintf("/tenantPolicyTemplate/template/igmpInterfacePolicies/%d", policyIndex)

	payloadCont := container.New()
	payloadCont.Array()

	if d.HasChange("name") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/name", updatePath), d.Get("name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("description") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/description", updatePath), d.Get("description").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("version3_asm") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/enableV3Asm", updatePath), d.Get("version3_asm").(bool))
		if err != nil {
			return err
		}
	}

	if d.HasChange("fast_leave") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/enableFastLeaveControl", updatePath), d.Get("fast_leave").(bool))
		if err != nil {
			return err
		}
	}

	if d.HasChange("report_link_local_groups") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/enableReportLinkLocalGroups", updatePath), d.Get("report_link_local_groups").(bool))
		if err != nil {
			return err
		}
	}

	if d.HasChange("igmp_version") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/igmpQuerierVersion", updatePath), d.Get("igmp_version").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("group_timeout") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/groupTimeout", updatePath), d.Get("group_timeout").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("query_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/queryInterval", updatePath), d.Get("query_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("query_response_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/queryResponseInterval", updatePath), d.Get("query_response_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("last_member_count") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/lastMemberCount", updatePath), d.Get("last_member_count").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("last_member_response_time") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/lastMemberResponseInterval", updatePath), d.Get("last_member_response_time").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("startup_query_count") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/startQueryCount", updatePath), d.Get("startup_query_count").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("startup_query_interval") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/startQueryInterval", updatePath), d.Get("startup_query_interval").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("querier_timeout") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/querierTimeout", updatePath), d.Get("querier_timeout").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("robustness_variable") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/robustnessFactor", updatePath), d.Get("robustness_variable").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("maximum_multicast_entries") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/maximumMulticastEntries", updatePath), d.Get("maximum_multicast_entries").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("reserved_multicast_entries") {
		err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/reservedMulticastEntries", updatePath), d.Get("reserved_multicast_entries").(int))
		if err != nil {
			return err
		}
	}

	if d.HasChange("state_limit_route_map_uuid") {
		uuid := d.Get("state_limit_route_map_uuid").(string)
		if uuid != "" {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/stateLimitRouteMapRef", updatePath), uuid)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/stateLimitRouteMapRef", updatePath), nil)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("report_policy_route_map_uuid") {
		uuid := d.Get("report_policy_route_map_uuid").(string)
		if uuid != "" {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/reportPolicyRouteMapRef", updatePath), uuid)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/reportPolicyRouteMapRef", updatePath), nil)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("static_report_route_map_uuid") {
		uuid := d.Get("static_report_route_map_uuid").(string)
		if uuid != "" {
			err := addPatchPayloadToContainer(payloadCont, "replace", fmt.Sprintf("%s/staticReportRouteMapRef", updatePath), uuid)
			if err != nil {
				return err
			}
		} else {
			err := addPatchPayloadToContainer(payloadCont, "remove", fmt.Sprintf("%s/staticReportRouteMapRef", updatePath), nil)
			if err != nil {
				return err
			}
		}
	}
	err = doPatchRequest(msoClient, fmt.Sprintf("api/v1/templates/%s", templateId), payloadCont)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("templateId/%s/IGMPInterfacePolicy/%s", templateId, d.Get("name").(string)))
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Update Complete: %v", d.Id())
	return resourceMSOIGMPInterfacePolicyRead(d, m)
}

func resourceMSOIGMPInterfacePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Beginning Delete: %v", d.Id())
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	templateCont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(templateCont, "uuid", d.Get("uuid").(string), "tenantPolicyTemplate", "template", "igmpInterfacePolicies")
	if err != nil {
		return err
	}

	payloadModel := models.GetRemovePatchPayload(fmt.Sprintf("/tenantPolicyTemplate/template/igmpInterfacePolicies/%d", policyIndex))

	_, err = msoClient.PatchbyID(fmt.Sprintf("api/v1/templates/%s", templateId), payloadModel)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[DEBUG] MSO IGMP Interface Policy Resource - Delete Complete: %v", d.Id())
	return nil
}
