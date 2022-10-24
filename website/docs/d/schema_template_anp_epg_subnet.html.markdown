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

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Subnet.
* `template_name` - (Required) Template where Anp Epg Subnet to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `ip` - (Required) Ip Address of Subnet.

## Attribute Reference ##

* `scope` - (Optional) Scope of Subnet.
* `shared` - (Optional) Whether the subnet should be shared or not.