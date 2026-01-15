package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOL3OutInterfaceRoutingPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOL3OutInterfaceRoutingPolicyRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the tenant policy template.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the L3Out Interface Routing Policy.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the L3Out Interface Routing Policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the L3Out Interface Routing Policy.",
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
			"bfd_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "BFD configuration.",
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
						"echo_receive_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Echo receive interval in microseconds.",
						},
						"echo_admin_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Echo administrative state.",
						},
						"interface_control": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Interface control.",
						},
					},
				},
			},
			"ospf_interface_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "OSPF interface configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network type.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "OSPF priority.",
						},
						"cost_of_interface": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "OSPF cost.",
						},
						"hello_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Hello interval in seconds.",
						},
						"dead_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dead interval in seconds.",
						},
						"retransmit_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retransmit interval in seconds.",
						},
						"transmit_delay": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Transmit delay in seconds.",
						},
						"advertise_subnet": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Advertise subnet.",
						},
						"bfd": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable BFD.",
						},
						"mtu_ignore": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Ignore MTU.",
						},
						"passive_participation": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Passive participation.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMSOL3OutInterfaceRoutingPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	policyName := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, policyName, "tenantPolicyTemplate", "template", "l3OutIntfPolGroups")
	if err != nil {
		return err
	}

	setL3OutInterfaceRoutingPolicyData(d, policy, templateId)
	log.Printf("[DEBUG] MSO L3Out Interface Routing Policy Data Source - Read Complete: %v", d.Id())
	return nil
}
