---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg"
sidebar_current: "docs-mso-resource-schema_site_anp_epg"
description: |-
  Manages MSO Schema Site Application Network Profiles Endpoint Groups.
---

# mso_schema_site_anp_epg #

Manages MSO Schema Site Application Network Profiles Endpoint Groups.

## Example Usage ##

```hcl
resource "mso_schema_site_anp_epg" "site_anp_epg" {
  schema_id = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id = "5c7c95d9510000cf01c1ee3d"
  anp_name = "ANP"
  epg_name = "DB"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg.
* `template_name` - (Required) Template where Anp Epg to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp Epg.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.

## Attribute Reference ##

No attributes are exported.