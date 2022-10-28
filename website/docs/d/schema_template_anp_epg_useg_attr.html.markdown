---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_useg_attr"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_useg_attr"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Groups Useg Attributes.
---

# mso_schema_template_anp_epg_useg_attr #

Data source for MSO Schema Template Application Network Profiles Endpoint Groups Useg Attributes.

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg_useg_attr" "useg_attrs" {
  schema_id     = data.mso_schema.schema1.id
  anp_name      = "sanp1"
  epg_name      = "nkuseg"
  template_name = "stemplate1"
  name          = "usg_test"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to create Anp Epg Useg Attributes .
* `template_name` - (Required) Template where Anp Epg Useg Attributes to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name of Useg Attributes.

## Attribute Reference ##

* `useg_type` - (Optional) Type of Useg Attribute.
* `description` - (Optional) String which describes this Useg Attribute.
* `operator` - (Optional) Comparison Operator used in the Useg Attribute.
* `category` - (Optional) Classifier Category. It's used with useg_type `tag`.
* `value` - (Optional) Value of Useg-Attribute.
* `useg_subnet` - (Optional) Whether the Useg Subnet is enabled or not. This field only works with the `useg_type` Ip.