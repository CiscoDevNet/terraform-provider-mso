---
layout: "mso"
page_title: "MSO: mso_schema_validate"
sidebar_current: "docs-mso-data-source-schema_validate"
description: |-
  Data source for MSO Schema Template Vrf
---

# mso_schema_validate #

Data source for MSO schema validate, to fetch the MSO schema validate details.

## Example Usage ##

```hcl
data "mso_schema_validate" "example" {
  schema_id = mso_schema.example.id
}
```

## Argument Reference ##

* `schema_id` - (Required) The schema-id which user want to validate.

## Attribute Reference ##

* `result` - (Optional) The validation result for schema_id given.
