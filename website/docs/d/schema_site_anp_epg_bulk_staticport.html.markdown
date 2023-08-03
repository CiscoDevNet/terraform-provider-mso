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

data "mso_schema_site_anp_epg_bulk_staticport" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  anp_name      = "ANP"
  epg_name      = "DB"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Static Port is deployed.
* `site_id` - (Required) The site ID under which the Static Port is deployed.
* `template_name` - (Required) The template name under which the Static Port is deployed.
* `anp_name` - (Required) The ANP name under which the Static Port is deployed.
* `epg_name` - (Required) The EPG name under which the Static Port is deployed.

## Attribute Reference ##

* `static_ports` - (Read-Only) A list of Static Ports.
    * `path_type` - (Read-Only) The path type of the Static Port.
    * `pod` - (Read-Only) The pod of the Static Port.
    * `leaf` - (Read-Only) The leaf of the Static Port.
    * `path` - (Read-Only) The path of the Static Port.
    * `fex` - (Read-Only) The fex-id of the Static Port.
    * `micro_seg_vlan` - (Read-Only) The microsegmentation VLAN id of the Static Port.
    * `mode` - (Read-Only) The mode of the Static Port.
    * `deployment_immediacy` - (Read-Only) The deployment immediacy of the Static Port.
    * `vlan` - (Read-Only) The port encapsulation VLAN id of the Static Port.

 
