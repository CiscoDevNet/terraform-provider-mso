---
layout: "mso"
page_title: "MSO: mso_fabric_policies_physical_domain"
sidebar_current: "docs-mso-resource-fabric_policies_physical_domain"
description: |-
  Manages Physical Domains on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_physical_domain #

Manages Physical Domains on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> Physical Domain

## Example Usage ##

```hcl
resource "mso_fabric_policies_physical_domain" "physical_domain" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "physical_domain"
  description     = "Example description"
  vlan_pool_uuid  = mso_fabric_policies_vlan_pool.vlan_pool.uuid
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the Physical Domain.
* `description` - (Optional) The description of the Physical Domain.
* `vlan_pool_uuid` - (Optional) The NDO UUID of the VLAN Pool to associate with the Physical Domain.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Physical Domain.
* `id` - (Read-Only) The unique Terraform identifier of the Physical Domain.

## Importing ##

An existing MSO Physical Domain can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_physical_domain.physical_domain templateId/{template_id}/physicalDomain/{name}
```
