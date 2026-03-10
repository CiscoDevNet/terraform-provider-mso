---
layout: "mso"
page_title: "MSO: mso_fabric_resource_policies_virtual_port_channel_interface"
sidebar_current: "docs-mso-datasource-fabric_resource_policies_virtual_port_channel_interface"
description: |-
  Data source for Virtual Port Channel (VPC) Interfaces on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_resource_policies_virtual_port_channel_interface #

Data source for Virtual Port Channel (VPC) Interfaces on Cisco Nexus Dashboard Orchestrator (NDO).

## GUI Information ##

* `Location` - Manage -> Fabric Templates -> Fabric Resource Policies -> Virtual Port Channel Interface

## Example Usage ##

```hcl
data "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
  template_id = mso_template.fabric_resource_template.id
  name        = "tf_vpc_if"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Resource template.
* `name` - (Required) The name of the Virtual Port Channel Interface.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Virtual Port Channel Interface.
* `description` - (Read-Only) The description of the Virtual Port Channel Interface.
* `node_1` - (Read-Only) The first node ID.
* `node_2` - (Read-Only) The second node ID.
* `node_1_interfaces` - (Read-Only) List of interface IDs (or ranges) for node 1.
* `node_2_interfaces` - (Read-Only) List of interface IDs (or ranges) for node 2.
* `interface_policy_group_uuid` - (Read-Only) UUID of the Port Channel Interface Policy Group.
* `interface_descriptions` - (Read-Only) List of interface description entries.
  * `node` - (Read-Only) Node ID.
  * `interface` - (Read-Only) Interface ID.
  * `description` - (Read-Only) Interface description.
