---
layout: "mso"
page_title: "MSO: mso_schema_site_bd"
sidebar_current: "docs-mso-resource-schema_site_bd"
description: |-
 Manages MSO Schema Site Bridge Domain(BD).
---

# mso_schema_site_bd #

 Manages MSO Schema Site Bridge Domain(bd)

## Example Usage ##

```hcl

resource "mso_schema_site_bd" "bd1" {
  schema_id     = mso_schema.schema1.id
  bd_name       = mso_schema_template_bd.bridge_domain.name
  template_name = "Template1"
  site_id       = mso_schema_site.schema_site.site_id
  host_route    = false
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Bd.
* `template_name` - (Required) Template where Site Bd to be created.
* `site_id` - (Required) SiteID under which you want to deploy Bd.
* `bd_name` - (Required) Name of Site Bd. The name of the Bd should be present in the Bd list of the given `schema_id` and `template_name`
* `host_route` - (Optional) Value to check whether host-based routing is enabled. Default value is `false`.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Bridge Domain(BD) can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_bd.bd1 {schema_id}/site/{site_id}/bd/{bd_name}
```