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
  schema_id             = data.mso_schema.schema1.id
  site_id               = data.mso_site.site1.id
  template_name         = "Template1"
  contract_name         = "UntitledContract1"
  service_graph_name    = "sg1"  
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Service Graph.
* `site_id` - (Required) The site ID of the Service Graph.
* `template_name` - (Required) The template name of the Service Graph.
* `contract_name` - (Required) The name of the Contract.
* `service_graph_name` - (Required) The name of Service Graph.


## Attribute Reference ##

* `node_relationship` - (Read-Only) A list of node relationships of the Service Graph.
    * `provider_connector_bd_name` - (Read-Only) The BD name of the provider connector at template level.
    * `provider_connector_bd_schema_id` - (Read-Only) The BD schema ID of the provider connector at template level.
    * `provider_connector_bd_template_name` - (Read-Only) The BD template name of the provider connector at template level.
    * `consumer_connector_bd_name` - (Read-Only) The BD name of the consumer connector at template level.
    * `consumer_connector_bd_schema_id` - (Read-Only) The BD schema ID of the consumer connector at template level.
    * `consumer_connector_bd_template_name` - (Read-Only) The BD template name of the consumer connector at template level.
    * `provider_connector_cluster_interface` - (Read-Only) The cluster interface for the provider connector at site level. 
    * `consumer_connector_cluster_interface` - (Read-Only) The cluster interface for the consumer connector at site level.
    * `provider_connector_redirect_policy_tenant` - (Read-Only) The tenant redirection policy for the provider connector at site level. 
    * `provider_connector_redirect_policy` - (Read-Only) The redirection policy for the provider connector at site level.
    * `consumer_connector_redirect_policy_tenant` - (Read-Only) The tenant redirection policy for the consumer connector at site level. 
    * `consumer_connector_redirect_policy` - (Read-Only) The redirection policy for the consumer connector at site level.
    * `provider_subnet_ips` - (Read-Only) A list of subnet ips associated with the provider connector at site level.
    * `consumer_subnet_ips` - (Read-Only) A list of subnet ips associated with the consumer connector at site level.
