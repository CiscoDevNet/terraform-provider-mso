package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceMSOInterfaceSetting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMSOInterfaceSettingRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cdp_admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lldp_receive_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lldp_transmit_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"llfc_receive_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"llfc_transmit_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pfc_admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l2_interface_qinq": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l2_interface_reflective_relay": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stp_bpdu_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stp_bpdu_guard": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"speed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_negotiation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link_level_bring_up_delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"link_level_debounce_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"link_level_fec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mcp_admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mcp_strict_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mcp_initial_delay_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mcp_transmission_frequency_sec": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mcp_transmission_frequency_msec": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mcp_grace_period_sec": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mcp_grace_period_msec": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_channel_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"port_channel_min_links": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_channel_max_links": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_balance_hashing": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"synce_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domains_uuid": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"access_macsec_policy_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMSOInterfaceSettingRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] MSO Interface Setting Data Source - Beginning Read")
	msoClient := m.(*client.Client)

	templateId := d.Get("template_id").(string)
	name := d.Get("name").(string)

	response, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/templates/%s", templateId))
	if err != nil {
		return err
	}

	policy, err := GetPolicyByName(response, name, "fabricPolicyTemplate", "template", "interfacePolicyGroups")
	if err != nil {
		return err
	}

	err = setInterfaceSettingData(d, policy, templateId)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] MSO Interface Setting Data Source - Read Complete: %s", d.Id())
	return nil
}
