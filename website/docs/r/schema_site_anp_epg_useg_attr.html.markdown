---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_useg_attr"
sidebar_current: "docs-mso-resource-schema_site_anp_epg_useg_attr"
description: |-
   Resource for MSO Schema Site Application Network Profiles Endpoint Groups Useg Attributes.
---

# mso_schema_site_anp_epg_useg_attr #

Resource for MSO Schema Site Application Network Profiles Endpoint Groups Useg Attributes.

## Example Usage ##

```hcl
resource "mso_schema_site_anp_epg_useg_attr" "example" {
  schema_id     = mso_site.example.schema_id
  site_id       = mso_site.example.site_id
  template_name = mso_site.example.template_name
  anp_name      = mso_schema_site_anp.example.anp_name
  epg_name      = mso_schema_site_anp_epg.example.epg_name
  useg_name     = "useg_site_test"
  useg_type     = "tag"
  operator      = "startsWith"
  category      = "tagger"
  value         = "10.2.3.4"
}
```

## Argument Reference ##

* `schema_id`     - (Required) SchemaID under which you want to create Anp Epg Useg Attributes.
* `site_id`       - (Required) SiteID under which you want to create Anp Epg Useg Attributes.
* `template_name` - (Required) Template where Anp Epg Useg Attributes to be created.
* `anp_name`      - (Required) Name of Application Network Profiles.
* `epg_name`      - (Required) Name of Endpoint Group.
* `useg_name`     - (Required) Name of Useg Attributes.
* `useg_type`     - (Required) Type of Useg Attribute. Allowed values are `ip`, `mac`, `vm-name` (VM Name), `rootContName` (VM Data Center), `hv` (Hypervisor), `guest-os` (VM Operating System), `tag` (VM Tag), `vm` (VM Identifier), `domain` (VMM Domain), `vnic` (Vnic DN).
* `description`   - (Optional) String which describes this Useg Attribute.
* `operator`      - (Optional) Comparison Operator used in the Useg Attribute. Allowed values are `equals`, `startsWith`, `endsWith`, and `contains`. Default to `equals`. With `useg_type` in [ip, mac] only `equals` operator will be used. Operator passed in the terraform file will be ignored.
* `category`      - (Optional) Classifier Category. It's used with useg_type `tag`.
* `value`         - (Required) Value of Useg-Attribute.
* `fv_subnet`     - (Optional) Whether the Fv Subnet is enabled or not. This field only works with the `useg_type` Ip.

## Attribute Reference ##

The only attribute exported is `id`. Which is set to the name of Useg Attribute.

## Importing ##

An existing MSO Schema Site Application Network Profiles Endpoint Groups Useg Attributes can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_anp_epg_useg_attr.example {schema_id}/site/{site_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}/uSegAttr/{useg_name}
```

