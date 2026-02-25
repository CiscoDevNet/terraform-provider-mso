---
layout: "mso"
page_title: "MSO: mso_fabric_policies_interface_setting"
sidebar_current: "docs-mso-resource-fabric_policies_interface_setting"
description: |-
  Manages Interface Settings on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_interface_setting #

Manages Interface Settings on Cisco Nexus Dashboard Orchestrator (NDO). This resource is only supported NDO v4.3 and later.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> Interface Settings

## Example Usage ##

```hcl
resource "mso_fabric_policies_interface_setting" "test" {
    template_id                     = mso_template.test.id
    type                            = "physical"
    name                            = "test"
    auto_negotiation                = "off"
    cdp_admin_state                 = "disabled"
    description                     = "Test"
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
    synce_uuid                      = mso_fabric_policies_synce_interface_policy.test.uuid
    access_macsec_policy_uuid       = mso_fabric_policies_macsec_policy.test.uuid
    domain_uuids                    = [ mso_fabric_policies_l3_domain.test.uuid ]
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `type` - (Required) The type of the Interface Setting. Allowed values are `physical` or `portchannel`.
* `name` - (Required) The name of the Interface Setting.
* `description` - (Optional) The description of the Interface Setting.
* `cdp_admin_state` - (Optional) The administrative state of the CDP (Cisco Discovery Protocol) protocol. Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `lldp_receive_state` - (Optional) The receive state for the Link Layer Discovery Protocol (LLDP). Allowed values are `enabled` or `disabled`. Defaults to `enabled` when unset during creation.
* `lldp_transmit_state` - (Optional) The transmit state for the LLDP. Allowed values are `enabled` or `disabled`. Defaults to `enabled` when unset during creation.
* `llfc_receive_state` - (Optional) The receive state for the LLFC (Link Level Flow Control). Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `llfc_transmit_state` - (Optional) The transmit state for the LLFC. Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `pfc_admin_state` - (Optional) The administrative state of the PFC (Priority Flow Control). Allowed values are `off`, `on`, or `auto`. Defaults to `auto` when unset during creation.
* `l2_interface_qinq` - (Optional) The QinQ mode for the L2 interfaces. Allowed values are `double_q_tag_port`, `core_port`, `edge_port`, or `disabled`. Defaults to `disabled` when unset during creation.
* `l2_interface_reflective_relay` - (Optional) Enable or disable reflective relay for the L2 interfaces. Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `vlan_scope` - (Optional) The VLAN scope for the interface. Allowed values are `global` or `portlocal`. Defaults to `global` when unset during creation.
* `stp_bpdu_filter` - (Optional) Enable or disable BPDU (Bridge Protocol Data Unit) filter for the STP (Spanning Tree Protocol). Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `stp_bpdu_guard` - (Optional) Enable or disable BPDU guard for the STP. Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `speed` - (Optional) The speed of the interface. Allowed values are `100M`, `1G`, `10G`, `25G`, `40G`, `50G`, `100G`, `200G`, `400G`, or `inherit`. Defaults to `inherit` when unset during creation.
* `auto_negotiation` - (Optional) The auto-negotiation state of the interface. Allowed values are `on`, `off`, or `on_enforce`. Defaults to `on` when unset during creation.
* `link_level_bring_up_delay` - (Optional) The bring-up delay time in milliseconds. Valid range: 0-10000. Defaults to 0 when unset during creation.
* `link_level_debounce_interval` - (Optional) The debounce interval in milliseconds. Valid range: 0-5000. Defaults to 0 when unset during creation.
* `link_level_fec` - (Optional) The FEC (Forward Error Correction) mode. Allowed values are `inherit`, `cl74_fc_fec`, `cl91_rs_fec`, `cons16_rs_fec`, `ieee_rs_fec`, `kp_fec`, or `disable_fec`. Defaults to `inherit` when unset during creation.
* `mcp_admin_state` - (Optional) The administrative state of the MCP (Missed Class Protocol). Allowed values are `enabled` or `disabled`. Defaults to `enabled` when unset during creation.
* `mcp_strict_mode` - (Optional) Enable or disable strict mode for the MCP. Allowed values are `on` or `off`. Defaults to `off` when unset during creation.
* `mcp_initial_delay_time` - (Optional) The initial delay time for the MCP in seconds. Valid range: 0-1800. Defaults to 0 when unset during creation.
* `mcp_transmission_frequency_sec` - (Optional) The MCP transmission frequency in seconds. Valid range: 0-300. Defaults to 0 when unset during creation.
* `mcp_transmission_frequency_msec` - (Optional) The MCP transmission frequency in milliseconds. Valid range: 0-999. Defaults to 0 when unset during creation.
* `mcp_grace_period_sec` - (Optional) The MCP grace period in seconds. Valid range: 0-300. Defaults to 0 when unset during creation.
* `mcp_grace_period_msec` - (Optional) The MCP grace period in milliseconds. Valid range: 0-999. Defaults to 0 when unset during creation.
* `port_channel_mode` - (Optional) The port channel mode for the Interface Setting. This applies only when the Interface Setting has `type=portchannel`. Allowed values are `lacp_active`, `lacp_passive`, `static_channel_mode_on`, `mac_pinning`, `mac_pinning_physical_nic_load`, or `use_explicit_failover_order`. Defaults to `static_channel_mode_on` when unset during creation.
* `controls` - (Optional) A list of port channel controls. Allowed values are `fast_sel_hot_stdby`, `graceful_conv`, `susp_individual`, `load_defer`, or `symmetric_hash`.
* `port_channel_min_links` - (Optional) The minimum number of active links for the port channel. This applies only when the Interface Setting has `type=portchannel`. Valid range: 1-16. Defaults to 1 when unset during creation.
* `port_channel_max_links` - (Optional) The maximum number of links for the port channel. This applies only when the Interface Setting has `type=portchannel`. Valid range: 1-64. Defaults to 16 when unset during creation.
* `load_balance_hashing` - (Optional) The load balancing hashing algorithm for the Interface Setting. This applies only when the Interface Setting has `type=portchannel`. Allowed values are `destination_ip`, `layer_4_destination_ip`, `layer_4_source_ip`, or `source_ip`.
* `synce_uuid` - (Optional) The UUID of the SyncE Interface Policy to associate with this Interface Setting.
* `domain_uuids` - (Optional) A list of UUIDs of the L3 Domains to associate with this Interface Setting.
* `access_macsec_policy_uuid` - (Optional) The UUID of the access‑type MACsec policy to be associated with this Interface Setting.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Interface Setting.
* `id` - (Read-Only) The unique Terraform identifier of the Interface Setting.

## Importing ##

An existing MSO Interface Setting can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_interface_setting.test templateId/{template_id}/InterfaceSetting/{name}
```
