---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_subnet"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_subnet"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Groups Subnet.
---

# mso_schema_template_anp_epg_subnet #

Data source for MSO Schema Template Application Network Profiles Endpoint Groups Subnet.

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg_subnet" "subnet1" {
  schema_id = data.mso_schema.schema1.id
  anp_name  = "WoS-Cloud-Only-2"
  epg_name  = "EPG4"
  template  = "Template1"
  ip        = "31.101.102.0/8"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Subnet.
* `template` - (Required) Template name under which you want to deploy Subnet.
* `anp_name` - (Required) ANP name under which you want to deploy Subnet.
* `epg_name` - (Required) EPG name under which you want to deploy Subnet.
* `ip` - (Required) The IP range in CIDR notation.

## Attribute Reference ##

* `description` - (Read-Only) The description of this subnet.
* `scope` - (Read-Only) The scope of the subnet. Allowed values are `private` and `public`.
* `shared` - (Read-Only) Whether this subnet is shared between VRFs.
* `querier` - (Read-Only) Whether this subnet is an IGMP querier.
* `no_default_gateway` - (Read-Only) Whether this subnet has a default gateway.
* `primary` - (Read-Only) Whether this subnet is the primary subnet.
