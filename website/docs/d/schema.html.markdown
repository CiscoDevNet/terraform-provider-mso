---
layout: "mso"
page_title: "MSO: mso_schema"
sidebar_current: "docs-mso-data-source-schema"
description: |-
  Data source for MSO Schema
---

# mso_schema #

Data source for MSO schema  

## Example Usage ##

```hcl
data "mso_schema" "sample_schema" {
  name  = "schema1"
}
```

## Argument Reference ##

* `name` - (Required) name of the schema.

## Attribute Reference ##

* `template_name` - (Optional) name of templates for this schema.
* `tenant_id` - (Optional) temant_id for this schema.
