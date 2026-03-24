---
layout: "mso"
page_title: "MSO: mso_fabric_resource_policies_physical_interface"
sidebar_current: "docs-mso-resource-fabric_resource_policies_physical_interface"
description: |-
  Manages Fabric Resource Policies Physical Interface on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_resource_policies_physical_interface #

Manages Fabric Resource Policies Physical Interface on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.1 or higher.

## GUI Information ##

* `Location`: Manage -> Fabric Resource Template -> Fabric Resource Policies -> Physical Interfaces

## Example Usage ##

### Physical Interface with Interface Policy

```hcl
resource "mso_fabric_resource_policies_physical_interface" "physical" {
  template_id           = mso_template.fabric_resource_template.id
  name                  = "physical_interface"
  description           = "Terraform test Physical Interface"
  nodes                 = ["101", "102"]
  interfaces            = ["1/1", "1/2"]
  interface_policy_uuid = mso_fabric_policies_interface_setting.physical_interface.id
  interface_descriptions {
    interface   = "1/1"
    description = "Interface Description 1/1"
  }
  interface_descriptions {
    interface   = "1/2"
    description = "Interface Description 1/2"
  }
}
```

### Physical Interface with Breakout Mode

```hcl
resource "mso_fabric_resource_policies_physical_interface" "breakout" {
  template_id   = mso_template.fabric_resource_template.id
  name          = "breakout_mode_physical_interface"
  description   = "Terraform test Physical Interface with breakout mode"
  nodes         = ["101", "102"]
  interfaces    = ["1/1", "1/2"]
  breakout_mode = "4x100G"
  interface_descriptions {
    interface   = "1/1"
    description = "Interface Description 1/1"
  }
  interface_descriptions {
    interface   = "1/2"
    description = "Interface Description 1/2"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the Fabric Resource template.
* `name` - (Required) The name of the Physical Interface Policy.
* `description` - (Optional) The description of the Physical Interface Policy.
* `nodes` - (Required) A list of node IDs where this Physical Interface Policy will be applied. This is required when creating or updating a Physical Interface Policy.
* `interfaces` - (Required) A list of interfaces where this Physical Interface Policy will be applied. This is required when creating or updating a Physical Interface Policy.
* `interface_policy_uuid` - (Optional) The UUID of the (physical) interface settings policy to associate with the Physical Interface Policy. This policy will be applied to every interface listed in the `interfaces` attribute. Either `interface_policy_uuid` or `breakout_mode` must be specified when creating or updating a Physical Interface Policy.
* `breakout_mode` - (Optional) The breakout mode of the Physical Interface Policy. Valid values are `4x100G`, `4x25G`, and `4x10G`. Either `interface_policy_uuid` or `breakout_mode` must be specified when creating or updating a Physical Interface Policy.
* `interface_descriptions` - (Optional) A list of interface descriptions of the Physical Interface Policy.
  * `interface` - (Required) The interface ID of the member interface. Must match an interface defined in the `interfaces` attribute.
  * `description` - (Optional) The description of the member interface.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Physical Interface Policy.
* `uuid` - (Read-Only) The NDO UUID of the Physical Interface Policy.
* `policy_group_type` - (Read-Only) The policy group type of the Physical Interface Policy.

## Importing ##

An existing MSO Fabric Resource Policies Physical Interface Policy can be [imported][docs-import] into this resource via its ID, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```bash
terraform import mso_fabric_resource_policies_physical_interface.example templateId/{template_id}/PhysicalInterface/{name}
```