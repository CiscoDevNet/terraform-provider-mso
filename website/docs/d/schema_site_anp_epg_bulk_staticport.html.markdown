---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_bulk_staticport"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_bulk_staticport"
description: |-
  Data source for MSO Schema Site ANP EPG Bulk Static Port.
---

# mso_schema_site_anp_epg_bulk_staticport #

Data source for MSO Schema Site ANP EPG Bulk Static Port.

## Example Usage ##

```hcl

data "mso_schema_site_anp_epg_bulk_staticport" "static_port" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  anp_name      = "ANP"
  epg_name      = "DB"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Static Port.
* `site_id` - (Required) SiteID under which you want to deploy Static Port.
* `template_name` - (Required) Template name under which Static Port is deployed.
* `anp_name` - (Required) ANP name under which you want to deploy Static Port.
* `epg_name` - (Required) EPG name under which you want to deploy Static Port.


## Attribute Reference ##

* `static_ports` - (Optional) A block representing a Static Port object. Type: Block.
    * `path_type` - (Optional) The type of the static port.
    * `pod` - (Optional) The pod of the static port.
    * `leaf` - (Optional) The leaf of the static port.
    * `path` - (Optional) The path of the static port.
    * `fex` - (Optional) Fex-id to be used. This parameter will work only with the `path_type` as `port`.
    * `micro_seg_vlan` - (Optional) The microsegmentation VLAN id of the static port.
    * `mode` - (Optional) The mode of the static port.
    * `deployment_immediacy` - (Optional) The deployment immediacy of the static port.
    * `vlan` - (Optional) The port encap VLAN id of the static port.

 
