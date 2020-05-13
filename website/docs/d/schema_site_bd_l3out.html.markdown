---
layout: "mso"
page_title: "MSO: mso_schema_site_bd_l3out"
sidebar_current: "docs-mso-data-source-schema_site_bd_l3out"
description: |-
  Data source for MSO Schema Site Bridge Domain L3out.
---

# mso_schema_site_bd_l3out #

Data source for MSO Schema Site Bridge Domain L3out.

## Example Usage ##

```hcl
data "mso_schema_site_bd_l3out" "bdL3out" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  site_id = "5c7c95b25100008f01c1ee3c"
  bd_name = "WebServer-Finance"
  l3out_name = "ccc"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Bd L3out.
* `site_id` - (Required) SiteID under which you want to deploy Bd L3out.
* `bd_name` - (Required) Name of Bridge Domain.
* `l3out_name` - (Required) Name of L3out to manage.

## Attribute Reference ##

* `template_name` - (Optional) Template where Bd L3out to be created.
