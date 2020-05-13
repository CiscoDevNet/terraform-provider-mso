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
  site_id = "5c7c95d9510000cf01c1ee3d"
  schema_id ="5c6c16d7270000c710f8094d"
  vrf_name = "vrf5810"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Vrf.
* `site_id` - (Required) SiteID under which you want to deploy Vrf.
* `template_name` - (Required) Template where Site Vrf to be created.
* `vrf_name` - (Required) Name of Site Vrf. The name of the VRF should be present in the VRF list of the given `schema_id` and `template_name`

## Attribute Reference ##

No attributes are exported




