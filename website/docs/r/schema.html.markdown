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

### When template blocks are provided ###

```hcl
resource "mso_schema" "schema1" {
  name          = "Schema3"
  template {
    name          = "Template1"
    display_name  = "TEMP1"
    tenant_id     = "623316531d0000abdd50343a"
  }
  template {
    name          = "Template2"
    display_name  = "TEMP2"
    tenant_id     = "623316531d0000abdd50343a"
  }
  template {
    name          = "Template3"
    display_name  = "TEMP3"
    tenant_id     = "0000ffff0000000000000010"
  }
}  

```

### When template_name and tenant_id are used(DEPRECATED) ###

```hcl
resource "mso_schema" "schema1" {
  name          = "nkp12"
  template_name = "template1"
  tenant_id     = "5ea000bd2c000058f90a26ab"
}

```

## Argument Reference ##

* `name` - (Required) name of the schema.
* `template_name` - (Optional) **Deprecated**. Name of template attached to the schema.
* `tenant_id` - (Optional) **Deprecated**. tenant_id for the schema.
* `template` - (Optional) A block that represents the template associated with the schema. Multiple templates can be created using this attribute. Type - Block.
  * `name` - Name of template.
  * `display_name` - Display name for the template.
  * `tenant_id` - tenant_id for the template.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of schema created.

## Importing ##

An existing MSO Schema can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema.schema1 {schema_id}
```