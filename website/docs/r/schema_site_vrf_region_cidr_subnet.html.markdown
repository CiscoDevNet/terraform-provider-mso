---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_cidr_subnet"
sidebar_current: "docs-mso-resource-schema_site_vrf_region_cidr_subnet"
description: |-
  Manages MSO Schema Site Vrf Region Cidr Subnet.
---

# mso_schema_site_vrf_region_cidr_subnet #

Manages MSO Schema Site Vrf Region Cidr Subnet.

## Example Usage ##

```hcl

resource "mso_schema_site_vrf_region_cidr_subnet" "sub1" {
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
  site_id       = mso_schema_site.schema_site.site_id
  vrf_name      = mso_schema_site_vrf_region_cidr.vrfRegionCidr.vrf_name
  region_name   = mso_schema_site_vrf_region_cidr.vrfRegionCidr.region_name
  cidr_ip       = mso_schema_site_vrf_region_cidr.vrfRegionCidr.ip
  ip            = "203.168.240.1/24"
  zone          = "West"
  usage         = "gateway"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region Cidr Subnet.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region Cidr Subnet.
* `template_name` - (Required)  Template under which you want to deploy Vrf Region.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.
* `cidr_ip` - (Required) The IP range of for the region CIDR where Vrf Region Cidr Subnet to be created.
* `ip` - (Required) The IP subnet of this region CIDR.
* `zone` - (Optional) The name of the availability zone for the region CIDR subnet. This argument is required for AWS sites.
* `name` - (Optional) The name for the region CIDR Subnet.
* `usage` - (Optional) The usage for the region CIDR Subnet.
* `subnet_group` - (Optional) The subnet group for the region CIDR Subnet.

## Attribute Reference ##

No attributes are exported.

## Note ##
Multiple Subnets with same Ip are allowed, but the operations will take place on the first found Subnet with the given Ip.

## Importing ##

An existing MSO Schema Site Vrf Region Cidr Subnet can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_vrf_region_cidr_subnet.sub1 {schema_id}/site/{site_id}/template/{template_name}/vrf/{vrf_name}/region/{region_name}/cidrIP/{cidr_ip}/subnet/{ip}
```
