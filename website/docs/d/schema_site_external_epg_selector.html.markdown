---
layout: "mso"
page_title: "MSO: mso_schema_site_external_epg_selector"
sidebar_current: "docs-mso-data-source-schema_site_external_epg_selector"
description: |-
  Data source for MSO Schema Site external Endpoint Groups Selector.
---

# mso_schema_site_external_epg_selector #

Data source for MSO Schema Site external Endpoint Groups Selector.

```hcl

data "mso_schema_site_external_epg_selector" "check"{
  schema_id         = "${mso_schema_template_external_epg.template_externalepg.schema_id}"
  template_name     = "${mso_schema_template_external_epg.template_externalepg.template_name}"
  site_id           = "${mso_schema_template_external_epg.template_externalepg.site_id}"
  external_epg_name = "${mso_schema_template_external_epg.template_externalepg.external_epg_name}"
  name              = "second"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy external Epg Selector.
* `site_id` - (Required) site ID under which you want to deploy external Epg Selector.
* `template_name` - (Required) Template under above site id where external Epg Selector to be created.
* `external_epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name for the selector.

## Attribute Reference ##

* `ip` - (Optional) Ip address associated with the selector.