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
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  external_epg_name = "UntitledExternalEPG1"
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
