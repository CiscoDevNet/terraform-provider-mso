---
layout: "mso"
page_title: "MSO: mso_schema_template"
sidebar_current: "docs-mso-data-source-schema_template"
description: |-
  Data source for MSO Schema Template.
---

# mso_schema_template #

Data source for MSO schema Template.

## Example Usage ##

```hcl

data "mso_schema_template" "example" {
  schema_id = data.mso_schema.schema1.id
  name      = "template-name"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Template.
* `name` - (Required) The name of the Template.

## Attribute Reference ##

* `tenant_id` - (Read-Only) The tenant ID to associate with the Template.
* `display_name` - (Read-Only) The name of the Template as displayed on the MSO UI.
* `template_type` - (Read-Only) The type of the Template.
* `description` - (Read-Only) The description of the Template.
