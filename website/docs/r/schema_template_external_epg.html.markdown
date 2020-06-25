---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg"
sidebar_current: "docs-mso-resource-schema_template_external_epg"
description: |-
  Manages MSO Schema Template External Endpoint Group.
---

# mso_schema_template_external_epg #

Manages MSO Schema Template External Endpoint Group.

## Example Usage ##

```hcl
resource "mso_schema_template_external_epg" "template_externalepg" {
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  external_epg_name = "external_epg12"
  display_name = "external_epg12"
  vrf_name = "vrf1"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy External-epg.
* `template_name` - (Required) Template where External-epg to be created.
* `external_epg_name` - (Required) Name of External-epg.
* `display_name` - (Required) Display Name of the External-epg on the MSO UI.
* `external_epg_type` - (Optional) Type of External EPG. Allowed values are `on-premise` and `cloud`. Default to `on-premise`.
* `vrf_name` - (Required) The VRF associated to this External-epg. VRF must exist.
* `vrf_schema_id` - (Optional) SchemaID of VRF. schema_id of External-epg will be used if not provided. Should use this parameter when VRF is in different schema than external-epg.
* `vrf_template_name` - (Optional) Template Name of VRF. template_name of External-epg will be used if not provided.
* `include_in_preferred_group` - (Optional) This parameter indicates whether EPG is included in preferred group or not. Default to false.
* `l3out_name` - (Optional) Name of L3out to attach. Should use this parameter with `external_epg_type` as `on-premise`.
* `l3out_schema_id` - (Optional) SchemaId of L3out. `schema_id` will be used if not provided. Should use this parameter with `external_epg_type` as `on-premise`.
* `l3out_template_name` - (Optional) Template name of L3out. `template_name` will be used if not provided. Should use this parameter with `external_epg_type` as `on-premise`.

NOTE: SchemaID and Template Name for VRF and L3out must be same.


## Attribute Reference ##

No attributes are exported.
