---
layout: "mso"
page_title: "MSO: mso_schema_site_anp"
sidebar_current: "docs-mso-resource-schema_site_anp"
description: |-
  MSO Schema Site Application Network Profile(ANP) Resource
---

# mso_schema_site_anp #

 MSO Schema Site Application Network Profile(ANP) Resource.

## Example Usage ##

```hcl
resource "mso_schema_site_anp" "anp1" {
  schema_id = "5c6c16d7270000c710f8094d"
  anp_name = "AP1234"
  template_name = "Template1"
  site_id = "5c7c95d9510000cf01c1ee3d"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Site Anp.
* `template_name` - (Required) Template where Site Anp to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp.
* `anp_name` - (Required) Name of Site Anp.  The name of the ANP should be present in the ANP list of the given `schema_id` and `template_name`

## Attribute Reference ##

No attributes are exported.