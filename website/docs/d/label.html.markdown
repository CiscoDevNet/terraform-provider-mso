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
data "mso_label" "label1" {
  label = "hello3"
}
```

## Argument Reference ##

* `label` - (Required) name of the label.

## Attribute Reference ##

* `type` - (Optional) type of the label.
