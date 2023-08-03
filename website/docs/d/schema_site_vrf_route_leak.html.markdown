---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_route_leak"
sidebar_current: "docs-mso-data-source-schema_site_vrf_route_leak"
description: |-
  Data source for MSO Schema Site VRF Route Leak.
---

# mso_schema_site_vrf_route_leak #

Data source for MSO Schema Site VRF Route Leak.

## Example Usage ##

```hcl

data "mso_schema_site_vrf_route_leak" "example" {
  schema_id       = mso_schema.demo_schema.id
  template_name   = "Template1"
  site_id         = mso_schema_site.demo_site.site_id
  vrf_name        = mso_schema_template_vrf.vrf1.name
  target_vrf_name = mso_schema_template_vrf.vrf2.name
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Route Leak is deployed.
* `site_id` - (Required) The site ID under which the Route Leak is deployed.
* `template_name` - (Required) The template name under which the Route Leak is deployed.
* `vrf_name` - (Required) The name of the VRF under which the Route Leak is deployed.
* `target_vrf_schema_id` - (Optional)  The schema ID of the target vrf. The `schema_id` of the VRF will be used if not provided. 
* `target_vrf_template_name` - (Optional) The template name of the target vrf. The `template_name` of the VRF will be used if not provided. 
* `target_vrf_name` - (Required) The name of the target VRF.

## Attribute Reference ##

* `tenant_name` - (Read-Only) The name of the tenant.
* `type` - (Read-Only) The type of the Route Leak. 
* `subnet_ips` - (Read-Only) The list of subnet ips which need to be leaked.
