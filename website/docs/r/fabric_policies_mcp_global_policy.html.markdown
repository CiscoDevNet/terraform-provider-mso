---
layout: "mso"
page_title: "MSO: mso_fabric_policies_mcp_global_policy"
sidebar_current: "docs-mso-resource-fabric_policies_mcp_global_policy"
description: |-
  Manages MCP (MisCabling Protocol) Global Policy on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_mcp_global_policy #

Manages MCP (MisCabling Protocol) Global Policy on Cisco Nexus Dashboard Orchestrator (NDO). This resource is only supported on NDO v4.3 and later.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> MCP Global Policy

## Example Usage ##

```hcl
resource "mso_fabric_policies_mcp_global_policy" "test" {
    template_id                       = mso_template.test.id
    name                              = "test"
    description                       = "Test MCP Global Policy"
    admin_state                       = "enabled"
    enable_mcp_pdu_per_vlan           = "enabled"
    key                               = "test_key_123"
    loop_detect_multiplication_factor = 5
    port_disable_protection           = "enabled"
    initial_delay_time                = 360
    transmission_frequency_sec        = 5
    transmission_frequency_msec       = 500
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the MCP Global Policy.
* `description` - (Optional) The description of the MCP Global Policy.
* `admin_state` - (Optional) The administrative state of the MCP Global Policy. Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `enable_mcp_pdu_per_vlan` - (Optional) Enable MCP PDU per VLAN. This enables MCP to send packets on a per-EPG basis. Allowed values are `enabled` or `disabled`. Defaults to `disabled` when unset during creation.
* `key` - (Optional) The key to uniquely identify the MCP packets within the fabric. This must be provided when `admin_state` is set to `enabled`.
* `loop_detect_multiplication_factor` - (Optional) The number of MCP packets that will be received by the ACI fabric before the Loop Protection Action occurs. Valid range: 1-255. Defaults to `3` when unset during creation.
* `port_disable_protection` - (Optional) Enable or disable port disable protection for the MCP Global Policy. Allowed values are `enabled` or `disabled`. Defaults to `enabled` when unset during creation.
* `initial_delay_time` - (Optional) The time in seconds before MCP starts taking action. During this period, MCP will only generate syslog entries if a loop is detected, without taking protective action. Valid range: 0-1800. Defaults to `180` when unset during creation.
* `transmission_frequency_sec` - (Optional) The transmission frequency of the MCP packets in seconds. Valid range: 0-300. Defaults to `2` when unset during creation.
* `transmission_frequency_msec` - (Optional) The transmission frequency of the MCP packets in milliseconds. Valid range: 0-999. Defaults to `0` when unset during creation.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the MCP Global Policy.
* `id` - (Read-Only) The unique Terraform identifier of the MCP Global Policy.

## Importing ##

An existing MSO MCP Global Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_mcp_global_policy.test templateId/{template_id}/MCPGlobalPolicy/{name}
```