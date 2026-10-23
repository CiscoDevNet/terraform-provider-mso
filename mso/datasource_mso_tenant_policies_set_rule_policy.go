package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceMSOSetRulePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOSetRulePolicyRead,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Set Rule Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the Set Rule Policy.",
				Computed:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the Set Rule Policy.",
				Computed:    true,
			},
			"set_route_tag": {
				Type:        schema.TypeList,
				Description: "The route tag action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeInt,
							Description: "The action value. Range: 0-999999999.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"set_weight": {
				Type:        schema.TypeList,
				Description: "The weight action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeInt,
							Description: "The action value. Range: 0-9999.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"set_preference": {
				Type:        schema.TypeList,
				Description: "The preference action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeInt,
							Description: "The action value. Range: 0-999999999.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"set_metric": {
				Type:        schema.TypeInt,
				Description: "The metric action. Range: 1-999999999.",
				Computed:    true,
			},
			"set_next_hop": {
				Type:        schema.TypeString,
				Description: "The next-hop IPv4 or IPv6 address.",
				Computed:    true,
			},
			"set_metric_type": {
				Type:        schema.TypeString,
				Description: "The metric type. Valid values: none, ospf-type1, ospf-type2.",
				Computed:    true,
			},
			"set_next_hop_propagate": {
				Type:        schema.TypeBool,
				Description: "Enable next-hop propagation.",
				Computed:    true,
			},
			"set_multi_path": {
				Type:        schema.TypeBool,
				Description: "Enable multipath. Requires next-hop propagation.",
				Computed:    true,
			},
			"set_community": {
				Type:        schema.TypeList,
				Description: "The community action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"community": {
							Type:        schema.TypeString,
							Description: "The community value. Required for append and replace.",
							Computed:    true,
						},
						"criteria": {
							Type:        schema.TypeString,
							Description: "The community criteria. Valid values: none, append, replace.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"set_additional_community": {
				Type:        schema.TypeSet,
				Description: "The community values to append.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"set_dampening": {
				Type:        schema.TypeList,
				Description: "The route dampening action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"half_life": {
							Type:        schema.TypeInt,
							Description: "The half-life time in minutes. Range: 1-60.",
							Computed:    true,
						},
						"reuse_limit": {
							Type:        schema.TypeInt,
							Description: "The reuse limit. Range: 1-20000.",
							Computed:    true,
						},
						"suppress_limit": {
							Type:        schema.TypeInt,
							Description: "The suppress limit. Range: 1-20000.",
							Computed:    true,
						},
						"max_suppress_time": {
							Type:        schema.TypeInt,
							Description: "The maximum suppress time in minutes. Range: 1-255.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"set_as_path": {
				Type:        schema.TypeList,
				Description: "The AS-path prepend actions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prepend": {
							Type:        schema.TypeList,
							Description: "The action to prepend the configured AS numbers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_asn": {
										Type:        schema.TypeSet,
										Description: "The AS numbers to prepend. Order values determine the sequence.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"asn": {
													Type:        schema.TypeInt,
													Description: "The AS number. Range: 1-4294967295.",
													Computed:    true,
												},
												"order": {
													Type:        schema.TypeInt,
													Description: "The AS-path order. Range: 0-31.",
													Computed:    true,
												},
											},
										},
										Computed: true,
									},
								},
							},
							Computed: true,
						},
						"prepend_last_as": {
							Type:        schema.TypeList,
							Description: "The action to prepend the last AS number.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"asn_count": {
										Type:        schema.TypeInt,
										Description: "The number of times to prepend the last AS. Range: 1-10.",
										Computed:    true,
									},
								},
							},
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceMSOSetRulePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Set Rule Policy Data Source - Beginning Read")
	templateID := d.Get("template_id").(string)
	response, err := m.(*client.Client).GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateID))
	if err != nil {
		return err
	}
	policy, err := GetPolicyByName(response, d.Get("name").(string), "tenantPolicyTemplate", "template", "setRulePolicies")
	if err != nil {
		return err
	}
	return setSetRulePolicyData(d, policy, templateID)
}
