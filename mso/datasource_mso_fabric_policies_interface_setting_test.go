package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOFabricPoliciesInterfaceSettingPhysicalDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Data Source Interface Setting with type physical") },
				Config:    testAccMSOFabricPoliciesInterfaceSettingPhysicalDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "name", msoFabricPolicyTemplateInterfaceSettingName+"_physical_updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "description", "Physical interface for testing - updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "auto_negotiation", "off"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "cdp_admin_state", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "lldp_receive_state", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "llfc_receive_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "llfc_transmit_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "lldp_transmit_state", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "pfc_admin_state", "on"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "l2_interface_qinq", "core_port"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "l2_interface_reflective_relay", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "vlan_scope", "portlocal"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "stp_bpdu_filter", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "stp_bpdu_guard", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_admin_state", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_strict_mode", "off"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_initial_delay_time", "360"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_transmission_frequency_sec", "5"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_transmission_frequency_msec", "500"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_grace_period_sec", "10"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_grace_period_msec", "500"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_bring_up_delay", "500"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_debounce_interval", "200"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_fec", "cl91_rs_fec"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "speed", "25G"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "type", "physical"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "template_id"),
				),
			},
		},
	})
}

func TestAccMSOFabricPoliciesInterfaceSettingPortChannelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				PreConfig: func() { fmt.Println("Test: Data Source Interface Setting with type portchannel") },
				Config:    testAccMSOFabricPoliciesInterfaceSettingPortChannelDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "name", msoFabricPolicyTemplateInterfaceSettingName+"_portchannel_updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "description", "Port-Channel interface for testing - updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "auto_negotiation", "on_enforce"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "cdp_admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "lldp_receive_state", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "llfc_receive_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "llfc_transmit_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "lldp_transmit_state", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "pfc_admin_state", "off"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "l2_interface_qinq", "double_q_tag_port"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "l2_interface_reflective_relay", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "vlan_scope", "portlocal"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "stp_bpdu_filter", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "stp_bpdu_guard", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_strict_mode", "on"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_initial_delay_time", "100"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_transmission_frequency_sec", "250"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_transmission_frequency_msec", "980"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_grace_period_sec", "230"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_grace_period_msec", "999"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_bring_up_delay", "200"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_debounce_interval", "100"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_fec", "ieee_rs_fec"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "speed", "100G"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "type", "portchannel"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "template_id"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "load_balance_hashing", ""),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_max_links", "20"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_min_links", "6"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_mode", "mac_pinning_physical_nic_load"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "controls.#", "2"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "domain_uuids.#", "1"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "synce_uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "access_macsec_policy_uuid"),
				),
			},
		},
	})
}

func testAccMSOFabricPoliciesInterfaceSettingPhysicalDataSource() string {
	return fmt.Sprintf(`%[1]s
    data "mso_fabric_policies_interface_setting" "%[2]s_physical" {
		template_id                        = mso_template.%[3]s.id
		name                               = mso_fabric_policies_interface_setting.%[2]s_physical.name
	}`, testAccMSOFabricPoliciesInterfaceSettingPhysicalConfigUpdate(), msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName)
}

func testAccMSOFabricPoliciesInterfaceSettingPortChannelDataSource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_fabric_policies_interface_setting" "%[2]s_portchannel" {
		template_id                        = mso_template.%[3]s.id
		name                               = mso_fabric_policies_interface_setting.%[2]s_portchannel.name
	}`, testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigUpdate(), msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName)
}
