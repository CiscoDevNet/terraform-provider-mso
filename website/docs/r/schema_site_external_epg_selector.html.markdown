---
layout: "mso"
page_title: "MSO: mso_schema_site_external_epg_selector"
sidebar_current: "docs-mso-resource-schema_site_external_epg_selector"
description: |-
  Manages MSO Schema site external Endpoint Groups selectors.
---

# mso_schema_site_external_epg_selector#

Manages MSO Schema site external Endpoint Groups Selectors.

## Example Usage ##
```hcl

resource "mso_schema_site_external_epg_selector" "sel1" {
  schema_id         = "${mso_schema_template_external_epg.template_externalepg.schema_id}"
  template_name     = "${mso_schema_template_external_epg.template_externalepg.template_name}"
  site_id           = "${mso_schema_template_external_epg.template_externalepg.site_id}"
  external_epg_name = "${mso_schema_template_external_epg.template_externalepg.external_epg_name}"
  name              = "second"
  ip                = "12.25.70.50"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy External Epg Selector.
* `site_id` - (Required) site ID under which you want to deploy External Epg Selector.
* `template_name` - (Required) Template under above site id where External Epg Selector to be created.
* `external_epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name for the selector.
* `ip` - (Required) Ip address associated with the selector.

## Attribute Reference ##

No attributes are exported.