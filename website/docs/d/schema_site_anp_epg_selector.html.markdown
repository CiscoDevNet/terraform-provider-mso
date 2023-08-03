---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_selector"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_selector"
description: |-
  Data source for MSO Schema Site Application Network Profiles End Point Group Selector.
---

# mso_schema_site_anp_epg_selector #

Data source for MSO Schema Site Application Network Profiles End Point Group Selector.

```hcl

data "mso_schema_site_anp_epg_selector" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  anp_name      = "anp1"
  epg_name      = "epg1"
  name          = "check01"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Selector is deployed.
* `site_id` - (Required) The site ID under which the Selector is deployed.
* `template_name` - (Required) The template name under which the Selector is deployed.
* `anp_name` - (Required) The ANP name under which the Selector is deployed.
* `epg_name` - (Required) The EPG name under which the Selector is deployed.
* `name` - (Required) The name for the Selector.

## Attribute Reference ##

* `expressions` - (Read-Only) A list of expressions for the Selector.
    * `key` - (Read-Only) The key of the Selector expression.
    * `operator` - (Read-Only) The operator of the Selector expression.
    * `value` - (Read-Only) The value of the Selector expression.
