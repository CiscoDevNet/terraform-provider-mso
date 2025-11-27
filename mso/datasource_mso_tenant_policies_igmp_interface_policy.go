package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOIGMPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOIGMPInterfacePolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the IGMP Interface Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the IGMP Interface Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the IGMP Interface Policy.",
			},
			"version3_asm": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable IGMP version 3 ASM.",
			},
			"fast_leave": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable fast leave.",
			},
			"report_link_local_groups": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable reporting link-local groups.",
			},
			"igmp_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IGMP version (v2 or v3).",
			},
			"group_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The group timeout value in seconds.",
			},
			"query_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The query interval value in seconds.",
			},
			"query_response_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The query response interval value in seconds.",
			},
			"last_member_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The last member query count value.",
			},
			"last_member_response_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The last member query response time value in seconds.",
			},
			"startup_query_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The startup query count value.",
			},
			"startup_query_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The startup query interval value in seconds.",
			},
			"querier_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The querier timeout value in seconds.",
			},
			"robustness_variable": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The robustness variable value.",
			},
			"state_limit_route_map_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the state limit route map policy for multicast.",
			},
			"report_policy_route_map_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the report policy route map for multicast.",
			},
			"static_report_route_map_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the static report route map for multicast.",
			},
			"maximum_multicast_entries": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum multicast entries value.",
			},
			"reserved_multicast_entries": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The reserved multicast entries value.",
			},
		},
	}
}

func dataSourceMSOIGMPInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO IGMP Interface Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "igmpInterfacePolicies")
	if err != nil {
		return err
	}

	setIGMPInterfacePolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO IGMP Interface Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
