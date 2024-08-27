---
layout: "mso"
page_title: "MSO: mso_schema_site_external_epg"
sidebar_current: "docs-mso-resource-schema_site_external_epg"
description: |-
  Manages MSO Schema Site External Endpoint Groups.
---

# mso_schema_site_external_epg_selector#

Manages MSO Schema Site External Endpoint Groups.

## Example Usage ##
```hcl

resource "mso_schema_site_external_epg" "external_epg_1" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  site_id           = data.mso_site.site1.id
  external_epg_name = "external_epg_1"
}

```

## Argument Reference ##

* `schema_id` - (Required) Schema ID under which you want to deploy the External Endpoint Group.
* `site_id` - (Required) Site ID under which you want to deploy the External Endpoint Group.
* `template_name` - (Required) Template Name under which you want to define the External Endpoint Group.
* `external_epg_name` - (Required) Name of the External Endpoint Group.

## Attribute Reference ##

* `l3out_name` - (Optional) Name of the L3Out.
* `l3out_schema_id` - (Optional) Schema ID of the L3Out.
* `l3out_template_name` - (Optional) Template of the L3Out.
* `l3out_on_apic` - (Optional) Whether L3Out is created only locally on APIC.

## Importing ##

An existing MSO Schema Site External Endpoint Group can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_external_epg.extepg1 {schema_id}/site/{site_id}/externalEPG/{external_epg_name}
```