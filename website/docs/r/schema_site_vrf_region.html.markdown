---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region"
sidebar_current: "docs-mso-resource-schema_site_vrf_region"
description: |-
  Manages MSO Schema Site Vrf Region.
---

# mso_schema_site_vrf_region #

Manages MSO Schema Site Vrf Region.

## Example Usage ##

```hcl

resource "mso_schema_site_vrf_region" "vrfRegion" {
  schema_id           = mso_schema.schema1.id
  template_name       = "Template1"
  site_id             = mso_schema_site.schema_site.site_id
  vrf_name            = mso_schema_site_vrf.vrf1.vrf_name
  region_name         = "us-east-1"
  vpn_gateway         = true
  hub_network_enable  = true
  hub_network = {
    name        = "hub-fualt"
    tenant_name = "infra"
  }
  cidr {
    cidr_ip = "2.2.2.2/10"
    primary = true
    subnet {
      ip    = "1.20.30.4"
      name  = "subnet1"
      zone  = "us-east-1b"
      usage = "sdfkhsdkf"
    }
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Region.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Region.
* `vrf_name` - (Required) Name of Vrf.
* `region_name` - (Required) Name of Region to manage.
* `template_name` - (Required) Template where Vrf Region to be created.

* `cidr` - (Required) CIDR to set into region
* `cidr.cidr_ip` - (Required) IP address for CIDR.
* `cidr.primary` - (Required) Primary flag to set CIDR as primary. Only one CIDR can be set as primary.

* `cidr.subnet` - (Required) Subnets to associate with CIDR.
* `cidr.subnet.ip` - (Required) IP address for the subnet.
* `cidr.subnet.name` - (Required) Name for the subnet.
* `cidr.subnet.zone` - (Optional) The name of the availability zone for the subnet. This argument is required for AWS sites.
* `cidr.subnet.usage` - (Optional) Usage information of particular subnet.

* `vpn_gateway` - (Optional) VPN gateway flag.
* `hub_network_enable` - (Optional) Hub Network enable flag. To set hub network in region, this attribute should be true. this parameter is supported in MSO v3.0 or higher with Cloud APIC version 5.0 or higher.

* `hub_network` - (Optional) Hub Network to set into the region. this parameter is supported in MSO v3.0 or higher with Cloud APIC version 5.0 or higher.
* `hub_network.name` - (Required) The name of the hub network.
* `hub_network.tenant_name` - (Required) Tenant name for the hub network.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Vrf Region can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_vrf_region.vrfRegion {schema_id}/site/{site_id}/vrf/{vrf_name}/region/{region_name}
```
