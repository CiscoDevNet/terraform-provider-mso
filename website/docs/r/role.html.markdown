---
layout: "mso"
page_title: "MSO: mso_role"
sidebar_current: "docs-mso-resource-role"
description: |-
  Manages MSO Resource Role
---

# schema #

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
* `description` - (Required) Description of the role.
* `read_permissions` - (Required) Read permissions assigned to the role.
* `write_permissions` - (Required) Write permissions assigned to the role.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of site associated.