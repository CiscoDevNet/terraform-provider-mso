---
layout: "mso"
page_title: "MSO: mso_schema_site"
sidebar_current: "docs-mso-data-source-schema_site"
description: |-
  Data source for MSO Schema Site.
---

# mso_schema_site #

Data source for MSO Schema Site.

## Example Usage ##

```hcl

data "mso_schema_site" "example" {
  name          = "sitename"
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
}

```

## Argument Reference ##

* `name` - (Required) The name of the Site.
* `schema_id` - (Required) The schema ID where the Site is associated.
* `template_name` - (Required) The name of the template attached to the Site.

## Attribute Reference ##

* `site_id` - (Read-Only) The ID (UUID) of the Site.
