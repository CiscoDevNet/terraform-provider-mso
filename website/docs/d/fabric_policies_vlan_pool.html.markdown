---
layout: "mso"
page_title: "MSO: mso_fabric_policies_vlan_pool"
sidebar_current: "docs-mso-data-source-fabric_policies_vlan_pool"
description: |-
  Manages VLAN Pools on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_policies_vlan_pool #

Manages VLAN Pools on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> VLAN Pool

## Example Usage ##

```hcl
data "mso_fabric_policies_vlan_pool" "vlan_pool" {
  template_id = mso_template.fabric_policy_template.id
  name        = "vlan_pool"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the VLAN Pool.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the VLAN Pool.
* `id` - (Read-Only) The unique Terraform identifier of the VLAN Pool.
* `description` - (Read-Only) The description of the VLAN Pool.
* `allocation_mode` - (Read-Only) The allocation mode of the VLAN Pool.
* `vlan_range` - (Read-Only) The list of encapsulation blocks, each defining a range of VLAN IDs.
  * `vlan_range.from` - (Read-Only) The starting VLAN ID of the encapsulation block.
  * `vlan_range.to` - (Read-Only) The ending VLAN ID of the encapsulation block.
  * `vlan_range.allocation_mode` - (Read-Only) The allocation mode of the encapsulation block.
