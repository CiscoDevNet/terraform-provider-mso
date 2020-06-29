---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_selector"
sidebar_current: "docs-mso-resource-schema_template_anp_epg_selector"
description: |-
  Manages MSO Schema Template Application Network Profiles Endpoint Groups selectors.
---

# mso_schema_template_anp_epg_selector#

Manages MSO Schema Template Application Network Profiles Endpoint Groups Selectors.

## Example Usage ##
```hcl
resource "mso_schema_template_anp_epg_selector" "check" {
  schema_id     = "${mso_schema_template_anp_epg.anp_epg.schema_id}"
  template      = "${mso_schema_template_anp_epg.anp_epg.template_name}"
  anp_name      = "${mso_schema_template_anp_epg.anp_epg.anp_name}"
  epg_name      = "${mso_schema_template_anp_epg.anp_epg.name}"
  name          = "check01"
  expressions {
    key         = "one"
    operator    = "equals"
    value       = "1"
  }
  expressions {
    key         = "two"
    operator    = "equals"
    value       = "2"
  }
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Subnet.
* `template` - (Required) Template where Anp Epg Subnet to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name for the selector.
* `expressions` - (Optional) expressions of Selector.
* `expressions.key` - (Optional) expression key for the selector.
* `expressions.operator` - (Optional) expression operator for the selector. value should be from 
* `expressions.value` - (Optional) expression value for the selector.

## Attribute Reference ##

No attributes are exported.