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

data "mso_schema" "demo_schema" {
  name = "demo_schema"
}

```

## Argument Reference ##

* `name` - (Required) name of the schema.

## Attribute Reference ##

* `template_name` - (Read-Only) **Deprecated**. Name of template attached to the schema.
* `tenant_id` - (Read-Only) **Deprecated**. tenant_id for the schema.
* `description` - (Read-Only) The description of the schema.
* `template` - (Read-Only) A block that represents the template associated with the schema. Type - Block.
  * `name` - (Read-Only) The name of the template.
  * `display_name` - (Read-Only) The display name of the template.
  * `description` - (Read-Only) The description of the template.
  * `tenant_id` - (Read-Only) The tenant-id to associate with the template.
  * `template_type` - (Read-Only) The template type of the template.

