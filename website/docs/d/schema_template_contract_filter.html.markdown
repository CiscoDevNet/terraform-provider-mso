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

data "mso_schema_template_contract_filter" "example" {
  schema_id            = data.mso_schema.schema1.id
  template_name        = "Template1"
  contract_name        = "Web-to-DB"
  filter_type          = "provider_to_consumer"
  filter_name          = "Any"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Contract.
* `template_name` - (Required) The template name of the Contract.
* `contract_name` - (Required) The name of the Contract.
* `filter_type` - (Required) The type of the Filter. Allowed values are `bothWay`, `provider_to_consumer` and `consumer_to_provider`.
* `filter_name` - (Required) The name of the Filter.
* `filter_schema_id` - (Optional) The schema ID of the Filter. The `schema_id` of the Contract will be used if not provided. 
* `filter_template_name` - (Optional) The template name of the Filter. The `template_name` of the Contract will be used if not provided. 


## Attribute Reference ##

* `directives` - (Read-Only) A list of filter directives.
