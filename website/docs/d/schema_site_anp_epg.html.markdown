---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg"
description: |-
  Data source for MSO Schema Site ANP End Point Group (EPG).
---

# mso_schema_site_anp_epg #

Data source for MSO Schema Site ANP End Point Group (EPG).

## Example Usage ##

```hcl

data "mso_schema_site_anp_epg" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  site_id       = data.mso_site.site1.id
  anp_name      = "ANP"
  epg_name      = "DB"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Subnet is deployed.
* `site_id` - (Required) The site ID under which the Subnet is deployed.
* `template_name` - (Required) The template name under which the Subnet is deployed.
* `anp_name` - (Required) The name of the ANP.
* `epg_name` - (Required) The name of the EPG.
