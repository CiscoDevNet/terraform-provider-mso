package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOMLDSnoopingPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOMLDSnoopingPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the MLD Snooping Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the MLD Snooping Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the MLD Snooping Policy.",
			},
			"admin_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The administrative state of the MLD Snooping Policy.",
			},
			"fast_leave_control": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether fast leave control is enabled.",
			},
			"querier_control": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether querier control is enabled.",
			},
			"querier_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The querier version (v1 or v2).",
			},
			"query_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The query interval in seconds.",
			},
			"query_response_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The query response interval in seconds.",
			},
			"last_member_query_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The last member query interval in seconds.",
			},
			"start_query_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The start query interval in seconds.",
			},
			"start_query_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The start query count.",
			},
		},
	}
}

func dataSourceMSOMLDSnoopingPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO MLD Snooping Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "mldSnoopPolicies")
	if err != nil {
		return err
	}

	setMLDSnoopingPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO MLD Snooping Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
