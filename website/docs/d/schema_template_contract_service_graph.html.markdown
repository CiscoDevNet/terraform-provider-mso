---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_service_graph"
sidebar_current: "docs-mso-data-source-schema_template_contract_service_graph"
description: |-
  Data source for MSO Schema Template Contract Service Graph.
---

# mso_schema_template_contract_service_graph #

Data source for MSO Schema Template Contract Service Graph.

## Example Usage ##

```hcl

data "mso_schema_template_contract_service_graph" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = data.mso_schema_template.t1.name
  contract_name = data.mso_schema_template_contract.c1.name
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID to associate the Template Service Graph with Contract.
* `template_name` - (Required) Template name to associate the Template Service Graph with Contract.
* `contract_name` - (Required) Contract name to associate the Template Service Graph.

## Attribute Reference ##
* `service_graph_name` - (Read-Only) Name of the Template Service Graph.
* `service_graph_schema_id` - (Read-Only) SchemaID of the source Template Service Graph.
* `service_graph_template_name` - (Read-Only) Template name of the source Template Service Graph.
* `node_relationship` - (Read-Only) Contract Service Graph Node relationship information.
  * `provider_connector_bd_name` - (Read-Only) Name of the BD that has to be connected to a Provider Connector.
  * `provider_connector_bd_schema_id` - (Read-Only) SchemaID of the source Provider Connector BD.
  * `provider_connector_bd_template_name` - (Read-Only) Template name of the source Provider Connector BD.
  * `consumer_connector_bd_name` - (Read-Only) Name of the BD that has to be connected to a Consumer Connector.
  * `consumer_connector_bd_schema_id` - (Read-Only) SchemaID of the source Consumer Connector BD.
  * `consumer_connector_bd_template_name` - (Read-Only) Template name of the source Consumer Connector BD.

