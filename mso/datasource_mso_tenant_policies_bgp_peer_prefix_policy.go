package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOBGPPeerPrefixPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOBGPPeerPrefixPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the BGP Peer Prefix Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the BGP Peer Prefix Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the BGP Peer Prefix Policy.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action of the BGP Peer Prefix Policy (log, reject, restart, shutdown).",
			},
			"max_number_of_prefixes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of prefixes for the BGP Peer Prefix Policy.",
			},
			"threshold_percentage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The threshold percentage of the BGP Peer Prefix Policy.",
			},
			"restart_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The restart time of the BGP Peer Prefix Policy in seconds.",
			},
		},
	}
}

func dataSourceMSOBGPPeerPrefixPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "bgpPeerPrefixPolicies")
	if err != nil {
		return err
	}

	setBGPPeerPrefixPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO BGP Peer Prefix Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
