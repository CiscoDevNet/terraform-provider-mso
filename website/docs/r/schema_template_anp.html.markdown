---
layout: "mso"
page_title: "MSO: mso_schema_template_anp"
sidebar_current: "docs-mso-resource-schema_template_anp"
description: |-
  Manages MSO Resource Schema Template ANP
---

# mso_schema_template_anp #

Manages MSO Resource Schema Template ANP

## Example Usage ##

```hcl

resource "mso_schema_template_anp" "anp1" {
  schema_id    = mso_schema.schema1.id
  template     = mso_schema_template.st1.name
  name         = "anp123"
  display_name = "anp1234"
}

```

## Argument Reference ##


* `schema_id` - (Required) The schema-id where ANP is associated.
* `name` - (Required) Name of the ANP to add.
* `template` - (Required) Template associated with the ANP.
* `display_name` - (Required) The name as displayed on the MSO web interface.
* `description` - (Optional) The description of the ANP.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Resource Schema Template ANP can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_anp.anp1 {schema_id}/template/{template}/anp/{name}
```

