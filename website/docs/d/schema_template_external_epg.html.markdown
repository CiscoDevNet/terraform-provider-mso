---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg"
sidebar_current: "docs-mso-data-source-schema_template_external_epg"
description: |-
  MSO Schema Template External Endpoint Group Data source.
---

# mso_schema_template_external_epg #

MSO Schema Template External Endpoint Group Data source.

## Example Usage ##

```hcl

data "mso_schema_template_external_epg" "externalEpg" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  external_epg_name = "ExternalEPG1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy External-epg.
* `template_name` - (Required) Template where External-epg to be created.
* `external_epg_name` - (Required) Name of External-epg.

## Attribute Reference ##

* `display_name` - (Optional) Display Name of the External-epg on the MSO UI.
* `vrf_name` - (Optional) The VRF associated to this External-epg. VRF must exist.
* `vrf_schema_id` - (Optional) SchemaID of VRF. schema_id of External-epg will be used if not provided. Should use this parameter when VRF is in different schema than external-epg.
* `vrf_template_name` - (Optional) Template Name of VRF. template_name of External-epg will be used if not provided.
* `anp_name` - (Optional) Name of anp to attach.
* `anp_schema_id` - (Optional) SchemaId of anp. `schema_id` will be used if not provided.
* `anp_template_name` - (Optional) Template name of anp. `template_name` will be used if not provided.

* `site_id` - (Optional) List of ids of sites associated with the schema. Required when `external_epg_type` is "cloud".
* `selector_name` - (Optional) name of the selector for external epg. Required when `external_epg_type` is "cloud".
* `selector_ip` - (Optional) ip address for expression in selector. Required when `external_epg_type` is "cloud".
