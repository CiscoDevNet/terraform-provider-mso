---
layout: "mso"
page_title: "MSO: mso_schema_template_contract"
sidebar_current: "docs-mso-resource-schema_template_contract"
description: |-
  Manages MSO Schema Template Contract.
---

# mso_schema_template_contract #

Manages MSO Schema Template Contract.

## Example Usage ##

```hcl

resource "mso_schema_template_contract" "example" {
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
  contract_name = "C1"
  display_name  = "C1"
  filter_type   = "bothWay"
  scope         = "context"
  target_dscp   = "af11"
  filter_relationship {
    filter_schema_id     = mso_schema_template_filter_entry.filter_entry.schema_id
    filter_template_name = "Template1"
    filter_name          = mso_schema_template_filter_entry.filter_entry.name
    filter_type          = "bothWay"
    directives           = ["log", "no_stats"]
    action               = "deny"
    priority             = "level1"
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which you want to deploy Contract.
* `template_name` - (Required) The template name under which you want to deploy Contract.
* `contract_name` - (Required) The name of the Contract.
* `display_name` - (Optional) The display name of the Contract.
* `filter_type` - (Optional)  The type of filters assigned to the Contract. Allowed values are `bothWay` and `oneWay`. Default to `bothWay`.
* `scope` - (Optional) The scope of the Contract. Allowed values are `application-profile`, `tenant`, `context`, and `global`. Default to `context`.
* `target_dscp` - (Optional) The dscp value of the Contract. Allowed values are `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0`, `cs1`, `cs2`, `cs3`, `cs4`, `cs5`, `cs6`, `cs7`, `expeditedForwarding`, `voiceAdmit`, and `unspecified`. Defaults to `unspecified`.
* `priority` - (Optional) The priority of the Contract. Allowed values are `unspecified`, `level1`, `level2`, `level3`, `level4`, `level5`, and `level6`. Defaults to `unspecified`.
* `filter_relationship` - (Optional) A list of Filter Relationships for the Contract.
  * `filter_schema_id` - (Optional) The schema ID of the Filter associated with the Contract.
  * `filter_template_name` - (Optional) The template name of the Filter associated with the Contract.
  * `filter_name` - (Required) The name of the Filter associated with the Contract.
  * `filter_type` - (Optional) The type of the Filter associated with the Contract. Allowed values are `bothWay`, `consumer_to_provider` and `provider_to_consumer`. Defaults to `bothWay`.
  * `directives` - (Optional)  A list of filter directives associated with the Contract. Allowed values are `none`, `no_stats`, and `log`.
  * `action` - (Optional) The action of the Filter associated with the Contract. Allowed values are `deny` and `permit`. 
  * `priority` - (Optional) The override priority of the Filter associated with the Contract. Allowed values are `default`, `level1`, `level2`, and `level3`. 
  
* `filter_relationships` - (Optional) **Deprecated** A Map to provide one Filter Relationship. This attribute is deprecated, use `filter_relationship` instead. It is not allowed to use in combination with `filter_relationship`.
  * `filter_schema_id` - (Optional) The schemaId in which the filter is located.
  * `filter_template_name` - (Optional) The template name in which the filter is located.
  * `filter_name` - (Required) The filter to associate with this contract.
* `directives` -  (Optional) **Deprecated** A list of filter directives. Allowed values are `none`, `no_stats`, and `log`.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Contract can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_contract.example {schema_id}/templates/{template_name}/contracts/{contract_name}
```