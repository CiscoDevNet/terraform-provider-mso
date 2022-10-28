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

data "mso_schema" "schema1" {
  name  = "schema1"
}

```

## Argument Reference ##

* `name` - (Required) name of the schema.

## Attribute Reference ##

* `template_name` - (Optional) **Deprecated**. Name of template attached to the schema.
* `tenant_id` - (Optional) **Deprecated**. tenant_id for the schema.
* `template` - (Optional) A block that represents the template associated with the schema. Type - Block.
  * `name` - Name of template.
  * `display_name` - Display name for the template.
  * `tenant_id` - tenant_id for the template.

