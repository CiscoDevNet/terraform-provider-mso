---
layout: "mso"
page_title: "MSO: mso_schema_template_contract"
sidebar_current: "docs-mso-data-source-schema_template_contract"
description: |-
  Data source for MSO Schema Template Contract.
---

# mso_schema_template_contract #

Data source for MSO Schema Template Contract.

## Example Usage ##

```hcl
data "mso_schema_template_contract" "scontract1" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  contract_name = "c1"
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
* `filter_schema_id` - (Optional) The schemaId in which the filter is located.
* `filter_template_name` - (Optional) The template name in which the filter is located.
* `filter_name` - (Optional) The filter to associate with this contract.
* `directives` - (Optional) A list of filter directives. Allowed values are `log` and `none`.
