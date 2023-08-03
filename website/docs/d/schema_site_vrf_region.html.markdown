---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region"
sidebar_current: "docs-mso-data-source-schema_site_vrf_region"
description: |-
  Data source for MSO Schema Site VRF Region.
---

# mso_schema_site_vrf_region #

Data source for MSO Schema Site VRF Region.

## Example Usage ##

```hcl

data "mso_schema_site_vrf_region" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  vrf_name      = "Campus"
  region_name   = "westus"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Region is deployed.
* `site_id` - (Required) The site ID under which the Region is deployed.
* `template_name` - (Required) The template name under which the Region is deployed.
* `vrf_name` - (Required) The name of the VRF under which the Region is deployed.
* `region_name` - (Required) The name of the Region.

## Attribute Reference ##

* `vpn_gateway` - (Read-Only) The VPN gateway flag of the Region.
* `hub_network_enable` - (Read-Only) The Hub Network enable flag of the Region.
* `cidr` - (Read-Only) A list of CIDRs for the Region.
    * `cidr_ip` - (Read-Only) The IP range of the Region.
    * `primary` - (Read-Only) Whether this is the primary CIDR.
    * `subnet` - (Read-Only) A list of Subnets for the CIDR.
        * `ip` - (Read-Only) The P address of the subnet.
        * `name` - (Read-Only) The name of the subnet.
        * `zone` - (Read-Only) The availability zone name of the Subnet. 
        * `usage` - (Read-Only) The usage of the Subnet.
        * `subnet_group` - (Read-Only) The group of the Subnet.
* `hub_network` - (Read-Only) A list of Hub Networks for the Region.
    * `name` - (Read-Only) The name of the hub network.
    * `tenant_name` - (Read-Only) The tenant name of the hub network.
