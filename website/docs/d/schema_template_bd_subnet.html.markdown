---
layout: "mso"
page_title: "MSO: mso_schema_template_bd_subnet"
sidebar_current: "docs-mso-data-source-schema_template_bd_subnet"
description: |-
  Data source for MSO Schema Template Bridge Domain (BD) Subnet.
---

# mso_schema_template_bd_subnet #

Data source for MSO Schema Template Bridge Domain (BD) Subnet.

## Example Usage ##

```hcl

data "mso_schema_template_bd_subnet" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  bd_name       = "testBD"
  ip            = "10.23.13.0/8"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Subnet.
* `template_name` - (Required) The template name of the Subnet.
* `bd_name` - (Required) The name of the BD.
* `ip` - (Required) The IP range of the Subnet in CIDR notation.

## Attribute Reference ##

* `scope` - (Read-Only) The scope of the Subnet.
* `shared` - (Read-Only) Whether the Subnet is shared between VRFs.
* `description` - (Read-Only) The description of the Subnet.
* `no_default_gateway` - (Read-Only) Whether the Subnet has a default gateway.
* `querier` - (Read-Only) Whether the Subnet is an IGMP querier.
* `primary` - (Read-Only) Whether the Subnet is the primary Subnet.
