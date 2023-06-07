---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region"
sidebar_current: "docs-mso-data-source-schema_site_vrf_region"
description: |-
  Data source for MSO Schema Site Vrf Region.
---

# mso_schema_site_vrf_region #

Data source for MSO Schema Site Vrf Region.

## Example Usage ##

```hcl

data "mso_schema_site_vrf_region" "vrfRegion" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  vrf_name      = "Campus"
  region_name   = "westus"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region.
* `template_name` - (Required) Template under which you want to deploy Vrf Region.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.

## Attribute Reference ##

* `cidr` - (Read-Only) CIDR to set into region
* `cidr.cidr_ip` - (Read-Only) IP address for CIDR.
* `cidr.primary` - (Read-Only) Primary flag to set CIDR as primary.

* `cidr.subnet` - (Read-Only) Subnets to associate with CIDR.
* `cidr.subnet.ip` - (Read-Only) IP address for the subnet.
* `cidr.subnet.name` - (Read-Only) Name for the subnet.
* `cidr.subnet.zone` - (Read-Only) The name of the availability zone for the subnet.
* `cidr.subnet.usage` - (Read-Only) Usage information of particular subnet.
* `cidr.subnet.subnet_group` - (Read-Only) The name of the subnet group label for the subnet.

* `vpn_gateway` - (Read-Only) VPN gateway flag.
* `hub_network_enable` - (Read-Only) Hub Network enable flag.

* `hub_network` - (Read-Only) Hub Network to set into the region.
* `hub_network.name` - (Read-Only) Name of the hub network.
* `hub_network.tenant_name` - (Read-Only) Tenant name for the hub network.
