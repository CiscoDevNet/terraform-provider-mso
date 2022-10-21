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
		schema_id           = mso_schema.schema1.id
		template_name       = "Template1"
		external_epg_name   = "temp_epg"
    external_epg_type   = "cloud"
		display_name        = "temp_epg"
		vrf_name            = "Myvrf"
    anp_name            = "ap1"
    l3out_name          = "temp"
    site_id             = ["5c7c95d9510000cf01c1ee3d"]
    selector_name       = "check02"
    selector_ip         = "12.23.34.45"
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
* `anp_name` - (Optional) Name of anp to attach.
* `anp_schema_id` - (Optional) SchemaId of anp. `schema_id` will be used if not provided.
* `anp_template_name` - (Optional) Template name of anp. `template_name` will be used if not provided.

* `site_id` - (Optional) List of ids of sites associated with the schema. Required when `external_epg_type` is "cloud".
* `selector_name` - (Optional) name of the selector for external epg. Required when `external_epg_type` is "cloud".
* `selector_ip` - (Optional) ip address for expression in selector. Required when `external_epg_type` is "cloud".

NOTE: SchemaID and Template Name for VRF and L3out must be same.


## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template External Endpoint Group can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_external_epg.template_externalepg {schema_id}/template/{template_name}/externalEPG/{external_epg_name}
```