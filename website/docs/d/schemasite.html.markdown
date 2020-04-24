---
layout: "mso"
page_title: "MSO: mso_schema_site"
sidebar_current: "docs-mso-data-source-schema_site"
description: |-
  Data source for MSO Schema Site
---

# mso_schema_site #

Data source for MSO schema site, to fetch the MSO schema site details.

## Example Usage ##

```hcl
data "mso_schema_site" "sample_schema_site" {
  name       = "sitename"
  schema_id  = "schema-id"
}
```

## Argument Reference ##

* `name` - (Required) name of the schema.
* `schema_id` - (Required) The name of the template.


## Attribute Reference ##

* `template_name` - (Required) The name of the template deployed to the site.
* `site_id` - (Optional) Site id is set to the MSO site UUID.