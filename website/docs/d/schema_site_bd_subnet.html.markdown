---
layout: "mso"
page_title: "MSO: mso_schema_site_bd_subnet"
sidebar_current: "docs-mso-data-source-schema_site_bd_subnet"
description: |-
  Data source for MSO Schema Site Bridge Domain (BD) Subnet.
---

# mso_schema_site_bd_subnet #

Data source for MSO Schema Site Bridge Domain (BD) Subnet.

## Example Usage ##

```hcl

data "mso_schema_site_bd_subnet" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  bd_name       = "WebServer-Finance"
  ip            = "200.168.240.1/24"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Subnet is deployed.
* `site_id` - (Required) The site ID under which the Subnet is deployed.
* `template_name` - (Required) The template name under which the Subnet is deployed.
* `bd_name` - (Required)  The bridge domain name under which the Subnet is deployed.
* `ip` - (Required) The IP of the Subnet.

## Attribute Reference ##

* `scope` - (Read-Only) The scope of the Subnet.
* `shared` - (Read-Only) Whether the Subnet is shared between VRFs.
* `querier` - (Read-Only) Whether the Subnet is an IGMP querier.
* `no_default_gateway` - (Read-Only) Whether the Subnet has a default gateway.
* `description` - (Read-Only) The description of the Subnet.
* `primary` - (Read-Only) Whether the Subnet is the primary Subnet.
* `virtual` - (Read-Only) Whether the Subnet is virtual.
