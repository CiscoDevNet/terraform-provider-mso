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

data "mso_schema_template_contract_service_graph" "name" {
  schema_id             = data.mso_schema.schema1.id
  site_id               = data.mso_site.site1.id
  template_name         = "Template1"
  contract_name         = "UntitledContract1"
  service_graph_name    = "sg1"  
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Contract Service Graph.
* `site_id` - (Required) SiteID under which you want to deploy Contract Service Graph.
* `template_name` - (Required) Template where Contract Service Graph to be created.
* `contract_name` - (Required) The name of the contract to manage. There should be an existing contract with this name.
* `service_graph_name` - (Required) The name of service graph.


## Attribute Reference ##

* `node_relationship` - (Optional) Service graph node relationship information.
* `node_relationship.provider_connector_bd_name` - (Optional) bd name for provider connector at template level.
* `node_relationship.consumer_connector_bd_name` - (Optional) bd name for consumer connector at template level.
* `node_relationship.provider_connector_cluster_interface` - (Optional) cluster interface for provider connector to attach with node at site level. 
* `node_relationship.consumer_connector_cluster_interface` - (Optional) cluster interface for consumer connector to attach with node at site level.
* `node_relationship.provider_connector_redirect_policy_tenant` - (Optional) tenant for redirection policy for provider connector at site level. It is required to set redirection policy for provider connector.
* `node_relationship.provider_connector_redirect_policy` - (Optional) redirection policy for provider connector at site level.
* `node_relationship.consumer_connector_redirect_policy_tenant` - (Optional) tenant for redirection policy for consumer connector at site level. It is required to set redirection policy for consumer connector.
* `node_relationship.consumer_connector_redirect_policy` - (Optional) redirection policy for consumer connector at site level.
* `node_relationship.provider_subnet_ips` - (Optional) subnet ips which will be associated with provider connector at site level. It should be in CIDR format.
* `node_relationship.consumer_subnet_ips` - (Optional) subnet ips which will be associated with consumer connector at site level. It should be in CIDR format.

