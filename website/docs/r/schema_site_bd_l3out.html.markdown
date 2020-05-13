---
layout: "mso"
page_title: "MSO: mso_schema_site_bd_l3out"
sidebar_current: "docs-mso-resource-schema_site_bd_l3out"
description: |-
  Manages MSO Schema Site Bridge Domain L3out.
---

# mso_schema_site_bd_l3out #

Manages MSO Schema Site Bridge Domain L3out.

## Example Usage ##

```hcl
resource "mso_schema_site_bd_l3out" "bdL3out" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
  bd_name = "WebServer-Finance"
  l3out_name = "zzz"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Bd L3out.
* `site_id` - (Required) SiteID under which you want to deploy Bd L3out.
* `bd_name` - (Required) Name of Bridge Domain.
* `l3out_name` - (Required) Name of L3out to manage.
* `template_name` - (Required) Template where Bd L3out to be created.

## Attribute Reference ##

No attributes are exported.
