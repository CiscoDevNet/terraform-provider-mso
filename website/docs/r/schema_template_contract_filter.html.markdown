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
resource "mso_schema_template_contract_filter" "filter1" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  contract_name = "c200"
  display_name = "c200"
  filter_type = "bothWay"
  scope = "context"
  filter_relationships_provider_to_consumer = {
   provider_to_consumer_schema_id = "5c4d5bb72700000401f80948"
   provider_to_consumer_template_name = "Template1"
   provider_to_consumer_name = "many"
  }
  provider_to_consumer_directives = ["log","none","log"]
  filter_relationships_consumer_to_provider = {
    consumer_to_provider_schema_id = "5c4d5bb72700000401f80948"
    consumer_to_provider_template_name = "Template1"
    consumer_to_provider_name = "many"
  }
  consumer_to_provider_directives = ["log","none"]
 
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Contract.
* `template_name` - (Required) Template where Contract to be created.
* `contract_name` - (Required) The name of the contract to manage.
* `contract_schema_id` - (Optional) The schemaID that defines the referenced Contract.
* `contract_schema_template` - (Optional) The Template name that defines the referenced Contract.
* `display_name` - (Optional) Display Name of the contract on the MSO UI.
* `filter_type` - (Optional) The type of filters defined in this contract. Allowed values are `bothWay` and `oneWay`. Default to `bothWay`
* `scope` - (Optional) The scope of the contract.
* `filter_relationships` - (Optional) Map to provide Filter Relationships.
* `filter_schema_id` - (Optional) The schemaId in which the filter is located.
* `filter_template_name` - (Optional) The template name in which the filter is located.
* `filter_name` - (Required) The filter to associate with this contract.
* `directives` - (Required) A list of filter directives. Allowed values are `log` and `none`.
* `filter_relationships_provider_to_consumer` - (Required) Map to provide Filter Relationships Provider to Consumer.
* `provider_to_consumer_schema_id` - (Optional) The schemaId in which the filter is located.
* `provider_to_consumer_template_name` - (Optional) The template name in which the filter is located.
* `provider_to_consumer_name` - (Required) The filter to associate with this contract.
* `provider_to_consumer_directives` - (Optional) A list of filter directives. Allowed values are `log` and `none`.
* `filter_relationships_consumer_to_provider` - (Required) Map to provide Filter Relationships Consumer to Provider.
* `consumer_to_provider_schema_id` - (Optional) The schemaId in which the filter is located.
* `consumer_to_provider_template_name` - (Optional) The template name in which the filter is located.
* `consumer_to_provider_name` - (Required) The filter to associate with this contract.
* `consumer_to_provider_directives` - (Optional) A list of filter directives. Allowed values are `log` and `none`.

## Attribute Reference ##

No attributes are exported.
