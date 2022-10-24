---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf"
sidebar_current: "docs-mso-resource-schema_site_vrf"
description: |-
  Manages MSO Schema Site VRF.
---

# mso_schema_site_vrf #

 Manages MSO Schema Site VRF.

## Example Usage ##

```hcl

resource "mso_schema_site_vrf" "vrf1" {
  template_name = "Template1"
  site_id       = mso_schema_site.schema_site.site_id
  schema_id     = mso_schema.schema1.id
  vrf_name      = mso_shema_template_vrf.vrf1.name
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Vrf.
* `site_id` - (Required) SiteID under which you want to deploy Vrf.
* `template_name` - (Required) Template where Site Vrf to be created.
* `vrf_name` - (Required) Name of Site Vrf. The name of the VRF should be present in the VRF list of the given `schema_id` and `template_name`

## Attribute Reference ##

No attributes are exported

## Importing ##

An existing MSO Schema Site Vrf can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_vrf.vrf1 {schema_id}/site/{site_id}/vrf/{vrf_name}
```