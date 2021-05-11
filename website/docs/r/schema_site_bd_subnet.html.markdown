---
layout: "mso"
page_title: "MSO: mso_schema_site_bd_subnet"
sidebar_current: "docs-mso-resource-schema_site_bd_subnet"
description: |-
  Manages MSO Schema Site Bridge Domain(BD) Subnet.
---

# mso_schema_site_bd_subnet #

Manages MSO Schema Site Bridge Domain(BD) Subnet.

## Example Usage ##

```hcl
resource "mso_schema_site_bd_subnet" "sub1" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
  bd_name = "WebServer-Finance"
  ip = "200.168.240.1/24"
  description = "Subnet 1"
  shared = false
  scope = "private"
  querier = false
  no_default_gateway = false
  }
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Subnet.
* `site_id` - (Required) SiteID under which you want to deploy Subnet.
* `bd_name` - (Required) Bd name under which you want to deploy Subnet. The Bd Name Reference should have `l2Stretch` set to `false` to be able to add a subnet.
* `ip` - (Required) The IP of the Subnet.
* `template_name` - (Required) Template name under which you want to deploy Subnet.
* `scope` - (Optional) The scope of the subnet. Allowed values are `private` and `public`.
* `shared` - (Optional) Whether this subnet is shared between VRFs.
* `querier` - (Optional) Whether this subnet is an IGMP querier.
* `no_default_gateway` - (Optional) Whether this subnet has a default gateway.
* `description` - (Optional) The description of this subnet. 

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Bridge Domain(BD) Subnet can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_bd_subnet.sub1 {schema_id}/site/{site_id}/bd/{bd_name}/subnet/{ip}
```