---
layout: "mso"
page_title: "MSO: mso_rest"
sidebar_current: "docs-mso-resource-rest"
description: |-
  MSO Rest Resource to manage the MSO objects via REST API.
---

# mso_rest #

MSO Rest Resource to manage the MSO objects via REST API.

## Example Usage ##

```hcl

resource "mso_rest" "sample_rest" {
    path = "api/v1/schemas/5ebb9f682c0000da45812937"
    method = "PATCH"
    payload = <<EOF
[
  {
    "op": "remove",
    "path": "/templates/Template3",
    "value": {
      "name": "Template3",
      "displayName": "Templat2",
      "tenantId": "5e9d09482c000068500a269a"
    }
  }
]
EOF  
}

```

## Argument Reference ##

* `path` - (Required) MSO REST endpoint, where the data is being sent.
* `method` - (Optional) HTTP method, allowed values are POST, PATCH, GET, DELETE and PUT.
* `payload` - (Required) JSON encoded payload data.

NOTE: This resource will not work well in the case of Terraform destroy if there is a change in the terraform configuration required to destroy the object from the MSO, as Destroy only has the access to the data in the state file. To destroy the objects created via mso_rest in such cases modify the payload and use the Terraform apply instead.

## Attribute Reference ##

No Attributes are Exported.
