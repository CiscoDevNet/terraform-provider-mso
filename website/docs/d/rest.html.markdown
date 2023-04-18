---
layout: "mso"
page_title: "MSO: mso_rest"
sidebar_current: "docs-mso-data-source-rest"
description: |-
  MSO Rest data source to read MSO objects via REST API.
---

# mso_rest #

MSO Rest data source to read MSO objects via REST API.

## Example Usage ##

```hcl
data "mso_rest" "system_config" {
  path = "api/v1/platform/systemConfig"
}
```

## Argument Reference ##

* `path` - (Required) MSO REST endpoint, where the data is being read.

## Attribute Reference ##

* `content` - (Read-Only) JSON response as a string.
