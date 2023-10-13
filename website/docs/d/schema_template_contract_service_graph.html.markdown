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
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
  contract_name = "C1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Contract Service Graph.
* `template_name` - (Required) The template name of the Contract Service Graph.
* `contract_name` - (Required) The contract name of the Contract Service Graph.

## Attribute Reference ##

* `service_graph_name` - (Read-Only) The name of the Service Graph.
* `service_graph_schema_id` - (Read-Only) The schema ID of the Service Graph.
* `service_graph_template_name` - (Read-Only) The template name of the Service Graph.
* `node_relationship` - (Read-Only) The Contract Service Graph Node relationship information.
  * `provider_connector_bd_name` - (Read-Only) The name of the BD that has to be connected to a Provider Connector.
  * `provider_connector_bd_schema_id` - (Read-Only) The schema ID of the BD that has to be connected to a Provider Connector.
  * `provider_connector_bd_template_name` - (Read-Only) The template name of the BD that has to be connected to a Provider Connector.
  * `consumer_connector_bd_name` - (Read-Only) The name of the BD that has to be connected to a Consumer Connector.
  * `consumer_connector_bd_schema_id` - (Read-Only) The schema ID of the BD that has to be connected to a Consumer Connector.
  * `consumer_connector_bd_template_name` - (Read-Only) The template name of the BD that has to be connected to a Consumer Connector.
