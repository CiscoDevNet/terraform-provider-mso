---
layout: "mso"
page_title: "MSO: mso_schema_site_anp"
sidebar_current: "docs-mso-data-source-schema_site_anp"
description: |-
  Data source for MSO Schema Site Application Network Profile (ANP).
---

# mso_schema_site_anp #

 Data source for MSO Schema Site Application Network Profile (ANP).

## Example Usage ##

```hcl

data "mso_schema_site_anp" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  anp_name      = "anp1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the ANP is deployed.
* `site_id` - (Required) The site ID under which the ANP is deployed.
* `template_name` - (Required) The template name under which the ANP is deployed.
* `anp_name` - (Required) The name of the ANP.
