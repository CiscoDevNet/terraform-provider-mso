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
  display_name = "mso"
}
```

## Argument Reference ##

* `name` - (Required) The name of the tenant.
* `display_name` - (Required) The name of the tenant to be displayed in the web UI.

## Attribute Reference ##

* `description` - (Optional) The description for this tenant.
* `user_associations` - (Optional) A list of associated users for this tenant.
* `site_association` - (Optional) A list of associated sites for this tenant.
