---
layout: "mso"
page_title: "MSO: mso_fabric_policies_interface_setting"
sidebar_current: "docs-mso-data-source-fabric_policies_interface_setting"
description: |-
  Data source for Interface Settings on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_interface_setting #

Data source for Interface Settings on Cisco Nexus Dashboard Orchestrator (NDO). This resource is only supported NDO v4.3 and later.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> Interface Settings

## Example Usage ##

```hcl
data "mso_fabric_policies_interface_setting" "test" {
    template_id                     = mso_template.test.id
    name                            = "test"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the Interface Setting.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Interface Setting.
* `id` - (Read-Only) The unique Terraform identifier of the Interface Setting.
* `type` - (Read-Only) The type of the Interface Setting.
* `description` - (Read-Only) The description of the Interface Setting.
* `cdp_admin_state` - (Read-Only) The administrative state of the CDP (Cisco Discovery Protocol) protocol.
* `lldp_receive_state` - (Read-Only) The receive state for the Link Layer Discovery Protocol (LLDP).
* `lldp_transmit_state` - (Read-Only) The transmit state for the LLDP.
* `llfc_receive_state` - (Read-Only) The receive state for the LLFC (Link Level Flow Control).
* `llfc_transmit_state` - (Read-Only) The transmit state for the LLFC.
* `pfc_admin_state` - (Read-Only) The administrative state of the PFC (Priority Flow Control).
* `l2_interface_qinq` - (Read-Only) The QinQ mode for the L2 interfaces.
* `l2_interface_reflective_relay` - (Read-Only) Enable or disable reflective relay for the L2 interfaces.
* `vlan_scope` - (Read-Only) The VLAN scope for the interface.
* `stp_bpdu_filter` - (Read-Only) Enable or disable BPDU (Bridge Protocol Data Unit) filter for the STP (Spanning Tree Protocol).
* `stp_bpdu_guard` - (Read-Only) Enable or disable BPDU guard for the STP.
* `speed` - (Read-Only) The speed of the interface.
* `auto_negotiation` - (Read-Only) The auto-negotiation state of the interface.
* `link_level_bring_up_delay` - (Read-Only) The bring-up delay time in milliseconds.
* `link_level_debounce_interval` - (Read-Only) The debounce interval in milliseconds.
* `link_level_fec` - (Read-Only) The FEC (Forward Error Correction) mode.
* `mcp_admin_state` - (Read-Only) The administrative state of the MCP (Missed Class Protocol).
* `mcp_strict_mode` - (Read-Only) Enable or disable strict mode for the MCP.
* `mcp_initial_delay_time` - (Read-Only) The initial delay time for the MCP in seconds.
* `mcp_transmission_frequency_sec` - (Read-Only) The MCP transmission frequency in seconds.
* `mcp_transmission_frequency_msec` - (Read-Only) The MCP transmission frequency in milliseconds.
* `mcp_grace_period_sec` - (Read-Only) The MCP grace period in seconds.
* `mcp_grace_period_msec` - (Read-Only) The MCP grace period in milliseconds.
* `port_channel_mode` - (Read-Only) The port channel mode for the Interface Setting.
* `controls` - (Read-Only) A list of port channel controls.
* `port_channel_min_links` - (Read-Only) The minimum number of active links for the port channel.
* `port_channel_max_links` - (Read-Only) The maximum number of links for the port channel.
* `load_balance_hashing` - (Read-Only) The load balancing hashing algorithm for the Interface Setting.
* `synce_uuid` - (Read-Only) The UUID of the SyncE Interface Policy to associate with this Interface Setting.
* `domain_uuids` - (Read-Only) A list of UUIDs of the L3 Domains to associate with this Interface Setting.
* `access_macsec_policy_uuid` - (Read-Only) The UUID of the access‑type MACsec policy to be associated with this Interface Setting.
