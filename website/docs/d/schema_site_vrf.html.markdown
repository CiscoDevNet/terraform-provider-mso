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
  site_id = "5c7c95d9510000cf01c1ee3d"
  schema_id ="5c6c16d7270000c710f8094d"
  vrf_name = "vrf5810"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Vrf.
* `site_id` - (Required) SiteID under which you want to deploy Vrf.
* `vrf_name` - (Required) Name of Site Vrf. The name of the VRF should be present in the VRF list of the given `schema_id` and `template_name`

## Attribute Reference ##

* `template_name` - (Optional) Template where Site Vrf to be created.
