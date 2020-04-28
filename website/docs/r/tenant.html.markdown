---
layout: "mso"
page_title: "MSO: mso_tenant"
sidebar_current: "docs-mso-resource-tenant"
description: |-
  Manages MSO Tenant
---

# schema #

Manages MSO Tenant

## Example Usage ##

```hcl
resource "mso_tenant" "tenant1" {
  name = "mso"
  display_name = "mso"
}
```

## Argument Reference ##

* `name` - (Required) The name of the tenant.
* `display_name` - (Required) The name of the tenant to be displayed in the web UI.

## Attribute Reference ##

* `description` - (Optional) The description for this tenant.
* `users` - (Optional) A list of associated users for this tenant.
* `sites` - (Optional) A list of associated sites for this tenant.
