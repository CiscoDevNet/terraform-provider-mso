---
layout: "mso"
page_title: "MSO: mso_rest"
sidebar_current: "docs-mso-data-source-rest"
description: |-
  Data source for reading MSO objects via REST API.
---

# mso_rest #

Data source for reading MSO objects via REST API.

## Example Usage ##

```hcl

data "mso_rest" "example" {
  path = "api/v1/platform/systemConfig"
}

```

## Argument Reference ##

* `path` - (Required) The MSO REST endpoint, where the data is being read.

## Attribute Reference ##

* `content` - (Read-Only) JSON response as a string.
