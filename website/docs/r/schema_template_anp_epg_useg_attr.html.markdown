---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_useg_attr"
sidebar_current: "docs-mso-resource-schema_template_anp_epg_useg_attr"
description: |-
   Resource for MSO Schema Template Application Network Profiles Endpoint Groups Useg Attributes.
---

# mso_schema_template_anp_epg_useg_attr #

Resource for MSO Schema Template Application Network Profiles Endpoint Groups Useg Attributes.

## Example Usage ##

```hcl

resource "mso_schema_template_anp_epg_useg_attr" "useg_attrs" {
  schema_id     = mso_schema.schema1.id
  anp_name      = "sanp1"
  epg_name      = mso_schema_template_anp_epg.anp_epg.name
  template_name = "stemplate1"
  name          = "usg_test"
  useg_type     = "tag"
  operator      = "startsWith"
  category      = "tagger"
  value         = "10.2.3.4"
  useg_subnet   = true
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to create Anp Epg Useg Attributes .
* `template_name` - (Required) Template where Anp Epg Useg Attributes to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `name` - (Required) Name of Useg Attributes.
* `useg_type` - (Required) Type of Useg Attribute. Allowed values are `ip`, `mac`, `dns`, `vm-name` (VM Name), `rootContName` (VM Data Center), `hv` (Hypervisor), `guest-os` (VM Operating System), `tag` (VM Tag), `vm` (VM Identifier), `domain` (VMM Domain), `vnic` (Vnic DN).
* `description` - (Optional) String which describes this Useg Attribute.
* `operator` - (Optional) Comparison Operator used in the Useg Attribute. Allowed values are `equals`, `startsWith`, `endsWith`, and `contains`. Default to `equals`. With `useg_type` in [ip, mac, dns] only `equals` operator will be used. Operator passed in the terraform file will be ignored.
* `category` - (Optional) Classifier Category. It's used with useg_type `tag`.
* `value` - (Required) Value of Useg-Attribute.
* `useg_subnet` - (Optional) Whether the Useg Subnet is enabled or not. This field only works with the `useg_type` Ip.

## Attribute Reference ##

The only attribute exported is `id`. Which is set to the name of Useg Attribute.

## Importing ##

An existing MSO Schema Template Application Network Profiles Endpoint Groups Useg Attributes can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_anp_epg_useg_attr.useg_attrs {schema_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}/useg/{name}
```

