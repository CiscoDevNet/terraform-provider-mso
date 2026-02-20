---
layout: "mso"
page_title: "MSO: mso_fabric_policies_mcp_global_policy"
sidebar_current: "docs-mso-data-source-fabric_policies_mcp_global_policy"
description: |-
  Data source for MCP (MisCabling Protocol) Global Policy.
---

# mso_fabric_policies_mcp_global_policy #

Data source for MCP (MisCabling Protocol) Global Policy. This data source is only supported on NDO v4.3 and later.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> MCP Global Policy

## Example Usage ##

```hcl
data "mso_fabric_policies_mcp_global_policy" "test" {
    template_id                       = mso_template.test.id
    name                              = "test"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the MCP Global Policy.

## Attribute Reference ##

* `description` - (Read-Only) The description of the MCP Global Policy.
* `admin_state` - (Read-Only) The administrative state of the MCP Global Policy.
* `enable_mcp_pdu_per_vlan` - (Read-Only) Enable MCP PDU per VLAN. This enables MCP to send packets on a per-EPG basis.
* `key` - (Read-Only) The key to uniquely identify the MCP packets within the fabric.
* `loop_detect_multiplication_factor` - (Read-Only) The number of MCP packets that will be received by the ACI fabric before the Loop Protection Action occurs.
* `port_disable_protection` - (Read-Only) Enable or disable port disable protection for the MCP Global Policy.
* `initial_delay_time` - (Read-Only) The time in seconds before MCP starts taking action. During this period, MCP will only generate syslog entries if a loop is detected, without taking protective action.
* `transmission_frequency_sec` - (Read-Only) The transmission frequency of the MCP packets in seconds.
* `transmission_frequency_msec` - (Read-Only) The transmission frequency of the MCP packets in milliseconds.
* `uuid` - (Read-Only) The NDO UUID of the MCP Global Policy.
* `id` - (Read-Only) The unique Terraform identifier of the MCP Global Policy.
