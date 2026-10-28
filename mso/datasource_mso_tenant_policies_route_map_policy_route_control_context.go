package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSORouteMapPolicyContext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSORouteMapPolicyContextRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"route_map_policy_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the parent Route Map Policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Route Control context.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the context.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action of the context.",
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The order of the context.",
			},
			"match_rule_uuids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A set of Match Rule Policy UUIDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"set_rule_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Set Rule Policy.",
			},
		},
	}
}

func dataSourceMSORouteMapPolicyContextRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Context Data Source - Beginning Read")
	msoClient := m.(*client.Client)
	templateID := d.Get("template_id").(string)
	policyUUID := d.Get("route_map_policy_uuid").(string)
	contextName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}

	policyIndex, err := GetPolicyIndexByKeyAndValue(response, "uuid", policyUUID, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	policy := response.S("tenantPolicyTemplate", "template", "routeMapPolicies").Index(policyIndex)
	match, err := getRouteMapPolicyContext(policy, contextName, policyUUID)
	if err != nil {
		return err
	}
	if err := setRouteMapPolicyContextData(d, templateID, policyUUID, match); err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Route Map Policy Context Data Source - Read Complete: %v", d.Id())
	return nil
}
