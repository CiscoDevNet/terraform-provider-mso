---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_useg_attr"
sidebar_current: "docs-mso-data-source-mso_schema_site_anp_epg_useg_attr"
description: |-
  Data source for MSO Schema Site Application Network Profiles Endpoint Groups Useg Attributes.
---

# mso_schema_site_anp_epg_useg_attr #

Data source for MSO Schema Site Application Network Profiles Endpoint Groups Useg Attributes.

## Example Usage ##

```hcl
data "mso_schema_site_anp_epg_useg_attr" "useg_attrs" {
  schema_id     = mso_site.example.schema_id
  site_id       = mso_site.example.site_id
  template_name = mso_site.example.template_name
  anp_name      = mso_schema_site_anp.example.anp_name
  epg_name      = mso_schema_site_anp_epg.example.epg_name
  useg_name     = "useg_site_test"
}
```

## Argument Reference ##

* `schema_id`     - (Required) SchemaID under which you want to create Anp Epg Useg Attributes.
* `site_id`       - (Required) SiteID under which you want to create Anp Epg Useg Attributes.
* `template_name` - (Required) Template where Anp Epg Useg Attributes to be created.
* `anp_name`      - (Required) Name of Application Network Profiles.
* `epg_name`      - (Required) Name of Endpoint Group.
* `useg_name`     - (Required) Name of Useg Attributes.

## Attribute Reference ##

* `useg_type`   - (Optional) Type of Useg Attribute.
* `description` - (Optional) String which describes this Useg Attribute.
* `operator`    - (Optional) Comparison Operator used in the Useg Attribute.
* `category`    - (Optional) Classifier Category. It's used with useg_type `tag`.
* `value`       - (Optional) Value of Useg-Attribute.
* `fv_subnet`   - (Optional) Whether the Useg Subnet is enabled or not. This field only works with the `useg_type` Ip.