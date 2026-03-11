package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSORouteMapPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSORouteMapPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Route Map Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Route Map Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the Route Map Policy.",
			},
		},
	}
}

func dataSourceMSORouteMapPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Route Map Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "routeMapPolicies")
	if err != nil {
		return err
	}

	setRouteMapPolicyData(d, policy, templateId)

	log.Printf("[DEBUG] MSO Route Map Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
