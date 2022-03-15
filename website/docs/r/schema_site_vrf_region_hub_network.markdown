---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_hub_network"
sidebar_current: "docs-mso-resource-schema_site_vrf_region_hub_network"
description: |-
  Manages MSO Schema Site VRF Region Hub Network
---

# mso_schema_site_vrf_region_hub_network #

Manages MSO Schema Site VRF Region Hub Network.

## Example Usage ##

```hcl
resource "mso_schema_site_vrf_region_hub_network" "example"{
    schema_id = mso_schema_site.example.schema_id
    template_name = mso_schema_site.example.template_name
    site_id = mso_schema_site.example.site_id
    vrf_name = mso_schema_site_vrf.example.vrf_name
    region_name = mso_schema_site_vrf_region.example.region_name
    name = "example"
    tenant_name = data.mso_tenant.example.id
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema-id where user wants to add Vrf Region Hub Network.
* `name` - (Required) Name of the Hub Network that user wants to add.
* `template_name` - (Required) Template name associated with the Vrf Region Hub Network.
* `vrf_name` - (Required) VRF name associated with the Vrf Region Hub Network.
* `site_id` - (Required) SiteID associated with the Vrf Region Hub Network.
* `tenant_name` - (Required) Tenant Name associated with Vrf Region Hub Network.
* `region_name` - (Required) Region Name associated with Vrf Region Hub Network.

## Attribute Reference ##
The only Attribute exposed for this resource is `id`. Which is set to the node name of Service Node created.

## Importing ##

An existing MSO Schema Site VRF Region Hub Network can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_vrf_region_hub_network.example {schema_id}/site/{site_id}/template/{template_name}/vrf/{vrf_name}/region/{region_name}/tenant/{tenant_name}/{name}
```