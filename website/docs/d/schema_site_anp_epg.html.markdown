---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg"
description: |-
  Data source for MSO Schema Site Application Network Profiles Endpoint Groups.
---

# mso_schema_site_anp_epg #

Data source for MSO Schema Site Application Network Profiles Endpoint Groups.

## Example Usage ##

```hcl

data "mso_schema_site_anp_epg" "anpEpg" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  site_id       = data.mso_site.site1.id
  anp_name      = "ANP"
  epg_name      = "DB"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg.

* `site_id` - (Required) SiteID under which you want to deploy Anp Epg.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.

## Attribute Reference ##

* `template_name` - (Optional) Template where Anp Epg to be created.
