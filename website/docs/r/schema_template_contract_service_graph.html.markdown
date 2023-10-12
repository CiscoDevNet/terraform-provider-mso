---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_service_graph"
sidebar_current: "docs-mso-resource-schema_template_contract_service_graph"
description: |-
  Manages MSO Schema Template Contract Service Graph.
---

# mso_schema_template_contract_service_graph #

Manages MSO Schema Template Contract Service Graph.

## Example Usage ##

```hcl

resource "mso_schema_template_contract_service_graph" "example" {
  schema_id          = data.mso_schema.schema1.id
  template_name      = data.mso_schema_template.t1.name
  contract_name      = data.mso_schema_template_contract.c1.name
  service_graph_name = data.mso_schema_template_service_graph.sg1.service_graph_name
  node_relationship {
    consumer_connector_bd_name          = data.mso_schema_template_bd.bd1.name
    provider_connector_bd_template_name = data.mso_schema_template.t2.name
    provider_connector_bd_schema_id     = data.mso_schema.schema2.id
    provider_connector_bd_name          = data.mso_schema_template_bd.bd2.name
  }
  node_relationship {
    consumer_connector_bd_name = data.mso_schema_template_bd.bd1.name
    provider_connector_bd_name = data.mso_schema_template_bd.bd1.name
  }
  node_relationship {
    consumer_connector_bd_name = data.mso_schema_template_bd.bd1.name
    provider_connector_bd_name = data.mso_schema_template_bd.bd1.name
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID to associate the Template Service Graph with Contract.
* `template_name` - (Required) Template name to associate the Template Service Graph with Contract.
* `contract_name` - (Required) Contract name to associate the Template Service Graph.
* `service_graph_name` - (Required) Name of the Template Service Graph.
* `service_graph_schema_id` - (Optional) SchemaID of the source Template Service Graph. The `schema_id` will be used if not provided.
* `service_graph_template_name` - (Optional) Template name of the source Template Service Graph. The `template_name` will be used if not provided.
* `node_relationship` - (Required) Contract Service Graph Node relationship information.
  * `provider_connector_bd_name` - (Required) Name of the BD that has to be connected to a Provider Connector.
  * `provider_connector_bd_schema_id` - (Optional) SchemaID of the source Provider Connector BD. The `schema_id` will be used if not provided.
  * `provider_connector_bd_template_name` - (Optional) Template name of the source Provider Connector BD. The `template_name` will be used if not provided.
  * `consumer_connector_bd_name` - (Required) Name of the BD that has to be connected to a Consumer Connector.
  * `consumer_connector_bd_schema_id` - (Optional) SchemaID of the source Consumer Connector BD. The `schema_id` will be used if not provided.
  * `consumer_connector_bd_template_name` - (Optional) Template name of the source Consumer Connector BD. The `template_name` will be used if not provided.


## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Contract Service Graph can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_contract_service_graph.example {schema_id}/template/{template_name}/contract/{contract_name}
```