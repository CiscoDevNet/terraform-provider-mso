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
  schema_id = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id = "5ce2de773700006a008a2678"
  vrf_name = "Campus"
  region_name = "westus"
  cidr_ip = "1.1.1.1/24"
  ip = "203.168.240.1/24"
  zone = "West"
  usage = "gateway"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region Cidr Subnet.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region Cidr Subnet.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.
* `template_name` - (Required) Template where Vrf Region Cidr Subnet to be created.
* `cidr_ip` - (Required) The IP range of for the region CIDR where Vrf Region Cidr Subnet to be created.
* `ip` - (Required) The IP subnet of this region CIDR.
* `zone` - (Required) The name of the zone for the region CIDR subnet.
* `usage` - (Optional) The usage for the region CIDR Subnet.

## Attribute Reference ##

No attributes are exported.

## Note ##
Multiple Subnets with same Ip are allowed, but the operations will take place on the first found Subnet with the given Ip.

## Importing ##

An existing MSO Schema Site Vrf Region Cidr Subnet can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_vrf_region_cidr_subnet.sub1 {schema_id}/site/{site_id}/vrf/{vrf_name}/region/{region_name}/cidrIP/{cidr_ip}/subnet/{ip}
```
