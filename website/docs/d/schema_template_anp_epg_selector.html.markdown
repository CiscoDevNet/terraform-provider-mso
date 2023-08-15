---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_selector"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_selector"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Group Selector.
---

# mso_schema_template_anp_epg_selector #

Data source for MSO Schema Template Application Network Profiles Endpoint Group Selector.

```hcl

data "mso_schema_template_anp_epg_selector" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  anp_name      = "anp1"
  epg_name      = "epg1
  name          = "subnet1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Selector.
* `template_name` - (Required) The template name of the Selector.
* `anp_name` - (Required) The name of the ANP.
* `epg_name` - (Required) The name of the EPG.
* `name` - (Required) The name of the Selector.

## Attribute Reference ##

* `expressions` - (Read-Only) A list of expressions for the Selector.
    * `key` - (Read-Only) The key of the Selector expression.
    * `operator` - (Read-Only) The operator of the Selector expression.
    * `value` - (Read-Only) The value of the Selector expression.
