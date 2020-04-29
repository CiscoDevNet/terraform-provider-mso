---
layout: "mso"
page_title: "MSO: mso_role"
sidebar_current: "docs-mso-resource-role"
description: |-
  Manages MSO Resource Role
---

# mso_role #

Manages MSO Role

## Example Usage ##

```hcl
resource "mso_role" "sample_role" {
  name = "UserManager"
  display_name = "UserManager"
  description = "hello"
  read_permissions = ["view-sites"]
  write_permissions = ["manage-sites","manage-tenants"]
  
}
```

## Argument Reference ##

* `name` - (Required) name of the role.
* `display_name` - (Required) Display name of the role.
* `read_permissions` - (Required) Read permissions assigned to the role.
Choices for read_permissions:
        "view-sites",
        "view-tenants",
        "view-schemas",
        "view-tenant-schemas",
        "view-users",
        "view-roles",
        "view-all-audit-records",
        "view-backup",
        "view-labels"
* `write_permissions` - (Required) Write permissions assigned to the role.
Choices for write_permissions:
        "manage-sites",
        "manage-tenants",
        "manage-labels",
        "manage-schemas",
        "manage-tenant-schemas",
        "manage-users",
        "manage-roles",
        "manage-audit-records",
        "manage-backup",
        "manage-labels"
* `description` - (Optional) Description of the role.

## Attribute Reference ##

No Attributes are Exported.
