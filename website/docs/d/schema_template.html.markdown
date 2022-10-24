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

data "mso_schema_template" "st10" {
  name      = "template-name"
  schema_id = data.mso_schema.schema1.id
}

```

## Argument Reference ##

* `name` - (Required) name of the template to fetch.
* `schema_id` - (Required) The schema-id where template is associated.

## Attribute Reference ##

* `tenant_id` - (Optional) Tenant id is set to the MSO template UUID.
* `display_name` - (Optional) The display name of the template deployed to the site
