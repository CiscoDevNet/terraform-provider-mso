---
layout: "mso"
page_title: "MSO: mso_schema_site_external_epg_selector"
sidebar_current: "docs-mso-data-source-schema_site_external_epg_selector"
description: |-
  Data source for MSO Schema Site External End Point Group Selector.
---

# mso_schema_site_external_epg_selector #

Data source for MSO Schema Site External End Point Group Selector.

```hcl

data "mso_schema_site_external_epg_selector" "example"{
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  site_id           = data.mso_site.site1.id
  external_epg_name = "external_epg1"
  name              = "second"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Selector is deployed.
* `site_id` - (Required) The site ID under which the Selector is deployed.
* `template_name` - (Required) The template name under which the Selector is deployed.
* `external_epg_name` - (Required) The name of the External EPG under which the Selector is deployed.
* `name` - (Required) The name of the Selector.

## Attribute Reference ##

* `ip` - (Read-Only) The IP address of the Selector.
