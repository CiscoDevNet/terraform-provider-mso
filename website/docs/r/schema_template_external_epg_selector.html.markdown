---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_selector"
sidebar_current: "docs-mso-resource-schema_template_external_epg_selector"
description: |-
  Manages MSO Schema Template External Endpoint Groups selectors.
---

# mso_schema_template_external_epg_selector#

Manages MSO Schema Template External Endpoint Groups Selectors.

## Example Usage ##
```hcl

resource "mso_schema_template_external_epg_selector" "selector1" {
	schema_id           = "${mso_schema_template_external_epg.template_externalepg.schema_id}"
	template_name       = "${mso_schema_template_external_epg.template_externalepg.template_name}"
	external_epg_name   = "${mso_schema_template_external_epg.template_externalepg.external_epg_name}"
	name                = "check"
    expressions {
      value = "1.20.30.44"
    }
    expressions{
      value = "5.6.7.8"
    }
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Subnet.
* `template_name` - (Required) Template where Anp Epg Subnet to be created.
* `external_epg_name` - (Required) Name of External Endpoint Group.
* `name` - (Required) Name for the selector.
* `expressions` - (Optional) expressions of Selector.
* `expressions.value` - (Optional) expression value for the selector.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template External Endpoint Groups Selector can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_external_epg_selector.selector1 {schema_id}/template/{template_name}/externalEPG/{external_epg_name}/selector/{name}
```