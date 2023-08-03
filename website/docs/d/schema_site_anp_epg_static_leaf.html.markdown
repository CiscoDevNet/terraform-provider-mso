---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_static_leaf"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_static_leaf"
description: |-
  Data source for MSO Schema Site Application Network Profiles End Point Group Static Leaf.
---

# mso_schema_site_anp_epg_static_leaf #

Data source for MSO Schema Site Application Network Profiles End Point Group Static Leaf.

## Example Usage ##

```hcl

data "mso_schema_site_anp_epg_static_leaf" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  site_id       = data.mso_site.site1.id
  anp_name      = "ANP"
  epg_name      = "Web"
  path          = "topology/pod-1/paths-103/pathep-[eth1/111]"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Static Leaf is deployed.
* `site_id` - (Required) The site ID under which the Static Leaf is deployed.
* `template_name` - (Required) The template name under which the Static Leaf is deployed.
* `anp_name` - (Required) The ANP name under which the Static Leaf is deployed.
* `epg_name` - (Required) The EPG name under which the Static Leaf is deployed.
* `path` - (Required) The Path of the Static Leaf.

## Attribute Reference ##

* `port_encap_vlan` - (Read-Only) The port encapsulation VLAN ID of the Static Leaf.
