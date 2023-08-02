---
layout: "mso"
page_title: "MSO: mso_schema_template_l3out"
sidebar_current: "docs-mso-data-source-schema_template_l3out"
description: |-
  Data source for MSO Schema Template L3Out.
---

# mso_schema_template_l3out #

Data source for MSO Schema Template L3Out.

## Example Usage ##

```hcl

data "mso_schema_template_l3out" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  l3out_name    = "Internet_L3Out"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the L3out.
* `template_name` - (Required) The template name of the L3out.
* `l3out_name` - (Required) The name of the L3Out.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the L3out as displayed on the MSO UI.
* `vrf_name` - (Read-Only) The name of the VRF associated with the L3out.
* `vrf_schema_id` - (Read-Only) The schema ID of the VRF associated with the L3out.
* `vrf_template_name` - (Read-Only) The template name of the VRF associated with the L3out.
