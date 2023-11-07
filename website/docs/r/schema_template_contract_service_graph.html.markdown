---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_service_graph"
sidebar_current: "docs-mso-resource-schema_template_contract_service_graph"
description: |-
  Manages MSO Schema Template Contract Service Graph.
---

# mso_schema_template_contract_service_graph #

Manages MSO Schema Template Contract Service Graph.

# Note: #
This resource is only compatible with NDO versions 3.7 and 4.2+. NDO versions 4.0 and 4.1 are not supported.

## Example Usage ##

```hcl

resource "mso_schema_template_contract_service_graph" "example" {
  schema_id          = mso_schema.schema1.id
  template_name      = "Template1"
  contract_name      = "C1"
  service_graph_name = "SG1"
  node_relationship {
    consumer_connector_bd_name          = "BD1"
    provider_connector_bd_schema_id     = mso_schema.schema2.id
    provider_connector_bd_template_name = "Template2"
    provider_connector_bd_name          = "BD2"
  }
  node_relationship {
    consumer_connector_bd_name = "BD1"
    provider_connector_bd_name = "BD2"
  }
  node_relationship {
    consumer_connector_bd_name = "BD1"
    provider_connector_bd_name = "BD2"
  }
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID of the Contract Service Graph.
* `template_name` - (Required) The template name of the Contract Service Graph.
* `contract_name` - (Required) The contract name of the Contract Service Graph.
* `service_graph_name` - (Required) The name of the Service Graph.
* `service_graph_schema_id` - (Optional) The schema ID of the Service Graph. The `schema_id` will be used if not provided.
* `service_graph_template_name` - (Optional) The template name of the Service Graph. The `template_name` will be used if not provided.
* `node_relationship` - (Required) The Contract Service Graph Node relationship information. The order of the node_relationship object should match the node types in the Service Graph.
  * `provider_connector_bd_name` - (Required) The name of the BD that has to be connected to a Provider Connector.
  * `provider_connector_bd_schema_id` - (Optional) The schema ID of the BD that has to be connected to a Provider Connector. The `schema_id` will be used if not provided.
  * `provider_connector_bd_template_name` - (Optional) The template name of the BD that has to be connected to a Provider Connector. The `template_name` will be used if not provided.
  * `consumer_connector_bd_name` - (Required) The name of the BD that has to be connected to a Consumer Connector.
  * `consumer_connector_bd_schema_id` - (Optional) The schema ID of the BD that has to be connected to a Consumer Connector. The `schema_id` will be used if not provided.
  * `consumer_connector_bd_template_name` - (Optional) The template name of the BD that has to be connected to a Consumer Connector. The `template_name` will be used if not provided.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Contract Service Graph can be [imported][docs-import] into this resource using its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_contract_service_graph.example {schema_id}/templates/{template_name}/contracts/{contract_name}
```