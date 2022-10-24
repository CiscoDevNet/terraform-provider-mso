---
layout: "mso"
page_title: "MSO: mso_user"
sidebar_current: "docs-mso-resource-user"
description: |-
  Manages MSO User
---

# mso_user #

Manages MSO User.

The mso_user resources is deprecated on ND-based MSO/NDO.
Use ND provider for manipulating users on ND-based MSO/NDO.

## Example Usage ##

```hcl

resource "mso_user" "user1" {
  username       = "name1"
  user_password  = "password"
  first_name     = "first_name"
  last_name      = "last_name"
  email          = "email@gmail.com"
  phone          = "12345678910"
  account_status = "active"
  roles {
    roleid = "0000ffff0000000000000031"
  }
}

```

## Argument Reference ##

* `username` - (Required) username of the user. It must contain at least 1 character in length.
* `user_password` - (Required) password of the user. It must contain at least 8 characters in length.
* `roles` - **Deprecated** (Optional) roles given to the user. This attribute is deprecated on ND-based MSO/NDO, use `user_rbac` instead.
  * `roles.roleid` - (Optional) id of roles given to the user.
  * `roles.access_type` - (Optional) access_type of roles given to the user.
* `user_rbac` - (Optional) roles given to the user.
  * `user_rbac.name` - (Optional) name of roles given to the user.
  * `user_rbac.user_priv` - (Optional) Privilege access given to users (WritePriv, ReadPriv)
* `first_name` - (Optional) firstname of the user.
* `last_name` - (Optional) lastname of the user.
* `email` - (Optional) email of the user.
* `phone` - (Optional) phone of the user.
* `account-status` - (Optional) account status of the user.
* `domain` - (Optional) domain status of the user.
* `roles.access_type` - (Optional) access_type of roles given to the user.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of user associated.

## Importing ##

An existing MSO User can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_user.user1 {user_id}
```