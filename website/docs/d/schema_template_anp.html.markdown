---
layout: "mso"
page_title: "MSO: mso_schema_template_anp"
sidebar_current: "docs-mso-data-source-schema_template_anp"
description: |-
  Data source for MSO Schema Template Anp
---

# mso_schema_template_anp #

Data source for MSO schema template Anp, to fetch the MSO schema template Anp details.

## Example Usage ##

```hcl
data "mso_schema_template_anp" "anp2" {
  schema_id="${mso_schema.schema1.id}"
  template= "template99"
  name = "anp123"
}
}
```

## Argument Reference ##

* `schema_id` - (Required) The schema-id where anp is associated.
* `name` - (Required) name of the anp to add.
* `template` - (Required) template associated with the anp.
* `display_name` - (Optional) The name as displayed on the MSO web interface.


## Attribute Reference ##

No attributes are exported.
