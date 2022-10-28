---
layout: "mso"
page_title: "MSO: mso_schema_site_anp"
sidebar_current: "docs-mso-data-source-schema_site_anp"
description: |-
  MSO Schema Site Application Network Profile(ANP) Data source.
---

# mso_schema_site_anp #

 MSO Schema Site Application Network Profile(ANP) Data source.

## Example Usage ##

```hcl

data "mso_schema_site_anp" "st10" {
  anp_name      = "anp1"
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Anp.
* `site_id` - (Required) SiteID under which you want to deploy Anp.
* `anp_name` - (Required) Name of Site Anp. The name of the ANP should be present in the ANP list of the given `schema_id` and `template_name`

## Attribute Reference ##

* `template_name` - (Optional) Template where Site Anp to be created.
