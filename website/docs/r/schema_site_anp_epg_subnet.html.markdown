---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_subnet"
sidebar_current: "docs-mso-resource-schema_site_anp_epg_subnet"
description: |-
  Manages MSO Schema Site Application Network Profiles Endpoint Groups Subnet.
---

# mso_schema_site_anp_epg_subnet #

Manages MSO Schema Site Application Network Profiles Endpoint Groups Subnet.

## Example Usage ##

```hcl
resource "mso_schema_site_anp_epg_subnet" "static_port" {
  schema_id = "5c4d5bb72700000401f80948"
  site_id = "5c7c95b25100008f01c1ee3c"
  template_name = "Template1"
  anp_name = "ANP"
  epg_name = "DB"
  ip = "10.7.0.1/8"
  scope = "public"
  shared = true

}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Subnet.
* `site_id` - (Required) SiteID under which you want to deploy Subnet.
* `template_name` - (Required) Template name under which you want to deploy Subnet.
* `anp_name` - (Required) ANP name under which you want to deploy Subnet.
* `epg_name` - (Required) EPG name under which you want to deploy Subnet.
* `ip` - (Required) The IP range in CIDR notation.
* `description` - (Optional) The description of this subnet.
* `scope` - (Required) The scope of the subnet. Allowed values are `private` and `public`.
* `shared` - (Required) Whether this subnet is shared between VRFs.
* `querier` - (Optional) Whether this subnet is an IGMP querier.
* `no_default_gateway` - (Optional) Whether this subnet has a default gateway.


## Attribute Reference ##

No attributes are exported.
