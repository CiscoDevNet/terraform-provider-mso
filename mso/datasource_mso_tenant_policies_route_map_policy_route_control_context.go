package mso

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSORouteMapPolicyContext() *schema.Resource {
	return &schema.Resource{
		Read: resourceMSORouteMapPolicyContextRead,

		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"match_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"set_rule_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
