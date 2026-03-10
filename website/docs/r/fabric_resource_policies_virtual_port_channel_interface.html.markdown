---
layout: "mso"
page_title: "MSO: mso_fabric_resource_policies_virtual_port_channel_interface"
sidebar_current: "docs-mso-resource-fabric_resource_policies_virtual_port_channel_interface"
description: |-
  Manages Virtual Port Channel (VPC) Interfaces on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_resource_policies_virtual_port_channel_interface #

Manages Virtual Port Channel (VPC) Interfaces on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Templates -> Fabric Resource Policies -> Virtual Port Channel Interface

## Example Usage ##

```hcl
resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
  template_id = mso_template.fabric_resource_template.id
  name        = "tf_vpc_if"
  description = "Example VPC Interface"
  node_1 = "101"
  node_2 = "102"
  node_1_interfaces = ["1/1", "1/10-11"]
  node_2_interfaces = ["1/2"]
  interface_policy_group_uuid = mso_fabric_policies_interface_setting.port_channel_interface.uuid

  interface_descriptions {
    node        = "101"
    interface   = "1/1"
    description = "Terraform example interface description"
  }

  interface_descriptions {
    node        = "102"
    interface   = "1/10"
    description = "Terraform example interface description"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Resource template.
* `name` - (Required) The name of the Virtual Port Channel Interface.
* `description` - (Optional) The description of the Virtual Port Channel Interface.
* `node_1` - (Required) The first node ID.
* `node_2` - (Required) The second node ID.
* `node_1_interfaces` - (Required) List of interface IDs (or ranges) for node 1.
* `node_2_interfaces` - (Optional) List of interface IDs (or ranges) for node 2.
* `interface_policy_group_uuid` - (Required) UUID of the Port Channel Interface Policy Group.
* `interface_descriptions` - (Optional) List of interface description entries. The provided list replaces the existing list in NDO.
  * `node` - (Required) Node ID.
  * `interface` - (Required) Interface ID.
  * `description` - (Optional) Interface description.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Virtual Port Channel Interface.
* `id` - (Read-Only) The unique Terraform identifier of the Virtual Port Channel Interface.

## Importing ##

An existing MSO Virtual Port Channel Interface can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if templateId/{template_id}/virtualPortChannelInterface/{name}
```
