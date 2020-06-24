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
* `vrf_name` - (Required) The VRF associated to this External-epg. VRF must exist.
* `vrf_schema_id` - (Optional) SchemaID of VRF. schema_id of External-epg will be used if not provided. Should use this parameter when VRF is in different schema than external-epg.
* `vrf_template_name` - (Optional) Template Name of VRF. template_name of External-epg will be used if not provided.

## Attribute Reference ##

No attributes are exported.
