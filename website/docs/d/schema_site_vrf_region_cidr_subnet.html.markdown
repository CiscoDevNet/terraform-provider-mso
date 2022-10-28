---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_cidr_subnet"
sidebar_current: "docs-mso-data-source-schema_site_vrf_region_cidr_subnet"
description: |-
  Data source for MSO Schema Site Vrf Region Cidr Subnet.
---

# mso_schema_site_vrf_region_cidr_subnet #

 Data source for MSO Schema Site Vrf Region Cidr Subnet.

## Example Usage ##

```hcl

 data "mso_schema_site_vrf_region_cidr_subnet" "vrfRegion" {
  schema_id   = data.mso_schema.schema1.id
  site_id     = data.mso_site.site1.id
  vrf_name    = "Campus"
  region_name = "westus"
  cidr_ip     = "1.1.1.1/24"
  ip          = "207.168.240.1/24"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region Cidr Subnet.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region Cidr Subnet.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.
* `cidr_ip` - (Required) The IP range of for the region CIDR where Vrf Region Cidr Subnet to be created.
* `ip` - (Required) The IP subnet of this region CIDR.


## Attribute Reference ##

* `template_name` - (Optional) Template where Vrf Region Cidr Subnet to be created.
* `zone` - (Optional) The name of the zone for the region CIDR subnet.
* `usage` - (Optional) The usage for the region CIDR Subnet.

## Note ##
Multiple Subnets with same Ip are allowed, but the operations will take place on the first found Subnet with the given Ip.
