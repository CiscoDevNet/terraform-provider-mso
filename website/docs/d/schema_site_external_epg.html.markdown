---
layout: "mso"
page_title: "MSO: mso_schema_site_external_epg"
sidebar_current: "docs-mso-data-source-schema_site_external_epg"
description: |-
  Data source for MSO Schema Site External End Point Group.
---

# mso_schema_site_external_epg #

Data source for MSO Schema Site External End Point Group.

```hcl

data "mso_schema_site_external_epg" "example" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  site_id           = data.mso_site.site1.id
  external_epg_name = "external_epg_1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the External EPG is deployed.
* `site_id` - (Required) The site ID under which the External EPG is deployed.
* `template_name` - (Required) The template name under which the External EPG is deployed.
* `external_epg_name` - (Required) The name of the External EPG.

## Attribute Reference ##

* `l3out_name` - (Read-Only) The name of the L3Out.
* `l3out_schema_id` - (Read-Only) The schema ID of the L3out.
* `l3out_template_name` - (Read-Only) The template name of the L3out.
* `l3out_dn` - (Read-Only) The DN of the L3out.
