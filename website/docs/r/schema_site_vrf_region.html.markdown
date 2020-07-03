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
  schema_id           = "5efd6ea60f00005b0ebbd643"
  template_name       = "Template1"
  site_id             = "5efeb3c4190000cc12d05376"
  vrf_name            = "Myvrf"
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
* `cidr.cidr_ip` - (Required) Ip address for cidr.
* `cidr.primary` - (Required) primary flag to set above ip as primary for cidr. Only one ip must be set as primary.

* `cidr.subnet` - (Required) subnets to associate with cidr.
* `cidr.subnet.ip` - (Required) ip address for subnet.
* `cidr.subnet.zone` - (Required) zone for the subnet.
* `cidr.subnet.usage` - (Optional) usage information of particular subnet.

* `vpn_gateway` - (Optional) VPN gateway flag.
* `hub_network_enable` - (Optional) Hub Network enable flag. To set hub network in region, this attribute should be true.

* `hub_network` - (Optional) Hub Network to set into the region.  this parameter is supported in MSO v3.0 or higher with Cloud APIC version 5.0 or higher.
* `hub_network.name` - (Required) name of the hub network.
* `hub_network.tenant_name` - (Required) Tenant name for the hub network.

## Attribute Reference ##

No attributes are exported.
