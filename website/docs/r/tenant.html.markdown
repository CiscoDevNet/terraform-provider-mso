---
layout: "mso"
page_title: "MSO: mso_tenant"
sidebar_current: "docs-mso-resource-tenant"
description: |-
  Manages MSO Tenant
---

# mso_tenant #

Manages MSO Tenant

## Example Usage ##

```hcl
resource "mso_tenant" "tenant1" {
  name = "m3"
  display_name = "m3"
  site_associations{site_id = "5c7c95b25100008f01c1ee3c"}
  user_associations{user_id = "0000ffff0000000000000020"}
}
```

## Argument Reference ##

* `name` - (Required) The name of the tenant.
* `display_name` - (Required) The name of the tenant to be displayed in the web UI.
* `description` - (Optional) The description for this tenant.
* `user_associations` - (Optional) A list of associated users for this tenant.
* `site_association` - (Optional) A list of associated sites for this tenant.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of tenant created.
