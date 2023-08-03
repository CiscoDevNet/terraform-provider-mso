---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_subnet"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_subnet"
description: |-
  Data source for MSO Schema Site ANP EPG Subnet.
---

# mso_schema_site_anp_epg_subnet #

Data source for MSO Schema Site ANP EPG Subnet.

## Example Usage ##

```hcl

data "mso_schema_site_anp_epg_subnet" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  anp_name      = "ANP"
  epg_name      = "DB"
  ip            = "10.7.0.1/8"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Subnet is deployed.
* `site_id` - (Required) The site ID under which the Subnet is deployed.
* `template_name` - (Required) The template name under which the Subnet is deployed.
* `anp_name` - (Required) The ANP name under which the Subnet is deployed.
* `epg_name` - (Required) The EPG name under which the Subnet is deployed.
* `ip` - (Required) The IP range in CIDR notation of the Subnet.

## Attribute Reference ##

* `description` - (Read-Only) The description of the Subnet.
* `scope` - (Read-Only) The scope of the Subnet.
* `shared` - (Read-Only) Whether the Subnet is shared between VRFs.
* `querier` - (Read-Only) Whether the Subnet is an IGMP querier.
* `no_default_gateway` - (Read-Only) Whether the Subnet has a default gateway.
* `primary` - (Read-Only) Whether the Subnet is the primary Subnet.
