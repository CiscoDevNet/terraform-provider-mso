---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_subnet"
sidebar_current: "docs-mso-data-source-schema_template_external_epg_subnet"
description: |-
  Data source for MSO Schema Template External EPG Subnet.
---

# mso_schema_template_external_epg_subnet #

Data source for MSO Schema Template External EPG Subnet.

## Example Usage ##

```hcl

data "mso_schema_template_external_epg_subnet" "example" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  external_epg_name = "Internet"
  ip                = "30.1.1.0/24"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the External EPG.
* `template_name` - (Required) The template name of the External EPG.
* `external_epg_name` - (Required) The name of the External EPG.
* `ip` - (Required) The IP range of the External EPG in CIDR notation.

## Attribute Reference ##

* `name` - (Read-Only) The name of Subnet.
* `scope` - (Read-Only) The scope of the subnet.
* `aggregate` - (Read-Only) The aggregate of the subnet.
