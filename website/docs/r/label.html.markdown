---
layout: "mso"
page_title: "MSO: mso_label"
sidebar_current: "docs-mso-resource-label"
description: |-
  Manages MSO Resource Label
---

# mso_label #

Manages MSO Label

## Example Usage ##

```hcl
 resource "mso_label" "label1" {
   label = "hello3"
   type  = "site"
 }
```

## Argument Reference ##

* `label` - (Required) name of the label.
* `type` - (Required) type of the label.

## Attribute Reference ##

No Attributes are Exported.
