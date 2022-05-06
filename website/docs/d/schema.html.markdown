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

* `template` - A block that represents the template associated with the schema. Type - Block.
  * `name` - Name of template.
  * `display_name` - Display name for the template.
  * `tenant_id` - tenant_id for the template.
