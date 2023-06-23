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
  schema_id     = data.mso_schema.schema1.id
  bd_name       = "bd4"
  template_name = "Template1"
  site_id       = data.mso_site.site1.id
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy the Site bridge domain.
* `site_id` - (Required) SiteID under which you want to deploy the bridge domain.
* `template_name` - (Required) Template where Site Bd is to be created.
* `bd_name` - (Required) Name of the Site bridge domain. The name of the bridge domain should be present in the bridge domain list of the given `schema_id` and `template_name`.

## Attribute Reference ##

* `host_route` - (Read-Only) Value to check whether host-based routing is enabled.
* `svi_mac` - (Read-Only) Value of the SVI MAC Address.
