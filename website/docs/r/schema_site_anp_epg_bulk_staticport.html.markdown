---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_bulk_staticport"
sidebar_current: "docs-mso-resource-schema_site_anp_epg_bulk_staticport"
description: |-
  Manages MSO Schema Site Application Network Profiles Endpoint Groups Bulk Static Port.
---

# mso_schema_site_anp_epg_bulk_staticport #

Manages MSO Schema Site Application Network Profiles Endpoint Groups Bulk Static Port.

## Example Usage ##

```hcl

resource "mso_schema_site_anp_epg_bulk_staticport" "static_port" {
  schema_id            = mso_schema.schema1.id
  site_id              = mso_schema_site.schema_site.site_id
  template_name        = "Template1"
  anp_name             = mso_schema_site_anp_epg.site_anp_epg.anp_name
  epg_name             = mso_schema_site_anp_epg.site_anp_epg.epg_name
  static_ports {
    path_type            = "vpc"
    deployment_immediacy = "lazy"
    pod                  = "pod-4"
    leaf                 = "105"
    path                 = "eth1/4"
    vlan                 = 207
    mode                 = "regular"
  }
  static_ports {
    path_type            = "port"
    deployment_immediacy = "immediate"
    pod                  = "pod-1"
    leaf                 = "102"
    path                 = "eth1/11"
    vlan                 = 200
    micro_seg_vlan       = 3
    mode                 = "untagged"
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which the Static Port is deployed.
* `site_id` - (Required) SiteID under which the Static Port is deployed.
* `template_name` - (Required) Template name under which the Static Port is deployed.
* `anp_name` - (Required) ANP name under which the Static Port is deployed.
* `epg_name` - (Required) EPG name under which the Static Port is deployed.
* `static_ports` - (Optional) A block representing a Static Port object. Type - Block.
    * `path_type` - (Required) The path type of the static port. Allowed values are `port`, `vpc` and `dpc`. Default to `port`.
    * `pod` - (Required) The pod of the static port.
    * `leaf` - (Required) The leaf of the static port. When `path_type` is `port` or `dpc`, then `leaf` is a string of the leaf ID; Example - '101'. When `path_type` is `vpc`, then `leaf` is a list with both leaf IDs; Example - '101-102'.
    * `path` - (Required) The path of the static port.
    * `mode` - (Required) The mode of the static port. Allowed values are `native`, `regular` and `untagged`.
    * `deployment_immediacy` - (Required) The deployment immediacy of the static port. Allowed values are `immediate` and `lazy`.
    * `vlan` - (Required) The port encapsulation VLAN id of the static port.
    * `micro_seg_vlan` - (Optional) The microsegmentation VLAN id of the static port.
    * `fex` - (Optional) Fex-id to be used. This parameter will work only with the `path_type` as `port`.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Application Network Profiles Endpoint Groups Bulk Static Port can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_anp_epg_bulk_staticport.static_port {schema_id}/site/{site_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}
```
