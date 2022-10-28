---
layout: "mso"
page_title: "MSO: mso_schema_site"
sidebar_current: "docs-mso-resource-schema_site"
description: |-
  Manages MSO Schema Site
---

# mso_schema_site #

Manages MSO Schema Site

## Example Usage ##

```hcl

resource "mso_schema_site" "schema_site" {
  schema_id       =  mso_schema.schema1.id
  site_id         =  mso_site.site_test_1.id
  template_name   =  "template1"
}

```
NOTE: To add multiple sites, the tenant associated to schema must also be associated to the sites.

## Example usage to add multiple sites to the same schema ##

```hcl
resource "mso_schema_site" "foo_schema_site" {
  schema_id       =  mso_schema.schema1.id
  site_id         =  mso_site.site_test_1.id
  template_name   =  "template1"
}

resource "mso_schema_site" "foo_schema_site_2" {
  schema_id       =  mso_schema.schema1.id
  site_id         =  mso_site.site_test_2.id
  template_name   =  "template1"
}
```

## Argument Reference ##

* `schema_id`     - (Required) name of the schema.
* `site_id`       - (Required) Site-id to associate.
* `template_name` - (Required) Template to be deployed on the site.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of schema site associated.

## Importing ##

An existing MSO Schema Site can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site.site1 {schema_id}/site/{site_name}
```