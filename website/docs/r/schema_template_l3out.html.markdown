---
layout: "mso"
page_title: "MSO: mso_schema_template_l3out"
sidebar_current: "docs-mso-resource-schema_template_l3out"
description: |-
  Manages MSO Schema Template L3Out.
---

# mso_schema_template_l3out #

Manages MSO Schema Template L3Out.

## Example Usage ##

```hcl
resource "mso_schema_template_l3out" "template_l3out" {
  schema_id = "5c6c16d7270000c710f8094d"
  template_name = "Template1"
  l3out_name = "l3out1"
  display_name = "l3out1"
  vrf_name = "vrf2"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy L3Out.
* `template_name` - (Required) Template where L3Out to be created.
* `l3out_name` - (Required) Name of L3Out.
* `display_name` - (Required) Display Name of the L3Out on the MSO UI.
* `vrf_name` - (Required) The VRF associated to this L3out. VRF must exist.
* `vrf_schema_id` - (Optional) SchemaID of VRF. schema_id of L3Out will be used if not provided. Should use this parameter when VRF is in different schema than l3out.
* `vrf_template_name` - (Optional) Template Name of VRF. template_name of L3Out will be used if not provided.

## Attribute Reference ##

No attributes are exported.
