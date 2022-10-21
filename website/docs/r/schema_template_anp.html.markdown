---
layout: "mso"
page_title: "MSO: mso_schema_template_anp"
sidebar_current: "docs-mso-resource-schema_template_anp"
description: |-
  Manages MSO Resource Schema Template Anp
---

# mso_schema_template_anp #

Manages MSO Resource Schema Template Anp

## Example Usage ##

```hcl

resource "mso_schema_template_anp" "anp1" {
  schema_id    = mso_schema.schema1.id
  template     = "template99"
  name         = "anp123"
  display_name = "anp1234"
}

```

## Argument Reference ##


* `schema_id` - (Required) The schema-id where anp is associated.
* `name` - (Required) name of the anp to add.
* `template` - (Required) template associated with the anp.
* `display_name` - (Required) The name as displayed on the MSO web interface.



## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Resource Schema Template Anp can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_anp.anp1 {schema_id}/template/{template}/anp/{name}
```

