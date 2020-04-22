---
layout: "mso"
page_title: "Provider: MSO"
sidebar_current: "docs-mso-index"
description: |-
  Cisco ACI Multi-Site Orchestrator (MSO) is responsible for provisioning, health monitoring, and managing the full lifecycle of Cisco ACI networking policies and tenant policies across Cisco ACI sites around the world. Terraform provider MSO is a Terraform plugin which will be used to manage the MSO Fabric Constructs on the Cisco MSO platform with leveraging advantages of Terraform. The provider needs to be configured with the proper credentials before it can be used.
---

MSO Provider
------------

Cisco ACI Multi-Site Orchestrator (MSO) is responsible for provisioning, health monitoring, and managing the full lifecycle of Cisco ACI networking policies and tenant policies across Cisco ACI sites around the world. Terraform provider MSO is a Terraform plugin which will be used to manage the MSO Fabric Constructs on the Cisco MSO platform with leveraging advantages of Terraform. The provider needs to be configured with the proper credentials before it can be used.

Authentication
--------------

Authentication with user-id and password.

Example Usage
------------

 ```hcl
 provider "mso" {
    # cisco-mso user name
    username = "admin"
    # cisco-mso password
    password = "password"
    # cisco-mso url
    url      = "https://173.36.219.193/"
    insecure = true
}

resource "mso_schema" "foo_schema" {
  name          = "nkp12"
  template_name = "template1"
  tenant_id     = "5ea000bd2c000058f90a26ab"
}
```

Argument Reference
------------------

Following arguments are supported with Cisco MSO terraform provider.

* `username` - (Required) This is the Cisco MSO username, which is required to authenticate with CISCO MSO.
* `password` - (Required) Password of the user mentioned in username argument. It is required when you want to use token basedauthentication.
* `url` - (Required) URL for CISCO MSO.
* `insecure` - (Optional) This determines whether to use insecure HTTP connection or not. Default value is `true`.
