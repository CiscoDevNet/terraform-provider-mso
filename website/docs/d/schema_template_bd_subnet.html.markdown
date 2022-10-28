---
layout: "mso"
page_title: "MSO: mso_schema_template_bd_subnet"
sidebar_current: "docs-mso-data-source-schema_template_bd_subnet"
description: |-
  Data source for MSO Schema Template Bridge Domain Subnet.
---

# mso_schema_template_bd_subnet #

Data source for MSO Schema Template Bridge Domain Subnet.

## Example Usage ##

```hcl

data "mso_schema_template_bd_subnet" "bridge_domain_subnet" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  bd_name       = "testBD"
  ip            = "10.23.13.0/8"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Bridge Domain.
* `template_name` - (Required) Template where Bridge Domain to be created.
* `bd_name` - (Required) Name of Bridge Domain.
* `ip` - (Required) The IP range in CIDR notation.

## Attribute Reference ##

* `scope` - (Optional) The scope of the subnet.
* `shared` - (Optional) Whether this subnet is shared between VRFs.
* `description` - (Optional) The description for the subnet.
* `no_default_gateway` - (Optional) Whether this subnet has a default gateway.
* `querier` - (Optional) Whether this subnet is an IGMP querier.
