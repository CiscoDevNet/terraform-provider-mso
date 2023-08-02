---
layout: "mso"
page_title: "MSO: mso_schema"
sidebar_current: "docs-mso-data-source-schema"
description: |-
  Data source for MSO Schema.
---

# mso_schema #

Data source for MSO schema.

## Example Usage ##

```hcl

data "mso_schema" "example" {
  name = "demo_schema"
}

```

## Argument Reference ##

* `name` - (Required) The name of the Schema.

## Attribute Reference ##

* `template_name` - (Read-Only) **Deprecated**. The template name of the Schema.
* `tenant_id` - (Read-Only) **Deprecated**. The tenant ID of the Schema.
* `description` - (Read-Only) The description of the Schema.
* `template` - (Read-Only) A list of templates for the Schema.
    * `name` - (Read-Only) The name of the Template.
    * `display_name` - (Read-Only) The name of the Template as displayed on the MSO UI.
    * `description` - (Read-Only) The description of the Template.
    * `tenant_id` - (Read-Only) The tenant ID of the Template.
    * `template_type` - (Read-Only) The type of the Template.

