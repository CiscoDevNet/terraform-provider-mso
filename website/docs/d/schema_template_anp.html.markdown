---
layout: "mso"
page_title: "MSO: mso_schema_template_anp"
sidebar_current: "docs-mso-data-source-schema_template_anp"
description: |-
  Data source for MSO Schema Template Application Network Profile (ANP).
---

# mso_schema_template_anp #

Data source for MSO Schema Template Application Network Profile (ANP).

## Example Usage ##

```hcl

data "mso_schema_template_anp" "example" {
  schema_id = data.mso_schema.schema1.id
  template  = "template99"
  name      = "anp123"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the ANP.
* `template` - (Required) The template name of the ANP.
* `name` - (Required) The name of the ANP.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the ANP as displayed on the MSO UI.
