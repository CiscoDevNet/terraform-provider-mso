---
layout: "mso"
page_title: "MSO: mso_schema_site_bd_l3out"
sidebar_current: "docs-mso-data-source-schema_site_bd_l3out"
description: |-
  Data source for MSO Schema Site Bridge Domain L3out.
---

# mso_schema_site_bd_l3out #

Data source for MSO Schema Site Bridge Domain L3out.

## Example Usage ##

```hcl

data "mso_schema_site_bd_l3out" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  bd_name       = "WebServer-Finance"
  l3out_name    = "ccc"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the L3out is deployed.
* `site_id` - (Required) The site ID under which the L3out is deployed.
* `template_name` - (Required) The template name under which the L3out is deployed.
* `bd_name` - (Required)  The bridge domain name under which the L3out is deployed.
* `l3out_name` - (Required) The name of the L3out.
