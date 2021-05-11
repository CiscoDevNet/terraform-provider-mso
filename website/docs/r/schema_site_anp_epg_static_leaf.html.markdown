---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_static_leaf"
sidebar_current: "docs-mso-resource-schema_site_anp_epg_static_leaf"
description: |-
 Manages MSO Schema Site Application Network Profiles Endpoint Groups StaticLeaf.
---

# mso_schema_site_anp_epg_static_leaf #

Manages MSO Schema Site Application Network Profiles Endpoint Groups StaticLeaf.

## Example Usage ##

```hcl
resource "mso_schema_site_anp_epg_static_leaf" "staticleaf1" {
  schema_id = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
  anp_name = "ANP"
  epg_name = "Web"
  path= "topology/pod-1/node-1001"
  port_encap_vlan = 100
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg StaticLeaf.
* `template_name` - (Required) Template where Anp Epg StaticLeaf to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp Epg StaticLeaf.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.
* `path` - (Required) Path Given to the StaticLeaf. ForceNew set to true.
* `port_encap_vlan` - (Required) The VLAN id of the static leaf. ForceNew set to true.


## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Application Network Profiles Endpoint Groups StaticLeaf can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_anp_epg_static_leaf.staticleaf1 {schema_id}/site/{site_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}/path/{static_leaf_path}
```