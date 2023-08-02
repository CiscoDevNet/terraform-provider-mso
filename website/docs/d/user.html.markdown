---
layout: "mso"
page_title: "MSO: mso_user"
sidebar_current: "docs-mso-resource-user"
description: |-
  Data source for MSO User.
---

# mso_user #

Data source for MSO User.

## Example Usage ##

```hcl

data "mso_user" "example" {
  username = "name"
}

```

## Argument Reference ##

* `username` - (Required) The username of the User.

## Attribute Reference ##

* `user_password` - (Read-Only) The password of the User.
* `roles` - **Deprecated** (Read-Only) The roles of the User. This attribute is deprecated on ND-based MSO/NDO.
    * `roleid` - (Read-Only) The role ID of the User.
    * `access_type` - (Read-Only) The acces type of the User.
* `user_rbac` - (Read-Only) The roles of the User. 
    * `name` - (Read-Only) The name of the role. 
    * `user_priv` - (Read-Only) The privilege access of the User.
* `first_name` - (Read-Only) The first name of the User.
* `last_name` - (Read-Only) The last name of the User.
* `email` - (Read-Only) The email of the User.
* `phone` - (Read-Only) The phone of the User. This attribute is deprecated on ND-based MSO/NDO.
* `account-status` - (Read-Only) The account status of the User.
* `domain` - (Read-Only) The domain status of the User.
