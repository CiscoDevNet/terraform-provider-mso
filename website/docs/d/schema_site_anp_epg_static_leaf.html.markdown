---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_static_leaf"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_static_leaf"
description: |-
  Data source for MSO Schema Site Application Network Profiles Endpoint Groups StaticLeaf.
---

# mso_schema_site_anp_epg_static_leaf #

Data source for MSO Schema Site Application Network Profiles Endpoint Groups StaticLeaf.

## Example Usage ##

```hcl

data "mso_schema_site_anp_epg_static_leaf" "st10" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  site_id       = data.mso_site.site1.id
  anp_name      = "ANP"
  epg_name      = "Web"
  path          = "topology/pod-1/paths-103/pathep-[eth1/111]"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg StaticLeaf.
* `site_id` - (Required) SiteID under which you want to deploy Anp Epg StaticLeaf.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.
* `path` - (Required) Path Given to the StaticLeaf.


## Attribute Reference ##

* `template_name` - (Optional) Template where Anp Epg StaticLeaf to be created.
* `port_encap_vlan` - (Optional) The VLAN id of the static leaf.


