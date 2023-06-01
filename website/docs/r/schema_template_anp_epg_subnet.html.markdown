---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_subnet"
sidebar_current: "docs-mso-resource-schema_template_anp_epg_subnet"
description: |-
  Manages MSO Schema Template Application Network Profiles Endpoint Groups Subnets.
---

# mso_schema_template_anp_epg_subnet #

Manages MSO Schema Template Application Network Profiles Endpoint Groups Subnets.

## Example Usage ##

```hcl

resource "mso_schema_template_anp_epg_subnet" "subnet1" {
  schema_id = mso_schema.schema1.id
  anp_name  = mso_schema_template_anp_epg.anp_epg.anp_name
  epg_name  = mso_schema_template_anp_epg.anp_epg.name
  template  = "Template1"
  ip        = "31.101.102.0/8"
  scope     = "public"
  shared    = true
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Subnet.
* `template_name` - (Required) Template name under which you want to deploy Subnet.
* `anp_name` - (Required) ANP name under which you want to deploy Subnet.
* `epg_name` - (Required) EPG name under which you want to deploy Subnet.
* `ip` - (Required) The IP range in CIDR notation.
* `description` - (Optional) The description of this subnet.
* `scope` - (Optional) The scope of the subnet. Allowed values are `private` and `public`.
* `shared` - (Optional) Whether this subnet is shared between VRFs.
* `querier` - (Optional) Whether this subnet is an IGMP querier.
* `no_default_gateway` - (Optional) Whether this subnet has a default gateway.
* `primary` - (Optional) Whether this subnet is the primary subnet.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Application Network Profiles Endpoint Groups Subnet can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_anp_epg_subnet.subnet1 {schema_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}/ip/{ip}
```
