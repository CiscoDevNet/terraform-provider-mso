---
layout: "mso"
page_title: "MSO: mso_fabric_policies_vlan_pool"
sidebar_current: "docs-mso-resource-fabric_policies_vlan_pool"
description: |-
  Manages VLAN Pools on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_policies_vlan_pool #

Manages VLAN Pools on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> VLAN Pool

## Example Usage ##

```hcl
resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
  template_id = mso_template.fabric_policy_template.id
  name        = "vlan_pool"
  description = "Example description"
	vlan_range {
    from            = 200
    to              = 202
	}
	vlan_range {
    from            = 204
    to              = 209
	}
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the VLAN Pool.
* `description` - (Optional) The description of the VLAN Pool.
* `vlan_range` - (Optional) The list of encapsulation blocks, each defining a range of VLAN IDs. At least one must be set when creating a VLAN Pool.
  * `vlan_range.from` - (Required) The starting VLAN ID of the encapsulation block.
  * `vlan_range.to` - (Required) The ending VLAN ID of the encapsulation block.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the VLAN Pool.
* `id` - (Read-Only) The unique Terraform identifier of the VLAN Pool.

## Importing ##

An existing MSO VLAN Pool can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_vlan_pool.vlan_pool templateId/{template_id}/VlanPool/{name}
```
