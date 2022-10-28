---
layout: "mso"
page_title: "MSO: mso_role"
sidebar_current: "docs-mso-data-source-role"
description: |-
  Data source for MSO Role
---

# mso_role #

Data source for MSO role  

## Example Usage ##

```hcl

data "mso_role" "role" {
  name  = "UserManager"
}

```

## Argument Reference ##

* `name` - (Required) name of the schema.

## Attribute Reference ##

* `display_name` - (Optional) Name displayed associated to this Role.
* `read_permissions` - (Optional) Read permissions assigned to the role.
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
* `write_permissions` - (Optional) Write permissions assigned to the role.
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
* `description` - (Optional) Description for this role.
