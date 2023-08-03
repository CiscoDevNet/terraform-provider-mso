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

* `schema_id` - (Required) The schema ID under which the Static Port is deployed.
* `site_id` - (Required) The site ID under which the Static Port is deployed.
* `template_name` - (Required) The template name under which the Static Port is deployed.
* `anp_name` - (Required) The ANP name under which the Static Port is deployed.
* `epg_name` - (Required) The EPG name under which the Static Port is deployed.
* `path_type` - (Required) The type of the Static Port.
* `pod` - (Required) The pod of the Static Port.
* `leaf` - (Required) The leaf of the Static Port.
* `path` - (Required) The path of the Static Port.
* `fex` - (Optional) The fex ID of the Static Port. This parameter will work only with the `path_type` as `port`.


## Attribute Reference ##

* `micro_seg_vlan` - (Read-Only) The microsegmentation VLAN ID of the Static Port.
* `mode` - (Read-Only) The mode of the Static Port.
* `deployment_immediacy` - (Read-Only) The deployment immediacy of the Static Port.
* `vlan` - (Read-Only) The VLAN ID of the Static Port.

 
