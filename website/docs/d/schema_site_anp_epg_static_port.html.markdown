---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_static_port"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_static_port"
description: |-
  Data source for MSO Schema Site ANP EPG Static Port.
---

# mso_schema_site_anp_epg_static_port #

Data source for MSO Schema Site ANP EPG Static Port.

## Example Usage ##

```hcl

data "mso_schema_site_anp_epg_static_port" "static_port" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  anp_name      = "ANP"
  epg_name      = "DB"
  path_type     = "port"
  pod           = "pod-7"
  leaf          = "109"
  path          = "eth1/10"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Static Port.
* `site_id` - (Required) SiteID under which you want to deploy Static Port.
* `template_name` - (Required) Template name under which Static Port is deployed.
* `anp_name` - (Required) ANP name under which you want to deploy Static Port.
* `epg_name` - (Required) EPG name under which you want to deploy Static Port.
* `path_type` - (Required) The type of the static port.
* `pod` - (Required) The pod of the static port.
* `leaf` - (Required) The leaf of the static port.
* `path` - (Required) The path of the static port.
* `fex` - (Optional) Fex-id to be used. This parameter will work only with the `path_type` as `port`.


## Attribute Reference ##

* `micro_seg_vlan` - (Optional) The microsegmentation VLAN id of the static port.
* `mode` - (Optional) The mode of the static port.
* `deployment_immediacy` - (Optional) The deployment immediacy of the static port.
* `vlan` - (Optional) The port encap VLAN id of the static port.

 
