---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg"
sidebar_current: "docs-mso-data-source-schema_template_external_epg"
description: |-
  Data source for MSO Schema Template External End Point Group.
---

# mso_schema_template_external_epg #

Data source for MSO Schema Template External End Point Group.

## Example Usage ##

```hcl

data "mso_schema_template_external_epg" "example" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  external_epg_name = "ExternalEPG1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the External EPG.
* `template_name` - (Required) The template name of the External EPG.
* `external_epg_name` - (Required) The name of the External EPG.

## Attribute Reference ##

* `uuid` - (Read-Only) The UUID of the External EPG.
* `display_name` - (Read-Only) The name of the External EPG as displayed on the MSO UI.
* `external_epg_type` - (Read-Only) The type of the External EPG.
* `vrf_name` - (Read-Only) The name of the VRF associated with the External EPG.
* `vrf_schema_id` - (Read-Only) The schema ID of the VRF associated with the External EPG.
* `vrf_template_name` - (Read-Only) The template name of the VRF associated with the External EPG.
* `anp_name` - (Read-Only) The name of the ANP associated with the External EPG.
* `anp_schema_id` - (Read-Only) The schema ID of the ANP associated with the External EPG.
* `anp_template_name` - (Read-Only) The template name of the ANP associated with the External EPG.
* `l3out_name` - (Read-Only) The name of the L3out associated with the External EPG.
* `l3out_schema_id` - (Read-Only) The schema ID of the L3out associated with the External EPG.
* `l3out_template_name` - (Read-Only) The template name of the L3out associated with the External EPG.
* `selector_name` - (Read-Only) The name of the External EPG selector.
* `selector_ip` - (Read-Only) The ip address of the External EPG selector.
* `description` - (Read-Only) The description of the External EPG selector.
