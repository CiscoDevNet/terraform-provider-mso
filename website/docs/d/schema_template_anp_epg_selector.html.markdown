---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_selector"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_selector"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Groups Selector.
---

# mso_schema_template_anp_epg_selector #

Data source for MSO Schema Template Application Network Profiles Endpoint Groups Selector.

```hcl

data "mso_schema_template_anp_epg_selector" "read_check" {
  schema_id       = data.mso_schema.schema1.id
  template_name   = "Template1"
  anp_name        = "anp1"
  epg_name        = "epg1
  name            = "subnet1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Subnet.
* `template_name` - (Required) Template where Anp Epg Subnet to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name of Subnet.

## Attribute Reference ##

* `expressions` - (Optional) expressions of Selector.
* `expressions.key` - (Optional) expression key for the selector.
* `expressions.operator` - (Optional) expression operator for the selector. value should be from "equals", "notEquals", "in", "notIn", "keyExist", "keyNotExist".
* `expressions.value` - (Optional) expression value for the selector.