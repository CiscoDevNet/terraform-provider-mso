---
layout: "mso"
page_title: "MSO: mso_schema_site"
sidebar_current: "docs-mso-data-source-schema_site"
description: |-
  Data source for MSO Schema Site
---

# mso_schema_site #

Data source for MSO schema site

## Example Usage ##

```hcl
data "mso_schema_site" "sample_schema_site" {
  schema  = "schema1"
  site  = "bdsol-pod51"
  template  = "template1"
}
```

## Argument Reference ##

* `schema` - (Required) name of the schema.
* `template` - (Required) The name of the template.
* `site` - (Required) The name of the site to manage.

## Attribute Reference ##

No attributes are exported
