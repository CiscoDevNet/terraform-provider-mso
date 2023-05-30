---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_route_leak"
sidebar_current: "docs-mso-resource-schema_site_vrf_route_leak"
description: |-
  Manages MSO Schema Site Vrf Route Leak.
---

# mso_schema_site_vrf_route_leak #

Manages MSO Schema Site Vrf Route Leak.

## Example Usage ##

```hcl

resource "mso_schema_site_vrf_route_leak" "vrf1" {
  schema_id       = mso_schema.demo_schema.id
  template_name   = "Template1"
  site_id         = mso_site.demo_site.id
  vrf_name        = mso_schema_template_vrf.vrf1.name
  target_vrf_name = mso_schema_template_vrf.vrf2.name
  tenant_name     = mso_tenant.demo_tenant.name
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
* `tenant_name` - (Required) Name of the tenant.
* `type` - (Optional) Type of the Vrf Route Leak. Allowed values are `leak_all`, `subnet_ip` and `all_subnet_ips`. Default to `leak_all`.
* `subnet_ips` - (Optional) List of subnet ips which need to be leaked.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Vrf Route Leak can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_vrf_route_leak.vrf1RouteLeak {schema_id}/site/{site_id}/template/{template_name}/vrf/{vrf_name}/routeleak/{target_vrf_schema_id}/{target_vrf_template_name}/{target_vrf_name}
```
