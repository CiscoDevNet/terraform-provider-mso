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
  schema_id  = data.mso_schema.schema1.id
}

```

## Argument Reference ##

* `name` - (Required) name of the site to fetch.
* `schema_id` - (Required) The schema-id where site is associated.

## Attribute Reference ##

* `template_name` - (Optional) The name of the template deployed to the site.
* `site_id` - (Optional) Site id is set to the MSO site UUID.
