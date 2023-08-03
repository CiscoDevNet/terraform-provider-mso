---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_useg_attr"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_useg_attr"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Group uSeg Attribute.
---

# mso_schema_template_anp_epg_useg_attr #

Data source for MSO Schema Template Application Network Profiles Endpoint Group uSeg Attribute.

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg_useg_attr" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "template1"
  anp_name      = "anp1"
  epg_name      = "nkuseg"
  name          = "usg_test"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the uSeg Attribute.
* `template_name` - (Required) The template name of the uSeg Attribute.
* `anp_name` - (Required) The name of the ANP.
* `epg_name` - (Required) The name of the EPG.
* `name` - (Required) The name of the uSeg Attribute.

## Attribute Reference ##

* `useg_type` - (Read-Only) The type of the uSeg Attribute.
* `description` - (Read-Only) The description of the uSeg Attribute.
* `operator` - (Read-Only) The operator of the uSeg Attribute.
* `category` - (Read-Only) The category of the uSeg Attribute.
* `value` - (Read-Only) The value of the uSeg Attribute.
* `useg_subnet` - (Read-Only) Whether the uSeg Subnet is enabled.
