---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_static_port"
sidebar_current: "docs-mso-resource-schema_template_anp_epg_static_port"
description: |-
  Manages MSO Schema Template Application Network Profiles Endpoint Groups Static Port.
---

# mso_schema_template_anp_epg_static_port #

Manages MSO Schema Template Application Network Profiles Endpoint Groups Static Port.

## Example Usage ##

```hcl
resource "mso_schema_site_anp_epg_static_port" "static_port" {
  schema_id = "5c4d5bb72700000401f80948"
  site_id = "5c7c95b25100008f01c1ee3c"
  template_name = "Template1"
  anp_name = "ANP"
  epg_name = "DB"
  path_type = "port"
  deployment_immediacy = "lazy"
  pod = "pod-7"
  leaf = "109"
  path = "eth1/10"
  vlan = 200
  mode = "untagged"

}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Static Port.
* `site_id` - (Required) SiteID under which you want to deploy Static Port.
* `template_name` - (Required) Template name under which you want to deploy Static Port.
* `anp_name` - (Required) ANP name under which you want to deploy Static Port.
* `epg_name` - (Required) EPG name under which you want to deploy Static Port.
* `path_type` - (Required) The type of the static port. Allowed value is `port`.
* `pod` - (Required) The pod of the static port.
* `leaf` - (Required) The leaf of the static port.
* `path` - (Required) The path of the static port.
* `mode` - (Required) The mode of the static port. Allowed values are `native`, `regular` and `untagged`.
* `deployment_immediacy` - (Required) The deployment immediacy of the static port. Allowed values are `immediate` and `lazy`.
* `vlan` - (Required) The port encap VLAN id of the static port.
* `micro_segvlan` - (Optional) The microsegmentation VLAN id of the static port.

## Attribute Reference ##

No attributes are exported.
