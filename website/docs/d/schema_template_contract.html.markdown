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

data "mso_schema_template_contract" "contract1" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  contract_name = "c1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Contract.
* `template_name` - (Required) Template where Contract to be created.
* `contract_name` - (Required) The name of the contract to manage.

## Attribute Reference ##

* `display_name` - (Optional) Display Name of the contract on the MSO UI.
* `filter_type` - (Optional) The type of filters defined in this contract. Allowed values are `bothWay` and `oneWay`. Default to `bothWay`
* `scope` - (Optional) The scope of the contract.
* `filter_relationships` - **Deprecated** (Required if filter_relationship is not used) Map to provide Filter Relationships. This attribute is deprecated, use `filter_relationship` instead.
  * `filter_relationships.filter_schema_id` - (Optional) The schemaId in which the filter is located.
  * `filter_relationships.filter_template_name` - (Optional) The template name in which the filter is located.
  * `filter_relationships.filter_name` - (Required) The filter to associate with this contract.

* `filter_relationship` - (Required if filter_relationships is not used) Map to provide Filter Relationships.
  * `filter_relationship.filter_schema_id` - (Optional) The schemaId in which the filter is located.
  * `filter_relationship.filter_template_name` - (Optional) The template name in which the filter is located.
  * `filter_relationship.filter_name` - (Required) The filter to associate with this contract.

* `directives` - (Optional) A list of filter directives. Allowed values are `log` and `none`.
