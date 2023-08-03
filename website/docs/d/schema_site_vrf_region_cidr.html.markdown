---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf_region_cidr"
sidebar_current: "docs-mso-data-source-schema_site_vrf_region_cidr"
description: |-
  Data source for MSO Schema Site VRF Region CIDR.
---

# mso_schema_site_vrf_region_cidr #

Data source for MSO Schema Site VRF Region CIDR.

## Example Usage ##

```hcl

data "mso_schema_site_vrf_region_cidr" "example" {
  schema_id   = data.mso_schema.schema1.id
  site_id     = data.mso_site.site1.id
  template_name = "Template1"
  vrf_name    = "Campus"
  region_name = "westus"
  ip          = "192.168.241.0/24"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the CIDR is deployed.
* `site_id` - (Required) The site ID under which the CIDR is deployed.
* `template_name` - (Required) The template name under which the CIDR is deployed.
* `vrf_name` - (Required) The name of the VRF under which the CIDR is deployed.
* `region_name` - (Required) The name of the VRF Region under which the CIDR is deployed.
* `ip` - (Required) The IP range of the VRF Region in CIDR notation.

## Attribute Reference ##

* `primary` - (Read-Only) Whether this is the primary CIDR.
