---
layout: "mso"
page_title: "MSO: mso_schema_template"
sidebar_current: "docs-mso-data-source-schema_template"
description: |-
  Data source for MSO Schema Template
---

# mso_schema_template #

Data source for MSO schema template, to fetch the MSO schema template details.

## Example Usage ##

```hcl

data "mso_schema_template" "demo_template" {
  name      = "template-name"
  schema_id = data.mso_schema.schema1.id
}

```

## Argument Reference ##

* `name` - (Required) The name of the template.
* `schema_id` - (Required) The schema-id where template is associated.

## Attribute Reference ##

* `tenant_id` - (Read-Only) The tenant-id to associate with the template.
* `display_name` - (Read-Only) The display name of the template.
* `template_type` - (Read-Only) The template type of the template.
* `description` - (Read-Only) The description of the template.
