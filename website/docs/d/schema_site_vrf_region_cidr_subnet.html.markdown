---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_cidr_subnet"
sidebar_current: "docs-mso-data-source-schema_site_vrf_region_cidr_subnet"
description: |-
  Data source for MSO Schema Site VRF Region CIDR Subnet.
---

# mso_schema_site_vrf_region_cidr_subnet #

 Data source for MSO Schema Site VRF Region CIDR Subnet.

## Example Usage ##

```hcl

 data "mso_schema_site_vrf_region_cidr_subnet" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  vrf_name      = "Campus"
  region_name   = "westus"
  cidr_ip       = "1.1.1.1/24"
  ip            = "207.168.240.1/24"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Subnet is deployed.
* `site_id` - (Required) The site ID under which the Subnet is deployed.
* `template_name` - (Required) The template name under which the Subnet is deployed.
* `vrf_name` - (Required) The name of the VRF under which the Subnet is deployed.
* `region_name` - (Required) The name of the VRF Region under which the Subnet is deployed.
* `cidr_ip` - (Required) The IP range of the VRF Region where the Subnet is deployed in CIDR notation..
* `ip` - (Required) The IP of the Subnet.

## Attribute Reference ##

* `zone` - (Read-Only) The availability zone name of the Subnet. 
* `name` - (Read-Only) The name Subnet of the Subnet.
* `usage` - (Read-Only) The usage of the Subnet.
* `subnet_group` - (Read-Only) The group of the Subnet.

## Note ##
Multiple Subnets with same Ip are allowed, but the operations will take place on the first found Subnet with the given Ip.
