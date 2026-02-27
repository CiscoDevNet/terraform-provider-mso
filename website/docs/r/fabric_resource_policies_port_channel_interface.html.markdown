---
layout: "mso"
page_title: "MSO: mso_fabric_resource_policies_port_channel_interface"
sidebar_current: "docs-mso-resource-fabric_resource_policies_port_channel_interface"
description: |-
  Manages Fabric Resource Policies Port Channel Interface on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_resource_policies_port_channel_interface #

Manages Fabric Resource Policies Port Channel Interface on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.1 or higher.

## GUI Information ##

* `Location`: Manage -> Fabric Resource Template -> Fabric Resource Policies -> Port Channel Interfaces

## Example Usage ##

```hcl
resource "mso_fabric_resource_policies_port_channel_interface" "example" {
  template_id           = mso_template.fabric_resource_template.id
  name                  = "example"
  description           = "example description"
  node                  = "101"
  interfaces            = ["1/1", "1/2"]
  interface_policy_uuid = mso_fabric_policies_interface_setting.port_channel_interface.id
  interface_descriptions {
    interface   = "1/1"
    description = "1/1 description"
  }
  interface_descriptions {
    interface   = "1/2"
    description = "1/2 description"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Resource template.
* `name` - (Required) The name of the Port Channel Interface.
* `description` - (Optional) The description of the Port Channel Interface.
* `node` - (Required) The node ID of the Port Channel Interface. This is required when creating or updating a Port Channel Interface.
* `interfaces` - (Required) The member interfaces of the Port Channel. This is required when creating or updating a Port Channel Interface.
* `interface_policy_uuid` - (Required) The UUID of the Port Channel/Virtual Port Channel Interface Policy to associate with the Port Channel Interface. This is required when creating a Port Channel Interface.
* `interface_descriptions` - (Optional) A list of interface descriptions for the Port Channel member interfaces.
  * `interface` - (Required) The interface ID of the member interface. Must match an interface defined in the `interfaces` attribute.
  * `description` - (Optional) The description of the member interface.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Port Channel Interface.
* `uuid` - (Read-Only) The NDO UUID of the Port Channel Interface.

## Importing ##

An existing MSO Fabric Resource Policies Port Channel Interface can be [imported][docs-import] into this resource via its ID, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```bash
terraform import mso_fabric_resource_policies_port_channel_interface.example templateId/{template_id}/PortChannelInterface/{name}
```
