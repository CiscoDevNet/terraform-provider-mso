package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOL3OutNodeRoutingPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOL3OutNodeRoutingPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the L3Out Node Routing Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the L3Out Node Routing Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the L3Out Node Routing Policy.",
			},
			"as_path_multipath_relax": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "BGP Best Path Control - AS path multipath relax.",
			},
			"bfd_multi_hop_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "BFD multi-hop configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Administrative state.",
						},
						"detection_multiplier": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detection multiplier.",
						},
						"min_receive_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum receive interval in microseconds.",
						},
						"min_transmit_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum transmit interval in microseconds.",
						},
					},
				},
			},
			"bgp_node_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "BGP node configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"graceful_restart_helper": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Graceful restart helper mode.",
						},
						"keep_alive_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP keepalive interval in seconds.",
						},
						"hold_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP hold interval in seconds.",
						},
						"stale_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP stale interval in seconds.",
						},
						"max_as_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum AS path limit.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOL3OutNodeRoutingPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "l3OutNodePolGroups")
	if err != nil {
		return err
	}

	setL3OutNodeRoutingPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO L3Out Node Routing Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
