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

resource "mso_schema_template" "demo_template" {
  schema_id     = mso_schema.demo_schema.id
  name          = "Template1"
  display_name  = "Template1"
  tenant_id     = mso_tenant.demo_tenant.id
  template_type = "aci_multi_cloud"
}

```

## Argument Reference ##

* `name` - (Required) The name of the template.
* `schema_id` - (Required) The schema-id where template is associated.
* `tenant_id` - (Required) The tenant-id to associate with the template.
* `display_name` - (Required) The display name of the template.
* `template_type` - (Required) The template type of the template.
* `description` - (Optional) The description of the template.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of schema template associated.

## Importing ##

An existing MSO Schema Template can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template.st1 {schema_id}/template/{name}
```
