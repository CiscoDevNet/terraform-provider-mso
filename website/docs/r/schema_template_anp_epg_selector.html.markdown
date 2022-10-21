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
  schema_id     = mso_schema_template_anp_epg.anp_epg.schema_id
  template_name = mso_schema_template_anp_epg.anp_epg.template_name
  anp_name      = mso_schema_template_anp_epg.anp_epg.anp_name
  epg_name      = mso_schema_template_anp_epg.anp_epg.name
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
* `template_name` - (Required) Template where Anp Epg Subnet to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name for the selector.
* `expressions` - (Optional) expressions of Selector.
* `expressions.key` - (Required) expression key for the selector.
* `expressions.operator` - (Required) expression operator for the selector. value should be from "equals", "notEquals", "in", "notIn", "keyExist", "keyNotExist".
* `expressions.value` - (Optional) expression value for the selector.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Application Network Profiles Endpoint Groups Selector can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_anp_epg_selector.check {schema_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}/selector/{name}
```