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
data "mso_schema_template_contract_filter" "filter1" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  contract_name = "c200"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Contract.
* `template_name` - (Required) Template where Contract to be created.
* `contract_name` - (Required) The name of the contract to manage.

## Attribute Reference ##

* `contract_schema_id` - (Optional) The schemaID that defines the referenced Contract.
* `contract_schema_template` - (Optional) The Template name that defines the referenced Contract.
* `display_name` - (Optional) Display Name of the contract on the MSO UI.
* `filter_type` - (Optional) The type of filters defined in this contract. Allowed values are `bothWay` and `oneWay`. Default to `bothWay`
* `scope` - (Optional) The scope of the contract.
* `filter_relationships` - (Optional) Map to provide Filter Relationships.
* `filter_relationships.filter_schema_id` - (Optional) The schemaId in which the filter is located.
* `filter_relationships.filter_template_name` - (Optional) The template name in which the filter is located.
* `filter_relationships.filter_name` - (Optional) The filter to associate with this contract.
* `directives` - (Optional) It is required if filter_relationships is passed, otherwise not required. A list of filter directives. Allowed values are `log` and `none`.
* `filter_relationships_provider_to_consumer` - (Optional) Map to provide Filter Relationships Provider to Consumer.
* `filter_relationships_provider_to_consumer.provider_to_consumer_schema_id` - (Optional) The schemaId in which the filter is located.
* `filter_relationships_provider_to_consumer.provider_to_consumer_template_name` - (Optional) The template name in which the filter is located.
* `filter_relationships_provider_to_consumer.provider_to_consumer_name` - (Optional) The filter to associate with this contract.
* `provider_to_consumer_directives` - (Optional) It is required if filter_relationships_provider_to_consumer is passed, otherwise not required. A list of filter directives. Allowed values are `log` and `none`.
* `filter_relationships_consumer_to_provider` - (Optional) Map to provide Filter Relationships Consumer to Provider.
* `filter_relationships_consumer_to_provider.consumer_to_provider_schema_id` - (Optional) The schemaId in which the filter is located.
* `filter_relationships_consumer_to_provider.consumer_to_provider_template_name` - (Optional) The template name in which the filter is located.
* `filter_relationships_consumer_to_provider.consumer_to_provider_name` - (Optional) The filter to associate with this contract.
* `consumer_to_provider_directives` - (Optional) It is required if filter_relationships_consumer_to_provider is passed, otherwise not required. A list of filter directives. Allowed values are `log` and `none`.
 
