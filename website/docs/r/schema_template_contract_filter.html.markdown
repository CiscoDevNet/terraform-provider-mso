---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_filter"
sidebar_current: "docs-mso-resource-schema_template_contract_filter"
description: |-
  Manages MSO Schema Template Contract Filter.
---

# mso_schema_template_contract_filter #

Manages MSO Schema Template Contract Filter.

## Example Usage ##

```hcl

resource "mso_schema_template_contract_filter" "filter3" {
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
  contract_name = "Web-to-DB"
  filter_type   = "provider_to_consumer"
  filter_name   = "Any"
  directives    = ["none", "log"]
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Contract Filter.
* `template_name` - (Required) Template where Contract Filter to be created.
* `contract_name` - (Required) The name of the contract to manage. There should be an existing contract with this name.
* `filter_type` - (Required) The type of filters defined in this contract. Allowed values are `bothWay`, `provider_to_consumer` and `consumer_to_provider`.
* `filter_schema_id` - (Optional) The schemaId in which the filter is located. Default is `schema_id`.
* `filter_template_name` - (Optional) The template name in which the filter is located.  Default is `template_name`.
* `filter_name` - (Required) The filter name to associate with this contract. Filter must exist with the given `filter_name`, `filter_schema_id` and `filter_template_name`.
* `directives` - (Required) A list of filter directives. Allowed values are `log` and `none`.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Contract Filter can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_contract_filter.filter3 {schema_id}/template/{template_name}/contract/{contract_name}/filter/{filter_name}/filterType/{filter_type}
```