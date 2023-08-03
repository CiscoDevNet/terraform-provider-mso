---
layout: "mso"
page_title: "MSO: mso_schema_site_bd"
sidebar_current: "docs-mso-data-source-schema_site_bd"
description: |-
   Data source for MSO Schema Site Bridge Domain (BD).
---

# mso_schema_site_bd #

  Data source for MSO Schema Site Bridge Domain (BD).

## Example Usage ##

```hcl

data "mso_schema_site_bd" "example" {
  schema_id     = data.mso_schema.schema1.id
  bd_name       = "bd4"
  template_name = "Template1"
  site_id       = data.mso_site.site1.id
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the BD is deployed.
* `site_id` - (Required) The site ID under which the BD is deployed.
* `template_name` - (Required) The template name under which the BD is deployed.
* `bd_name` - (Required) The name of the BD.

## Attribute Reference ##

* `host_route` - (Read-Only) Whether host-based routing is enabled for the BD.
* `svi_mac` - (Read-Only) The SVI MAC Address of the BD.
