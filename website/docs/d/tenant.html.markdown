---
layout: "mso"
page_title: "MSO: mso_tenant"
sidebar_current: "docs-mso-data-source-tenant"
description: |-
  Data source for MSO Tenant
---

# mso_tenant #

Data source for MSO tenant

## Example Usage ##

```hcl
data "mso_tenant" "tenant1" {
  name = "mso"
}
```

## Argument Reference ##

 `name` - (Required) The name of the tenant.

## Attribute Reference ##

* `display_name` - (Required) The name of the tenant to be displayed in the web UI.
* `description` - (Optional) The description for this tenant.
* `users` - (Optional) A list of associated users for this tenant.
* `sites` - (Optional) A list of associated sites for this tenant.
