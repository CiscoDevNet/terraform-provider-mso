---
layout: "mso"
page_title: "MSO: mso_fabric_policies_l3_domain"
sidebar_current: "docs-mso-data-source-fabric_policies_l3_domain"
description: |-
  Data source for L3 Domain.
---

# mso_fabric_policies_l3_domain #

Data source for Layer 3 (L3) Domain. This data source is only supported NDO v4.3 and later.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> L3 Domains

## Example Usage ##

```hcl
data "mso_fabric_policies_l3_domain" "l3_domain" {
  template_id = mso_template.fabric_template.id
  name        = "test_l3_domain"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the L3 Domain to retrieve.

## Attribute Reference ##

* `description` - (Read-Only) The description of the L3 Domain.
* `vlan_pool_uuid` - (Read-Only) The NDO UUID of the VLAN Pool to associate with this L3 Domain.
* `uuid` - (Read-Only) The NDO UUID of the L3 Domain.
* `id` - (Read-Only) The unique Terraform identifier of the L3 Domain in the template.
