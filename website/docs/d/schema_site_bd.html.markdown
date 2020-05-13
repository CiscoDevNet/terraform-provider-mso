---
layout: "mso"
page_title: "MSO: mso_schema_site_bd"
sidebar_current: "docs-mso-data-source-schema_site_bd"
description: |-
  MSO Schema Site Bridge Domain(BD) Data source.
---

# mso_schema_site_bd #

 MSO Schema Site Bridge Domain(bd) Data source.

## Example Usage ##

```hcl
data "mso_schema_site_bd" "st10" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  bd_name = "bd4"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Bd.
* `site_id` - (Required) SiteID under which you want to deploy Bd.
* `bd_name` - (Required) Name of Site Bd. The name of the Bd should be present in the Bd list of the given `schema_id` and `template_name`

## Attribute Reference ##

* `template_name` - (Optional) Template where Site Bd to be created.
* `host` - (Optional) Value to check whether host-based routing is enabled. Default value is `false`.
