package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSOMatchRulePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOMatchRulePolicyRead,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Match Rule Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the Match Rule Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the Match Rule Policy.",
			},
			"match_prefix": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The prefixes to match.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPv4 or IPv6 prefix to match.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the match prefix.",
						},
						"aggregate": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to match a range of prefix lengths.",
						},
						"from_prefix": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum prefix length.",
						},
						"to_prefix": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum prefix length.",
						},
					},
				},
			},
			"match_community_term": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The community terms to match.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the community term.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the community term.",
						},
						"match_community_factor": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The community factors within the term.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"community": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The regular, extended, or well-known community to match.",
									},
									"scope": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The community scope.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the community factor.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOMatchRulePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Match Rule Policy Data Source - Beginning Read")
	templateID := d.Get("template_id").(string)
	response, err := m.(*client.Client).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	policy, err := GetPolicyByName(response, d.Get("name").(string), "tenantPolicyTemplate", "template", "matchRulePolicies")
	if err != nil {
		return err
	}
	return setMatchRulePolicyData(d, policy, templateID)
}
