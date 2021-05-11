---
layout: "mso"
page_title: "MSO: mso_schema_template_bd_subnet"
sidebar_current: "docs-mso-resource-schema_template_bd_subnet"
description: |-
  Manages MSO Schema Template Bridge Domain Subnet.
---

# mso_schema_template_bd_subnet #

Manages MSO Schema Template Bridge Domain Subnet.

## Example Usage ##

```hcl
resource "mso_schema_template_bd_subnet" "bdsub1" {
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  bd_name = "testBD"
  ip = "10.23.13.0/8"
  scope = "public"
  description = "Description for the subnet"
  shared = true
  no_default_gateway = false
  querier = true
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Bridge Domain Subnet.
* `template_name` - (Required) Template where Bridge Domain Subnet to be created.
* `bd_name` - (Required) Name of Bridge Domain.
* `ip` - (Required) The IP range in CIDR notation.
* `scope` - (Required) The scope of the subnet. Allowed values are `private`, `public`.
* `description` - (Optional) The description of the subnet.
* `shared` - (Required) Whether this subnet is shared between VRFs.
* `no_default_gateway` - (Optional) Whether this subnet has a default gateway.
* `querier` - (Optional) Whether this subnet is an IGMP querier.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Bridge Domain Subnet can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_bd_subnet.bdsub1 {schema_id}/template/{template_name}/bd/{bd_name}/subnet/{ip}
```
