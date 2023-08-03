---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_subnet"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_subnet"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Group Subnet.
---

# mso_schema_template_anp_epg_subnet #

Data source for MSO Schema Template Application Network Profiles Endpoint Group Subnet.

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg_subnet" "example" {
  schema_id = data.mso_schema.schema1.id
  template  = "Template1"
  anp_name  = "WoS-Cloud-Only-2"
  epg_name  = "EPG4"
  ip        = "31.101.102.0/8"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Subnet.
* `template` - (Required) The template name of the Subnet.
* `anp_name` - (Required) The name of the ANP.
* `epg_name` - (Required) The name of the EPG.
* `ip` - (Required) The IP range in CIDR notation.

## Attribute Reference ##

* `description` - (Read-Only) The description of the Subnet.
* `scope` - (Read-Only) The scope of the Subnet.
* `shared` - (Read-Only) Whether the Subnet is shared between VRFs.
* `querier` - (Read-Only) Whether the Subnet is an IGMP querier.
* `no_default_gateway` - (Read-Only) Whether the Subnet has a default gateway.
* `primary` - (Read-Only) Whether the Subnet is the primary Subnet.
