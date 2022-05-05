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

* `template_name` - (Optional) name of template attached to the schema.
* `tenant_id` - (Optional) tenant_id for the schema.
