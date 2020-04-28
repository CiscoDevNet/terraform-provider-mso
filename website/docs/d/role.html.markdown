---
layout: "mso"
page_title: "MSO: mso_role"
sidebar_current: "docs-mso-data-source-role"
description: |-
  Data source for MSO Role
---

# mso_schema #

Data source for MSO role  

## Example Usage ##

```hcl
data "mso_role" "sample_role" {
  name  = "UserManager"
}
```

## Argument Reference ##

* `name` - (Required) name of the schema.

## Attribute Reference ##

* `display_name` - (Optional) Name displayed associated to this Role.
* `description` - (Optional) Description for this role.
* `read_permissions` - (Optional) Read permissions attached to this role.
* `write_permissions` - (Optional) Write permissions attached to this role.
