---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf"
sidebar_current: "docs-mso-data-source-schema_site_vrf"
description: |-
  MSO Schema Site VRF Data source.
---

# mso_schema_site_vrf #

 MSO Schema Site VRF Data source.

## Example Usage ##

```hcl

data "mso_schema_site_vrf" "v1" {
  site_id   = data.mso_site.site1.id
  schema_id = data.mso_schema.schema1.id
  vrf_name  = "vrf5810"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Vrf.
* `site_id` - (Required) SiteID under which you want to deploy Vrf.
* `vrf_name` - (Required) Name of Site Vrf. The name of the VRF should be present in the VRF list of the given `schema_id` and `template_name`

## Attribute Reference ##

* `template_name` - (Optional) Template where Site Vrf to be created.
