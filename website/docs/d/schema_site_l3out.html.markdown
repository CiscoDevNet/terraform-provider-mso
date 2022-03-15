---
layout: "mso"
page_title: "MSO: mso_schema_site_l3out"
sidebar_current: "docs-mso-data-source-schema_site_l3out"
description: |-
  Data source for MSO Schema Site L3out
---

# mso_schema_site_l3out #

Data source for MSO schema site L3out, to fetch the MSO schema site L3out details.

## Example Usage ##

```hcl
data "mso_schema_site_l3out" "exmple" {
    vrf_name = mso_schema_site_vrf.example.vrf_name
    l3out_name = mso_schema_site_l3out.example.l3out_name
    template_name = mso_site.example.template_name
    site_id = mso_site.example.site_id
    schema_id = mso_site.example.schema_id
}
```

## Argument Reference ##
* `schema_id` - (Required) The schema-id where L3out is added.
* `l3out_name` - (Required) Name of the added L3out.
* `template_name` - (Required) Template name associated with the L3out.
* `vrf_name` - (Required) VRF name associated with the L3out.
* `site_id` - (Required) SiteID associated with the L3out.

## Attribute Reference ##

No attributes are exported.