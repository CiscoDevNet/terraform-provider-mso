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

* `schema_id` - (Required) The schema ID under which the BD is deployed.
* `site_id` - (Required) The site ID under which the BD is deployed.
* `template_name` - (Required) The template name under which the BD is deployed.
* `bd_name` - (Required)  The name of the BD.
* `l3out_name` - (Required) The name of the L3out.
* `l3out_schema_id` - (Optional) The schema ID of the L3out. The `schema_id` of the BD will be used if not provided. 
* `l3out_template_name` - (Optional) The template name of the L3out. The `template_name` of the BD will be used if not provided. 
