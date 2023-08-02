---
layout: "mso"
page_title: "MSO: mso_label"
sidebar_current: "docs-mso-data-source-label;"
description: |-
  Data source for MSO Label
---

# mso_label #

Data source for MSO Label

## Example Usage ##

```hcl

data "mso_label" "example" {
  label = "hello3"
}

```

## Argument Reference ##

* `label` - (Required) The name of the Label.

## Attribute Reference ##

* `type` - (Read-Only) The type of the Label.
