---
layout: "mso"
page_title: "MSO: mso_schema_site_template_contract_service_graph"
sidebar_current: "docs-mso-resource-schema_site_template_contract_service_graph"
description: |-
  Manages MSO Site Template Contract Service Graph.
---

# mso_schema_site_template_contract_service_graph #

Manages MSO Site Template Contract Service Graph.

## Example Usage ##

```hcl

resource "mso_schema_site_template_contract_service_graph" "example" {
  schema_id          = mso_schema.schema.id
  template_name      = "Template1"
  contract_name      = "C1"
  service_graph_name = "SG1"
  site_id            = mso_site.id
  node_relationship {
    provider_connector_cluster_interface      = "example_provider_cluster_interface"
    provider_connector_redirect_policy_tenant = "example_tenant"
    provider_connector_redirect_policy        = "example_redirect_policy"
    consumer_connector_cluster_interface      = "example_consumer_cluster_interface"
    consumer_connector_redirect_policy_tenant = "example_tenant"
    consumer_connector_redirect_policy        = "example_redirect_policy"
    consumer_subnet_ips                       = ["1.1.1.1/24", "2.2.2.2/24"]
  }
  node_relationship {
    provider_connector_cluster_interface = "example_provider_cluster_interface"
    consumer_connector_cluster_interface = "example_consumer_cluster_interface"
  }
}

# Cloud Network Controller site configuration
resource "mso_schema_site_template_contract_service_graph" "example" {
  schema_id          = mso_schema.schema.id
  template_name      = "Template1"
  contract_name      = "C1"
  service_graph_name = "SG1"
  site_id            = mso_site.id
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID of the Contract Service Graph.
* `template_name` - (Required) The template name of the Contract Service Graph.
* `contract_name` - (Required) The contract name of the Contract Service Graph.
* `site_id` - (Required) The site ID under which the Contract Service Graph is deployed.
* `service_graph_name` - (Required) The name of the Service Graph.
* `service_graph_schema_id` - (Optional) The schema ID of the Service Graph. The `schema_id` will be used if not provided.
* `service_graph_template_name` - (Optional) The template name of the Service Graph. The `template_name` will be used if not provided.
* `node_relationship` - (Optional) The Site Template Contract Service Graph Node relationship information. The order of the node_relationship object should match the node types in the Service Graph. **The `node_relationship` is not supported for the Cloud Network Controller site.**
  * `provider_connector_cluster_interface` - (Required) The name of the Cluster Interface that has to be connected to a Provider Connector.
  * `provider_connector_redirect_policy_tenant` - (Optional) The name of the Redirect Policy Tenant that has to be connected to a Provider Connector.
  * `provider_connector_redirect_policy` - (Optional) The name of the Redirect Policy that has to be connected to a Provider Connector.
  * `consumer_connector_cluster_interface` - (Required) The name of the Cluster Interface that has to be connected to a Consumer Connector.
  * `consumer_connector_redirect_policy_tenant` - (Optional) The name of the Redirect Policy Tenant that has to be connected to a Consumer Connector.
  * `consumer_connector_redirect_policy` - (Optional) The name of the Redirect Policy that has to be connected to a Consumer Connector.
  * `consumer_subnet_ips` - (Optional) List of subnets connected to a Consumer Connector EPG. Only supported for the load balancer device.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Site Template Contract Service Graph can be [imported][docs-import] into this resource using its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_template_contract_service_graph.example {schema_id}/sites/{site_id}/templates/{template_name}/contracts/{contract_name}
```
