---
layout: "mso"
page_title: "MSO: mso_schema_template"
sidebar_current: "docs-mso-resource-schema_template"
description: |-
  Manages MSO Schema Template
---

# mso_schema_template #

Manages MSO Schema Template

## Example Usage ##

```hcl

resource "mso_schema_template" "st1" {
  schema_id    = mso_schema.schema1.id
  name         = "Temp1"
  display_name = "Temp1"
  tenant_id    = mso_tenant.tenant1.id
}

```

## Argument Reference ##

* `schema_id` - (Required) name of the schema.
* `tenant_id` - (Required) Tenant-id to associate.
* `name` - (Required) Name of the template.
* `display_name` - (Required) Display name of the Template to be deployed on the site.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of schema template associated.

## Importing ##

An existing MSO Schema Template can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template.st1 {schema_id}/template/{name}
```
