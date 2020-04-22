---
layout: "mso"
page_title: "MSO: mso_schema"
sidebar_current: "docs-mso-resource-schema"
description: |-
  Manages MSO Schema
---

# schema #

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
* `template_name` - (optional) name of templates for this schema.
* `tenant_id` - (optional) temant_id for this schema.

## Attribute Reference ##

No attributes are exported
