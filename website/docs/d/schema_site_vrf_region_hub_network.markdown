---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_hub_network"
sidebar_current: "docs-mso-data-source-schema_site_vrf_region_hub_network"
description: |-
  Data source for MSO Schema Site VRF Region Hub Network
---

# mso_schema_site_vrf_region_hub_network #

Data source for MSO schema site VRF Region Hub Network, to fetch the MSO schema site VRF Region Hub Network details.

## Example Usage ##

```hcl
data "mso_schema_site_vrf_region_hub_network" "example"{
    schema_id = mso_schema_site.example.schema_id
    template_name = mso_schema_site.example.template_name
    site_id = mso_schema_site.example.site_id
    vrf_name = mso_schema_site_vrf.example.vrf_name
    region_name = mso_schema_site_vrf_region.example.region_name
    name = mso_schema_site_vrf_region_hub_network.example.name
    tenant_name = data.mso_tenant.example.id
}
```

## Argument Reference ##
* `schema_id` - (Required) The schema-id where Vrf Region Hub Network is added.
* `name` - (Required) Name of the added Vrf Region Hub Network.
* `template_name` - (Required) Template name associated with the Vrf Region Hub Network.
* `vrf_name` - (Required) VRF name associated with the Vrf Region Hub Network.
* `site_id` - (Required) SiteID associated with the Vrf Region Hub Network.
* `tenant_name` - (Required) Tenant Name associated with Vrf Region Hub Network.
* `region_name` - (Required) Region Name associated with Vrf Region Hub Network.

## Attribute Reference ##

No attributes are exported.