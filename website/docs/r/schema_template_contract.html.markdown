---
layout: "mso"
page_title: "MSO: mso_schema_template_contract"
sidebar_current: "docs-mso-resource-schema_template_contract"
description: |-
  Manages MSO Schema Template Bridge Domain.
---

# mso_schema_template_contract #

Manages MSO Schema Template Bridge Domain.

## Example Usage ##

```hcl
resource "mso_schema_template_contract" "template_contract" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  contract_name = "C1"
  display_name = "C1"
  filter_type = "bothWay"
  scope = "context"
  filter_relationships = {
    filter_schema_id = "5ea809672c00003bc40a2799"
    filter_template_name = "Template1"
    filter_name = "filter1"
  }
  directives = ["none"]
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
* `filter_relationships` - (Required) Map to provide Filter Relationships.
* `filter_schema_id` - (Optional) The schemaId in which the filter is located.
* `filter_template_name` - (Optional) The template name in which the filter is located.
* `filter_name` - (Required) The filter to associate with this contract.
* `directives` - (Required) A list of filter directives. Allowed values are `log` and `none`.

## Attribute Reference ##

No attributes are exported.
