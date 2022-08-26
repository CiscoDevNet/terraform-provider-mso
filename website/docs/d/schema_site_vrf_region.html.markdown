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
  schema_id     = "5d5dbf3f2e0000580553ccce"
  site_id       = "5ce2de773700006a008a2678"
  vrf_name      = "Campus"
  region_name   = "westus"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.

## Attribute Reference ##

* `template_name` - (Optional) Template where Vrf Region to be created.

* `cidr` - (Optional) CIDR to set into region
* `cidr.cidr_ip` - (Optional) IP address for CIDR.
* `cidr.primary` - (Optional) primary flag to set CIDR as primary. Only one CIDR can be set as primary.

* `cidr.subnet` - (Optional) subnets to associate with CIDR.
* `cidr.subnet.ip` - (Optional) IP address for the subnet.
* `cidr.subnet.name` - (Optional) Name for the subnet.
* `cidr.subnet.zone` - (Optional) zone for the subnet.
* `cidr.subnet.usage` - (Optional) usage information of particular subnet.

* `vpn_gateway` - (Optional) VPN gateway flag.
* `hub_network_enable` - (Optional) Hub Network enable flag. To set hub network in region, this attribute should be true. this parameter is supported in MSO v3.0 or higher with Cloud APIC version 5.0 or higher.

* `hub_network` - (Optional) Hub Network to set into the region. this parameter is supported in MSO v3.0 or higher with Cloud APIC version 5.0 or higher.
* `hub_network.name` - (Optional) name of the hub network.
* `hub_network.tenant_name` - (Optional) Tenant name for the hub network.
