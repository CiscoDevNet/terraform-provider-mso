---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_filter"
sidebar_current: "docs-mso-data-source-schema_template_contract_filter"
description: |-
  Data source for MSO Schema Template Contract Filter.
---

# mso_schema_template_contract-filter #

Data source for MSO Schema Template Contract Filter.

## Example Usage ##

```hcl
data "mso_schema_template_contract_filter" "f1" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  contract_name = "Web-to-DB"
  filter_type = "provider_to_consumer"
  filter_name = "Any"
  filter_schema_id= "5c4d5bb72700000401f80948"
  filter_template_name = "Template1"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Contract Filter.
* `template_name` - (Required) Template where Contract Filter to be created.
* `contract_name` - (Required) The name of the contract to manage. There should be an existing contract with this name.
* `filter_type` - (Required) The type of filters defined in this contract. Allowed values are `bothWay`, `provider_to_consumer` and `consumer_to_provider`.
* `filter_schema_id` - (Required) The schemaId in which the filter is located.
* `filter_template_name` - (Required) The template name in which the filter is located.
* `filter_name` - (Required) The filter name to associate with this contract. Filter must exist with the given `filter_name`, `filter_schema_id` and `filter_template_name` Force New set to `true`.


## Attribute Reference ##

* `directives` - (Optional) A list of filter directives. Allowed values are `log` and `none`.


