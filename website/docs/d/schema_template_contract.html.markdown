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

data "mso_schema_template_contract" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  contract_name = "c1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Contract.
* `template_name` - (Required) The template name of the Contract.
* `contract_name` - (Required) The name of the Contract.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the Contract as displayed on the MSO UI.
* `filter_type` - (Read-Only) The type of filters of the Contract.
* `scope` - (Read-Only) The scope of the Contract.
* `target_dscp` - (Read-Only) The dscp value of the Contract.
* `priority` - (Read-Only) The priority override of the Filter.
* `filter_relationship` - (Read-Only) A List of Filter relationships.
    * `filter_schema_id` - (Read-Only) The schema ID of the Filter.
    * `filter_template_name` - (Read-Only) The template name of the Filter.
    * `filter_name` - (Read-Only) The name of the Filter.
    * `filter_type` - (Read-Only) The type of the Filter. 
    * `action` - (Read-Only) The action of the Filter.
    * `directives` - (Read-Only) The directives of the Filter.
    * `priority` - (Read-Only) The priority override of the Filter.

* `filter_relationships` - (Read-Only) **Deprecated** A map of the Filter relationship.
    * `filter_schema_id` - (Read-Only) The schema ID of the Filter.
    * `filter_template_name` - (Read-Only) The template name of the Filter.
    * `filter_name` - (Read-Only) The name of the Filter.
* `directives` - (Read-Only) **Deprecated** The directives of the Filter.
* `description` - (Read-Only) The description of the Contract.