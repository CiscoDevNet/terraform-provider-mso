---
layout: "mso"
page_title: "MSO: mso_schema_site_contract_service_graph"
sidebar_current: "docs-mso-data-source-schema_site_contract_service_graph"
description: |-
  Data source for MSO Site Template Contract Service Graph.
---

# mso_schema_site_contract_service_graph #

Data source for MSO Site Template Contract Service Graph.

## Example Usage ##

```hcl

data "mso_schema_site_contract_service_graph" "example" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  contract_name = "C1"
  site_id       = mso_site.id
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Contract Service Graph.
* `template_name` - (Required) The template name of the Contract Service Graph.
* `contract_name` - (Required) The contract name of the Contract Service Graph.
* `site_id` - (Required) The site ID under which the Contract Service Graph is deployed.


## Attribute Reference ##

* `service_graph_schema_id` - (Read-Only) The schema ID of the Service Graph.
* `service_graph_template_name` - (Read-Only) The template name of the Service Graph.
* `service_graph_name` - (Read-Only) The name of the Service Graph.
* `node_relationship` - (Read-Only) The Site Template Contract Service Graph Node relationship information.
  * `provider_connector_cluster_interface` - (Read-Only) The name of the Cluster Interface that has to be connected to a Provider Connector.
  * `provider_connector_redirect_policy_tenant` - (Read-Only) The name of the Redirect Policy Tenant that has to be connected to a Provider Connector.
  * `provider_connector_redirect_policy` - (Read-Only) The name of the Redirect Policy that has to be connected to a Provider Connector.
  * `consumer_connector_cluster_interface` - (Read-Only) The name of the Cluster Interface that has to be connected to a Consumer Connector.
  * `consumer_connector_redirect_policy_tenant` - (Read-Only) The name of the Redirect Policy Tenant that has to be connected to a Consumer Connector.
  * `consumer_connector_redirect_policy` - (Read-Only) The name of the Redirect Policy that has to be connected to a Consumer Connector.
  * `consumer_subnet_ips` - (Read-Only) List of subnets connected to a Consumer Connector EPG.
