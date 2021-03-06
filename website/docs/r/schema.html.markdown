---
layout: "mso"
page_title: "MSO: mso_schema"
sidebar_current: "docs-mso-resource-schema"
description: |-
  Manages MSO Schema
---

# mso_schema #

Manages MSO Schema

## Example Usage ##

```hcl
resource "mso_schema" "foo_schema" {
  name          = "nkp12"
  template_name = "template1"
  tenant_id     = "5ea000bd2c000058f90a26ab"
}

```

## Argument Reference ##

* `name` - (Required) name of the schema.
* `template_name` - (Required) name of template attached to this schema.
* `tenant_id` - (Required) tenant_id for this schema.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of schema created.

## Importing ##

An existing MSO Schema can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema.schema1 {schema_id}
```