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
  anp_name = "AP1234"
  schema_id = "5c6c16d7270000c710f8094d"
  site_id = "5c7c95d9510000cf01c1ee3d"
  template_name = "Template1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Anp.
* `template_name` - (Required) Template where Site Anp to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp.
* `anp_name` - (Required) Name of Site Anp.

## Attribute Reference ##

* `anp_schema_id` - (Optional) SchemaID of Anp. schema_id will be used if not provided. Should use this parameter when Anp is deployed to a different site.
* `anp_template_name` - (Optional) Template Name of Anp. template_name will be used if not provided.
