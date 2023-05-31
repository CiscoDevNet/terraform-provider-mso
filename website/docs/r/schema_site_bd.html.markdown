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
  svi_mac       = "00:22:BD:F8:19:FF"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy the Site bridge domain.
* `template_name` - (Required) Template where the Site bridge domain is to be created.
* `site_id` - (Required) SiteID under which you want to deploy the bridge domain.
* `bd_name` - (Required) Name of the Site bridge domain. The name of the bridge domain should be present in the bridge domain list of the given `schema_id` and `template_name`
* `host_route` - (Optional) Value to check whether the host-based routing is enabled. Default value is `false`.
* `svi_mac` - (Optional) Value of the SVI MAC Address.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Bridge Domain(BD) can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_bd.bd1 {schema_id}/site/{site_id}-{template_name}/bd/{bd_name}
```