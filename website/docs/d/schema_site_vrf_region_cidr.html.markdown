---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_cidr"
sidebar_current: "docs-mso-data-source-schema_site_vrf_region_cidr"
description: |-
  Data source for MSO Schema Site Vrf Region Cidr.
---

# mso_schema_site_vrf_region_cidr #

Data source for MSO Schema Site Vrf Region Cidr.

## Example Usage ##

```hcl
data "mso_schema_site_vrf_region_cidr" "vrfRegionCidr" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  site_id = "5ce2de773700006a008a2678"
  vrf_name = "Campus"
  region_name = "westus"
  ip = "192.168.241.0/24"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.
* `ip` - (Required) The name of the region CIDR to manage.

## Attribute Reference ##

* `template_name` - (Optional) Template where Vrf Region to be created.
* `primary` - (Optional) Whether this is the primary CIDR.
