---
layout: "mso"
page_title: "MSO: mso_fabric_resource_policies_physical_interface"
sidebar_current: "docs-mso-data-source-fabric_resource_policies_physical_interface"
description: |-
  Data source for Fabric Resource Policies Physical Interface on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_resource_policies_physical_interface #

Data source for Fabric Resource Policies Physical Interface on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.1 or higher.

## GUI Information ##

* `Location`: Manage -> Fabric Resource Template -> Fabric Resource Policies -> Physical Interfaces

## Example Usage ##

```hcl
data "mso_fabric_resource_policies_physical_interface" "example" {
  template_id = mso_template.fabric_resource_template.id
  name        = "physical_interface"
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the Fabric Resource template.
* `name` - (Required) The name of the Physical Interface Policy.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Physical Interface Policy.
* `uuid` - (Read-Only) The NDO UUID of the Physical Interface Policy.
* `description` - (Read-Only) The description of the Physical Interface Policy.
* `nodes` - (Read-Only) A list of node IDs where this Physical Interface Policy will be applied.
* `interfaces` - (Read-Only) A list of interfaces where this Physical Interface Policy will be applied.
* `interface_policy_uuid` - (Read-Only) The UUID of the (physical) interface settings policy to associate with the Physical Interface Policy. This policy will be applied to every interface listed in the `interfaces` attribute.
* `breakout_mode` - (Read-Only) The breakout mode of the Physical Interface Policy.
* `interface_descriptions` - (Read-Only) A list of interface descriptions of the Physical Interface Policy.
  * `interface` - (Read-Only) The interface ID of the member interface.
  * `description` - (Read-Only) The description of the member interface.
* `policy_group_type` - (Read-Only) The policy group type of the Physical Interface.
