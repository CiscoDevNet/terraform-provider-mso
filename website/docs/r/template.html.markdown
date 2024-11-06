---
layout: "mso"
page_title: "MSO: mso_template"
sidebar_current: "docs-mso-resource-template"
description: |-
  Manages MSO Template
---

# mso_template #

Manages MSO Template

## Example Usage ##

```hcl

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id     = data.mso_tenant.example_tenant.id
  sites         = [data.mso_site.site_1.id, data.mso_site.site_2.id]
}

```

## Argument Reference ##

* `template_name` - (Required) The name of the template.
* `template_type` - (Required) The type of the template. Allowed values are `tenant`, `l3out`, `fabric_policy`, `fabric_resource`, `monitoring_tenant`, `monitoring_access`, or `service_device`.
* `tenant_id` - (Optional) The ID of the tenant to associate with the template.
* `sites` - (Optional) A list of site IDs to associate with the template.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the ID of the template.

## Importing ##

An existing MSO Schema Template can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_template.tenant_template {id}
```
