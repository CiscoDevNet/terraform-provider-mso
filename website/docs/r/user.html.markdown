---
layout: "mso"
page_title: "MSO: mso_user"
sidebar_current: "docs-mso-resource-user"
description: |-
  Manages MSO User
---

# user #

Manages MSO User

## Example Usage ##

```hcl
resource "mso_user" "user1" {
  username      = "name1"
  user_password  = "password"
  first_name="first_name"
  last_name="last_name"
  email="email@gmail.com"     
  phone="12345678910"
  account_status="active"
  roles{
    roleid="0000ffff0000000000000031"
}
}

```

## Argument Reference ##

* `username` - (Required) username of the user.
* `user_password` - (Required) password of the user.
* `first_name` - (Optional) firstname of the user.
* `last_name` - (Optional) lastname of the user.
* `email` - (Optional) email of the user.
* `phone` - (Optional) phone of the user.
* `account-status` - (Optional) account status of the user.
* `domain` - (Optional) domain status of the user.
* `roles` - (Required) roles given to the user.
* `roles.roleid` - (Required) id of roles given to the user.
* `roles.access_type` - (Optional) access_type of roles given to the user.



## Attribute Reference ##

No attributes are exported