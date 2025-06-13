---
layout: "mso"
page_title: "MSO: mso_fabric_policies_physical_domain"
sidebar_current: "docs-mso-data-source-fabric_policies_physical_domain"
description: |-
  Data source for Physical Domains on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_physical_domain #

Data source for Physical Domains on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> Physical Domain

## Example Usage ##

```hcl
data "mso_fabric_policies_physical_domain" "physical_domain" {
  template_id = mso_template.fabric_policy_template.id
  name        = "physical_domain"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the Physical Domain.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Physical Domain.
* `id` - (Read-Only) The unique Terraform identifier of the Physical Domain.
* `description` - (Read-Only) The description of the Physical Domain.
* `vlan_pool_uuid` - (Read-Only) The NDO UUID of the associated VLAN Pool.
