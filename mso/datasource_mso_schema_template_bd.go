package mso

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceMSOTemplateBD() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOTemplateBDRead,

		SchemaVersion: version,

		Schema: (map[string]*schema.Schema{
			"schema_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1000),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"intersite_bum_traffic": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"optimize_wan_bandwidth": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"layer2_stretch": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"layer3_multicast": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"layer2_unknown_unicast": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"unknown_multicast_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"multi_destination_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_unknown_multicast_flooding": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"arp_flooding": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"virtual_mac_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"unicast_routing": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_schema_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dhcp_policy": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Configure dhcp policy in versions before NDO 3.2",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"dhcp_policies": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Configure dhcp policies in versions NDO 3.2 and higher",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dhcp_option_policy_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_option_policy_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func dataSourceMSOTemplateBDRead(d *schema.ResourceData, m interface{}) error {
	return resourceMSOTemplateBDRead(d, m)
}
