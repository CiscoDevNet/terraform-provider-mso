---
layout: "mso"
page_title: "MSO: mso_schema_site_bd_subnet"
sidebar_current: "docs-mso-data-source-schema_site_bd_subnet"
description: |-
  Data source for MSO Schema Site Bd Subnet.
---

# mso_schema_site_bd_subnet #

Data source for MSO Schema Site Bridge Domain(Bd) Subnet.

## Example Usage ##

```hcl
data "mso_schema_site_bd_subnet" "s1" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  site_id = "5c7c95b25100008f01c1ee3c"
  bd_name = "WebServer-Finance"
  ip = "200.168.240.1/24"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Subnet.
* `site_id` - (Required) SiteID under which you want to deploy Subnet.
* `bd_name` - (Required) Bd name under which you want to deploy Subnet.
* `ip` - (Required) The IP of the Subnet.

## Attribute Reference ##

* `template_name` - (Optional) Template name under which you want to deploy Subnet.
* `scope` - (Optional) The scope of the subnet.
* `shared` - (Optional) Whether this subnet is shared between VRFs.
* `querier` - (Optional) Whether this subnet is an IGMP querier.
* `no_default_gateway` - (Optional) Whether this subnet has a default gateway.
* `description` - (Optional) The description of this subnet. 

