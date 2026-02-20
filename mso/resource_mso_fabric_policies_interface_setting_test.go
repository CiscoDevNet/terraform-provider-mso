package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOFabricPoliciesInterfaceSettingPhysicalResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Fabric Policies Interface Setting with type physical") },
				Config:    testAccMSOFabricPoliciesInterfaceSettingPhysicalConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "name", msoFabricPolicyTemplateInterfaceSettingName+"_physical"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "auto_negotiation", "on"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "cdp_admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "lldp_receive_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "llfc_receive_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "llfc_transmit_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "lldp_transmit_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "pfc_admin_state", "auto"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "l2_interface_qinq", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "l2_interface_reflective_relay", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "vlan_scope", "global"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "stp_bpdu_filter", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "stp_bpdu_guard", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_strict_mode", "off"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_initial_delay_time", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_transmission_frequency_sec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_transmission_frequency_msec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_grace_period_sec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_grace_period_msec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_bring_up_delay", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_debounce_interval", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_fec", "inherit"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "speed", "inherit"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "type", "physical"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "template_id"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "domains_uuid.#", "1"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "synce_uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "access_macsec_policy_uuid"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Fabric Policies Interface Setting with type physical") },
				ResourceName:      "mso_fabric_policies_interface_setting." + msoFabricPolicyTemplateInterfaceSettingName + "_physical",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Fabric Policies Interface Setting with type physical") },
				Config:    testAccMSOFabricPoliciesInterfaceSettingPhysicalConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "name", msoFabricPolicyTemplateInterfaceSettingName+"_physical_updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "description", "Physical interface for testing - updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "auto_negotiation", "off"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "cdp_admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "lldp_receive_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "llfc_receive_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "llfc_transmit_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "lldp_transmit_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "pfc_admin_state", "on"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "l2_interface_qinq", "core_port"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "l2_interface_reflective_relay", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "vlan_scope", "portlocal"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "stp_bpdu_filter", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "stp_bpdu_guard", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_strict_mode", "off"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_initial_delay_time", "360"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_transmission_frequency_sec", "5"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_transmission_frequency_msec", "500"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_grace_period_sec", "10"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "mcp_grace_period_msec", "500"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_bring_up_delay", "500"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_debounce_interval", "200"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "link_level_fec", "cl91_rs_fec"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "speed", "25G"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "type", "physical"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "template_id"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "synce_uuid", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_physical", "access_macsec_policy_uuid", ""),
				),
			},
		},
	})
}

func TestAccMSOFabricPoliciesInterfaceSettingPortChannelResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Fabric Policies Interface Setting with type portchannel") },
				Config:    testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "name", msoFabricPolicyTemplateInterfaceSettingName+"_portchannel"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "auto_negotiation", "on"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "cdp_admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "lldp_receive_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "llfc_receive_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "llfc_transmit_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "lldp_transmit_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "pfc_admin_state", "auto"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "l2_interface_qinq", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "l2_interface_reflective_relay", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "vlan_scope", "global"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "stp_bpdu_filter", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "stp_bpdu_guard", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_strict_mode", "off"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_initial_delay_time", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_transmission_frequency_msec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_transmission_frequency_sec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_grace_period_sec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_grace_period_msec", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_bring_up_delay", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_debounce_interval", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_fec", "inherit"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "template_id"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "load_balance_hashing", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_max_links", "16"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_min_links", "1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_mode", "static_channel_mode_on"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "speed", "inherit"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "type", "portchannel"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Fabric Policies Interface Setting with type portchannel") },
				ResourceName:      "mso_fabric_policies_interface_setting." + msoFabricPolicyTemplateInterfaceSettingName + "_portchannel",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Fabric Policies Interface Setting with type portchannel") },
				Config:    testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "name", msoFabricPolicyTemplateInterfaceSettingName+"_portchannel_updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "description", "Port-Channel interface for testing - updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "auto_negotiation", "on_enforce"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "cdp_admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "lldp_receive_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "llfc_receive_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "llfc_transmit_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "lldp_transmit_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "pfc_admin_state", "off"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "l2_interface_qinq", "double_q_tag_port"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "l2_interface_reflective_relay", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "vlan_scope", "portlocal"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "stp_bpdu_filter", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "stp_bpdu_guard", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_strict_mode", "on"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_initial_delay_time", "100"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_transmission_frequency_sec", "250"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_transmission_frequency_msec", "980"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_grace_period_sec", "230"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "mcp_grace_period_msec", "999"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_bring_up_delay", "200"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_debounce_interval", "100"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "link_level_fec", "ieee_rs_fec"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "speed", "100G"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "type", "portchannel"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "template_id"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "load_balance_hashing", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_max_links", "20"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_min_links", "6"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "port_channel_mode", "mac_pinning_physical_nic_load"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "controls.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "domains_uuid.#", "1"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "synce_uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_interface_setting."+msoFabricPolicyTemplateInterfaceSettingName+"_portchannel", "access_macsec_policy_uuid"),
				),
			},
		},
	})
}

func testAccMSOFabricPoliciesInterfaceSettingPhysicalConfigCreate() string {
	return fmt.Sprintf(`
%[1]s
%[4]s
resource "mso_fabric_policies_interface_setting" "%[2]s_physical" {
	template_id               = mso_template.%[3]s.id
	type                      = "physical"
	name                      = "%[2]s_physical"
	synce_uuid                = mso_fabric_policies_synce_interface_policy.%[5]s.uuid
	access_macsec_policy_uuid = mso_fabric_policies_macsec_policy.%[6]s.uuid
	domains_uuid              = [ mso_fabric_policies_l3_domain.%[7]s.uuid ]
}
`, testFabricPolicyTemplateConfig(), msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName, testFabricPolicyTemplateL3DomainConfig()+testFabricPolicyTemplateSyncEInterfacePolicyConfig()+testFabricPolicyTemplateMacsecPolicyConfig(), msoFabricPolicyTemplateSyncEInterfacePolicyName, msoFabricPolicyTemplateMacsecPolicyName, msoFabricPolicyTemplateL3DomainName)
}

func testAccMSOFabricPoliciesInterfaceSettingPhysicalConfigUpdate() string {
	return fmt.Sprintf(`
%[1]s
%[4]s
resource "mso_fabric_policies_interface_setting" "%[2]s_physical" {
	template_id                     = mso_template.%[3]s.id
	type                            = "physical"
	name                            = "%[2]s_physical_updated"
	auto_negotiation                = "off"
	cdp_admin_state                 = "disabled"
	description                     = "Physical interface for testing - updated"
	l2_interface_qinq               = "core_port"
	l2_interface_reflective_relay   = "enabled"
	link_level_bring_up_delay       = 500
	link_level_debounce_interval    = 200
	link_level_fec                  = "cl91_rs_fec"
	lldp_receive_state              = "disabled"
	lldp_transmit_state             = "disabled"
	llfc_receive_state              = "enabled"
	llfc_transmit_state             = "enabled"
	mcp_admin_state                 = "disabled"
	mcp_grace_period_msec           = 500
	mcp_grace_period_sec            = 10
	mcp_initial_delay_time          = 360
	mcp_strict_mode                 = "off"
	mcp_transmission_frequency_msec = 500
	mcp_transmission_frequency_sec  = 5
	pfc_admin_state                 = "on"
	speed                           = "25G"
	stp_bpdu_filter                 = "enabled"
	stp_bpdu_guard                  = "enabled"
	vlan_scope                      = "portlocal"
	synce_uuid                      = ""
	access_macsec_policy_uuid       = ""
	domains_uuid                    = []
}
`, testFabricPolicyTemplateConfig(), msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName, testFabricPolicyTemplateL3DomainConfig()+testFabricPolicyTemplateSyncEInterfacePolicyConfig()+testFabricPolicyTemplateMacsecPolicyConfig())
}

func testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigCreate() string {
	return fmt.Sprintf(`
%[1]s

resource "mso_fabric_policies_interface_setting" "%[2]s_portchannel" {
	template_id = mso_template.%[3]s.id
	type        = "portchannel"
	name        = "%[2]s_portchannel"
}
`, testFabricPolicyTemplateConfig(), msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName)
}

func testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigUpdate() string {
	return fmt.Sprintf(`
%[1]s

%[4]s

resource "mso_fabric_policies_interface_setting" "%[2]s_portchannel" {
	template_id                     = mso_template.%[3]s.id
	type                            = "portchannel"
	name                            = "%[2]s_portchannel_updated"
	auto_negotiation                = "on_enforce"
	cdp_admin_state                 = "enabled"
	controls                        = [ "graceful_conv", "susp_individual" ]
	description                     = "Port-Channel interface for testing - updated"
	l2_interface_qinq               = "double_q_tag_port"
	l2_interface_reflective_relay   = "enabled"
	link_level_bring_up_delay       = 200
	link_level_debounce_interval    = 100
	link_level_fec                  = "ieee_rs_fec"
	lldp_receive_state              = "disabled"
	lldp_transmit_state             = "disabled"
	llfc_receive_state              = "enabled"
	llfc_transmit_state             = "enabled"
	load_balance_hashing            = null
	mcp_admin_state                 = "enabled"
	mcp_grace_period_msec           = 999
	mcp_grace_period_sec            = 230
	mcp_initial_delay_time          = 100
	mcp_strict_mode                 = "on"
	mcp_transmission_frequency_msec = 980
	mcp_transmission_frequency_sec  = 250
	pfc_admin_state                 = "off"
	port_channel_max_links          = 20
	port_channel_min_links          = 6
	port_channel_mode               = "mac_pinning_physical_nic_load"
	speed                           = "100G"
	stp_bpdu_filter                 = "enabled"
	stp_bpdu_guard                  = "enabled"
	vlan_scope                      = "portlocal"
	synce_uuid                      = mso_fabric_policies_synce_interface_policy.%[5]s.uuid
	access_macsec_policy_uuid       = mso_fabric_policies_macsec_policy.%[6]s.uuid
	domains_uuid                    = [ mso_fabric_policies_l3_domain.%[7]s.uuid ]
}
`, testFabricPolicyTemplateConfig(), msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName, testFabricPolicyTemplateL3DomainConfig()+testFabricPolicyTemplateSyncEInterfacePolicyConfig()+testFabricPolicyTemplateMacsecPolicyConfig(), msoFabricPolicyTemplateSyncEInterfacePolicyName, msoFabricPolicyTemplateMacsecPolicyName, msoFabricPolicyTemplateL3DomainName)
}
