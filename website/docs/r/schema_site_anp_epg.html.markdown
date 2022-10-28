---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg"
sidebar_current: "docs-mso-resource-schema_site_anp_epg"
description: |-
  Manages MSO Schema Site Application Network Profiles Endpoint Groups.
---

# mso_schema_site_anp_epg #

Manages MSO Schema Site Application Network Profiles Endpoint Groups.

## Example Usage ##

```hcl

resource "mso_schema_site_anp_epg" "site_anp_epg" {
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
  site_id       = mso_schema_site.schema_site.site_id
  anp_name      = mso_schema_site_anp.anp1.anp_name
  epg_name      = mso_schema_template_anp_epg.anp_epg.name
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg.
* `template_name` - (Required) Template where Anp Epg to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp Epg.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.
* `private_link_label` - (Optional) Private Link Label.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Application Network Profiles Endpoint Group can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_anp_epg.site_anp_epg {schema_id}/site/{site_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}
```