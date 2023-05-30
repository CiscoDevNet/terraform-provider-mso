---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_route_leak"
sidebar_current: "docs-mso-data-source-schema_site_vrf_route_leak"
description: |-
  Data source for MSO Schema Site Vrf Route Leak.
---

# mso_schema_site_vrf_route_leak #

Data source for MSO Schema Site Vrf Route Leak.

## Example Usage ##

```hcl

data "mso_schema_site_vrf_route_leak" "vrf1" {
  schema_id       = mso_schema.demo_schema.id
  template_name   = "Template1"
  site_id         = mso_site.demo_site.id
  vrf_name        = mso_schema_template_vrf.vrf1.name
  target_vrf_name = mso_schema_template_vrf.vrf2.name
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Vrf Route Leak.
* `template_name` - (Required) Template under which you want to deploy Vrf Route Leak.
* `site_id` - (Required) SiteID under which you want to deploy Vrf Route Leak.
* `vrf_name` - (Required) Vrf under which you want to deploy Vrf Route Leak.
* `target_vrf_schema_id` - (Optional)  SchemaID of the target vrf. The `schema_id` of the Vrf will be used if not provided. 
* `target_vrf_template_name` - (Optional) Template name of the target vrf. The `template_name` of the Vrf will be used if not provided. 
* `target_vrf_name` - (Required) Name of the target Vrf.

## Attribute Reference ##

* `tenant_name` - (Read-Only) Name of the tenant.
* `type` - (Read-Only) Type of the Vrf Route Leak. 
* `subnet_ips` - (Read-Only) List of subnet ips which need to be leaked.
