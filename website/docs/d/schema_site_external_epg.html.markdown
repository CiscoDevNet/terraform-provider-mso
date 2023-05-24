---
layout: "mso"
page_title: "MSO: mso_schema_site_external_epg"
sidebar_current: "docs-mso-data-source-schema_site_external_epg"
description: |-
  Data source for MSO Schema Site External Endpoint Groups.
---

# mso_schema_site_external_epg #

Data source for MSO Schema Site External Endpoint Groups.

```hcl

data "mso_schema_site_external_epg" "external_epg_1" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  site_id           = data.mso_site.site1.id
  external_epg_name = "external_epg_1"
}

```

## Argument Reference ##

* `schema_id` - (Required) Schema ID under which you want to deploy the External Endpoint Group.
* `site_id` - (Required) Site ID under which you want to deploy the External Endpoint Group.
* `template_name` - (Required) Template Name under which you want to define the External Endpoint Group.
* `external_epg_name` - (Required) Name of the External Endpoint Group.

## Attribute Reference ##

* `l3out_name` - (Read-Only) Name of the L3Out.