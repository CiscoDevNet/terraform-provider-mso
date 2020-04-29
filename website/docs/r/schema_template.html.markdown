---
layout: "mso"
page_title: "MSO: mso_schema_template"
sidebar_current: "docs-mso-resource-schema_template"
description: |-
  Manages MSO Schema Template
---

# mso_schema_template #

Manages MSO Schema Site

## Example Usage ##

```hcl
resource "mso_schema_template" "st1" {
  schema_id = "${mso_schema.s01.id}"
  name = "Temp1"
  display_name = "Temp1"
  tenant_id = "5c4d9f3d2700007e01f80949"
}
```

## Argument Reference ##

* `schema_id` - (Required) name of the schema.
* `tenant_id` - (Required) Tenant-id to associate.
* `name` - (Required) Name of the template.
* `display_name` - (Required) Display name of the Template to be deployed on the site.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of schema template associated.
