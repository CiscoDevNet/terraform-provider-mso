---
layout: "mso"
page_title: "MSO: mso_schema_site_bd_subnet"
sidebar_current: "docs-mso-data-source-schema_site_bd_subnet"
description: |-
  Data source for MSO Schema Site Bd Subnet.
---

# mso_schema_site_bd_subnet #

Data source for MSO Schema Site Bridge Domain(Bd) Subnet.

## Example Usage ##

```hcl

data "mso_schema_site_bd_subnet" "s1" {
  schema_id = data.mso_schema.schema1.id
  site_id   = data.mso_site.site1.id
  bd_name   = "WebServer-Finance"
  ip        = "200.168.240.1/24"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy the subnet.
* `site_id` - (Required) SiteID under which you want to deploy the subnet.
* `bd_name` - (Required) Bd name under which you want to deploy the subnet.
* `ip` - (Required) The IP of the subnet.
* `template_name` - (Required) Template name under which you want to deploy the subnet.

## Attribute Reference ##

* `scope` - (Read-Only) The scope of the subnet.
* `shared` - (Read-Only) Whether this subnet is shared between VRFs.
* `querier` - (Read-Only) Whether this subnet is an IGMP querier.
* `no_default_gateway` - (Read-Only) Whether this subnet has a default gateway.
* `description` - (Read-Only) The description of this subnet. 

