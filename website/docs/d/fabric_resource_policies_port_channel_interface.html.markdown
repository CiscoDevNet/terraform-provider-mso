---
layout: "mso"
page_title: "MSO: mso_fabric_resource_policies_port_channel_interface"
sidebar_current: "docs-mso-data-source-fabric_resource_policies_port_channel_interface"
description: |-
  Data source for Fabric Resource Policies Port Channel Interface on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_resource_policies_port_channel_interface #

Data source for Fabric Resource Policies Port Channel Interface on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.1 or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Resource Template -> Fabric Resource Policies -> Port Channel Interfaces

## Example Usage ##

```hcl
data "mso_fabric_resource_policies_port_channel_interface" "example" {
  template_id = mso_template.fabric_resource_template.id
  name        = "port_channel_interface_example"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Resource template.
* `name` - (Required) The name of the Port Channel Interface Policy.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Port Channel Interface Policy.
* `uuid` - (Read-Only) The NDO UUID of the Port Channel Interface Policy.
* `description` - (Read-Only) The description of the Port Channel Interface Policy.
* `node` - (Read-Only) The node ID of the Port Channel Interface Policy.
* `interfaces` - (Read-Only) A list of interfaces where this Port Channel Interface Policy will be applied.
* `interface_policy_group_uuid` - (Read-Only) The UUID of the (Port Channel/Virtual Port Channel) interface settings policy to associate with the Port Channel Interface Policy.
* `interface_descriptions` - (Read-Only) A list of interface descriptions of the Port Channel Interface Policy.
  * `interface` - (Read-Only) The interface ID of the member interface.
  * `description` - (Read-Only) The description of the member interface.
