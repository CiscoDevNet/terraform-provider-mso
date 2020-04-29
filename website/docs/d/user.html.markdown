---
layout: "mso"
page_title: "MSO: mso_user"
sidebar_current: "docs-mso-resource-user"
description: |-
  Data source for MSO User
---

# mso_user #

Data source for MSO User

## Example Usage ##

```hcl
data "mso_user" "schema10" {
  username = "name"
}
```

## Argument Reference ##

* `username` - (Required) username of the user. It must contain at least 1 character in length.

## Attribute Reference ##

* `user_password` - (Optional) password of the user. It must contain at least 8 characters in length.
* `roles` - (Optional) roles given to the user.
* `roles.roleid` - (Optional) id of roles given to the user.
* `first_name` - (Optional) firstname of the user.
* `last_name` - (Optional) lastname of the user.
* `email` - (Optional) email of the user.
* `phone` - (Optional) phone of the user.
* `account-status` - (Optional) account status of the user.
* `domain` - (Optional) domain status of the user.
* `roles.access_type` - (Optional) access_type of roles given to the user.
