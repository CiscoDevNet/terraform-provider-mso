---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_selector"
sidebar_current: "docs-mso-resource-schema_site_anp_epg_selector"
description: |-
  Manages MSO Schema site Application Network Profiles Endpoint Groups selectors.
---

# mso_schema_site_anp_epg_selector#

Manages MSO Schema site Application Network Profiles Endpoint Groups Selectors.

## Example Usage ##
```hcl

resource "mso_schema_site_anp_epg_selector" "check" {
  schema_id     = "${mso_schema_site_anp_epg.anp_epg.schema_id}"
  site_id       = "${mso_schema_site_anp_epg.anp_epg.site_id}"
  template      = "${mso_schema_site_anp_epg.anp_epg.template_name}"
  anp_name      = "${mso_schema_site_anp_epg.anp_epg.anp_name}"
  epg_name      = "${mso_schema_site_anp_epg.anp_epg.epg_name}"
  name          = "check01"
  expressions {
    key         = "one"
    operator    = "equals"
    value       = "1"
  }
  expressions {
    key         = "two"
    operator    = "notEquals"
    value       = "22"
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Selector.
* `site_id` - (Required) site ID under which you want to deploy Anp Epg Selector.
* `template` - (Required) Template under above site id where Anp Epg Selector to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name for the selector.
* `expressions` - (Optional) expressions of Selector.
* `expressions.key` - (Required) expression key for the selector.
* `expressions.operator` - (Required) expression operator for the selector. value should be from "equals", "notEquals", "in", "notIn", "keyExist", "keyNotExist".
* `expressions.value` - (Optional) expression value for the selector.

## Attribute Reference ##

No attributes are exported.