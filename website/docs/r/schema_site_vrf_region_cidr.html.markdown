---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_cidr"
sidebar_current: "docs-mso-resource-schema_site_vrf_region_cidr"
description: |-
  Manages MSO Schema Site Vrf Region Cidr.
---

# mso_schema_site_vrf_region_cidr #

Manages MSO Schema Site Vrf Region Cidr.

## Example Usage ##

```hcl
resource "mso_schema_site_vrf_region_cidr" "vrfRegionCidr" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id = "5ce2de773700006a008a2678"
  vrf_name = "Campus"
  region_name = "region1"
  ip = "2.2.2.2/2"
  primary = false
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.
* `ip` - (Required) The name of the region CIDR to manage.
* `template_name` - (Required) Template where Vrf Region to be created.
* `primary` - (Required) Whether this is the primary CIDR.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Vrf Region Cidr can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_vrf_region_cidr.vrfRegionCidr {schema_id}/site/{site_id}/vrf/{vrf_name}/region/{region_name}/cidrIP/{ip}
```

