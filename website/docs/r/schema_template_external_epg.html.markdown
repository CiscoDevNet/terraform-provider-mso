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
  schema_id         = mso_schema.schema1.id
  template_name     = "Template1"
  external_epg_name = "temp_epg"
  external_epg_type = "cloud"
  display_name      = "temp_epg"
  vrf_name          = mso_schema_template_vrf.vrf1.name
  anp_name          = mso_schema_template_anp.anp1.name
  l3out_name        = mso_schema_template_l3out.template_l3out.l3out_name
  site_id           = [mso_schema_site.schema_site.site_id]
  selector_name     = "check02"
  selector_ip       = "12.23.34.45"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy External EPG.
* `template_name` - (Required) Template where External EPG is to be created.
* `external_epg_name` - (Required) Name of the External EPG.
* `display_name` - (Required) Display Name of the External EPG on the MSO UI.
* `description` - (Optional) The description of the External EPG.
* `external_epg_type` - (Optional) The type of External EPG. Allowed values are `on-premise` and `cloud`. Default to `on-premise`.
* `vrf_name` - (Required) The VRF associated with the External EPG. VRF must exist.
* `vrf_schema_id` - (Optional) SchemaID of VRF. schema_id of External EPG will be used if not provided. This parameter should be used when VRF is in a different schema than external EPG.
* `vrf_template_name` - (Optional) Template Name of VRF. The template_name of External EPG will be used if not provided.
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

* `uuid` - The UUID of the External EPG.

## Importing ##

An existing MSO Schema Template External Endpoint Group can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_external_epg.template_externalepg {schema_id}/template/{template_name}/externalEPG/{external_epg_name}
```