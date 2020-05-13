---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_staticleaf"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_staticleaf"
description: |-
  Data source for MSO Schema Site Application Network Profiles Endpoint Groups StaticLeaf.
---

# mso_schema_site_anp_epg_staticleaf #

Data source for MSO Schema Site Application Network Profiles Endpoint Groups StaticLeaf.

## Example Usage ##

```hcl
data "mso_schema_site_anp_epg_staticleaf" "st10" {
  schema_id = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
  anp_name = "ANP"
  epg_name = "Web"
  path= "topology/pod-1/paths-103/pathep-[eth1/111]"
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


