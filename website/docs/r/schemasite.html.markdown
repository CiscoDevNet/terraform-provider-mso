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
  schema_id  = "${mso_schema.schema1.id}"
  site_id  = "bdsol-pod51"
  template_name  = "template1"
}
```

## Argument Reference ##

* `schema_id` - (Required) name of the schema.
* `site_id` - (Required) Site-id to associate.
* `template_name` - (Required) Template to be deployed on the site.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of site associated.