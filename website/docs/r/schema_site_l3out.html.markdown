---
layout: "mso"
page_title: "MSO: mso_schema_site_l3out"
sidebar_current: "docs-mso-resource-schema_site_l3out"
description: |-
  Manages MSO Schema Site L3out
---

# mso_schema_site_l3out #

Manages MSO Schema Site L3out.

## Example Usage ##

```hcl
resource "mso_schema_site_l3out" "example" {
    vrf_name = data.mso_schema_site_vrf.example.vrf_name
    l3out_name = "example"
    template_name = data.mso_site.example.template_name
    site_id = data.mso_site.example.site_id
    schema_id = data.mso_site.example.schema_id 
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema-id where user wants to add L3out.
* `l3out_name` - (Required) Name of the L3out that user wants to add.
* `template_name` - (Required) Template name associated with the L3out.
* `vrf_name` - (Required) VRF name associated with the L3out.
* `site_id` - (Required) SiteID associated with the L3out.

## Attribute Reference ##
The only Attribute exposed for this resource is `id`. Which is set to the node name of Service Node created.

## Importing ##

An existing MSO Schema Site L3out can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_l3out.example {schema_id}/site/{site_id}/template/{template_name}/vrf/{vrf_name}/l3out/{l3out_name}
```