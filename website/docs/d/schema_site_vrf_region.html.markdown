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
  schema_id   = data.mso_schema.schema1.id
  site_id     = data.mso_site.site1.id
  vrf_name    = "Campus"
  region_name = "westus"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region.
* `template_name` - (Required) Template where Vrf Region to be created.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.

## Attribute Reference ##

* `cidr` - (Read-Only) CIDR to set into region
* `cidr.cidr_ip` - (Read-Only) IP address for CIDR.
* `cidr.primary` - (Read-Only) primary flag to set CIDR as primary. Only one CIDR can be set as primary.

* `cidr.subnet` - (Read-Only) subnets to associate with CIDR.
* `cidr.subnet.ip` - (Read-Only) IP address for the subnet.
* `cidr.subnet.name` - (Read-Only) Name for the subnet.
* `cidr.subnet.zone` - (Read-Only) zone for the subnet.
* `cidr.subnet.usage` - (Read-Only) usage information of particular subnet.
* `cidr.subnet.subnet_group` - (Read-Only) The name of the subnet group label for the subnet.

* `vpn_gateway` - (Read-Only) VPN gateway flag.
* `hub_network_enable` - (Read-Only) Hub Network enable flag. To set hub network in region, this attribute should be true. this parameter is supported in MSO v3.0 or higher with Cloud APIC version 5.0 or higher.

* `hub_network` - (Read-Only) Hub Network to set into the region. this parameter is supported in MSO v3.0 or higher with Cloud APIC version 5.0 or higher.
* `hub_network.name` - (Read-Only) name of the hub network.
* `hub_network.tenant_name` - (Read-Only) Tenant name for the hub network.
