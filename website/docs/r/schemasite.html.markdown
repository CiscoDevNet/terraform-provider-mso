---
layout: "mso"
page_title: "MSO: mso_schema_site"
sidebar_current: "docs-mso-resource-schema_site"
description: |-
  Manages MSO Schema Site
---

# schema #

Manages MSO Schema Site

## Example Usage ##

```hcl
resource "mso_schema_site" "foo_schema_site" {
  schema  = "schema1"
  site  = "bdsol-pod51"
  template  = "template1"
}
```

## Argument Reference ##

* `schema` - (Required) name of the schema.
* `template` - (Required) The name of the template.
* `site` - (Required) The name of the site to manage.

## Attribute Reference ##

No attributes are exported
