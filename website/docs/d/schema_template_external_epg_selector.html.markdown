---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_selector"
sidebar_current: "docs-mso-data-source-schema_template_external_epg_selector"
description: |-
  Data source for MSO Schema Template External End Point Group Selector.
---

# mso_schema_template_external_epg_selector #

!> **Deprecated** This data source is deprecated: cloud-specific features are no longer supported in Nexus Dashboard 4.x (NDO 5.x) releases. Ensure the related configuration is removed from the NDO schema before upgrading Nexus Dashboard to 4.x (NDO 5.x).

Data source for MSO Schema Template External End Point Group Selector.

```hcl

data "mso_schema_template_external_epg_selector" "example" {
  schema_id          = data.mso_schema.schema1.id
  template_name      = "Template1"
  external_epg_name  = "epg1"
  name               = "check"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the External EPG.
* `template_name` - (Required) The template name of the External EPG.
* `external_epg_name` - (Required) The name of the External EPG.
* `name` - (Required) The name of the Selector.

## Attribute Reference ##

* `expressions` - (Read-Only) A list of expressions for the Selector.
    * `key` - (Read-Only) The key of the Selector expression.
    * `operator` - (Read-Only) The operator of the Selector expression.
    * `value` - (Read-Only) The value of the Selector expression.
