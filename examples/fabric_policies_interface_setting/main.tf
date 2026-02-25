terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO pwd>
  url      = "" # <MSO URL>
  insecure = true
}

resource "mso_template" "fabric_policy_template" {
	template_name = "fabric_policy_template"
	template_type = "fabric_policy"
}

resource "mso_fabric_policies_l3_domain" "l3_domain" {
  template_id    = mso_template.fabric_policy_template.id
  name           = "l3_domain"
}

resource "mso_fabric_policies_synce_interface_policy" "synce_interface_policy" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "synce_interface_policy"
}

resource "mso_fabric_policies_macsec_policy" "macsec_policy" {
  template_id            = mso_template.fabric_policy_template.id
  name                   = "macsec_policy"
  interface_type         = "access"
  cipher_suite           = "256GcmAes"
  window_size            = 128
  security_policy        = "shouldSecure"
  sak_expire_time        = 60
  confidentiality_offset = "offset30"
  key_server_priority    = 8
}

resource "mso_fabric_policies_interface_setting" "physical_interface" {
	template_id                     = mso_template.fabric_policy_template.id
	type                            = "physical"
	name                            = "physical_interface"
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
}

resource "mso_fabric_policies_interface_setting" "portchannel_interface" {
	template_id                     = mso_template.fabric_policy_template.id
	type                            = "portchannel"
	name                            = "portchannel_interface"
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
	synce_uuid                      = mso_fabric_policies_synce_interface_policy.synce_policy.uuid
	access_macsec_policy_uuid       = mso_fabric_policies_macsec_policy.macsec_policy.uuid
	domain_uuids                    = [ mso_fabric_policies_l3_domain.l3_domain.uuid ]
}
