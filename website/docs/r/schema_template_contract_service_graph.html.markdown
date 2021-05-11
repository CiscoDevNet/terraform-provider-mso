---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_service_graph"
sidebar_current: "docs-mso-resource-schema_template_contract_service_graph"
description: |-
  Manages MSO Schema Template Contract service graph.
---

# mso_schema_template_contract_service_graph #

Manages MSO Schema Template Contract service graph.

## Example Usage ##

```hcl

resource "mso_schema_template_contract_service_graph" "one" {
  schema_id               = "5f11b0e22c00001c4a812a2a"
  site_id                 = "5c7c95b25100008f01c1ee3c"
  template_name           = "Template1"
  contract_name           = "UntitledContract1"
  service_graph_name      = "sg1"
  service_graph_schema_id = "5f16a7c62c00006367812a2f"
  node_relationship {
    provider_connector_bd_name                = "BD1"
    consumer_connector_bd_name                = "BD2"
    provider_connector_cluster_interface      = "test"
    consumer_connector_cluster_interface      = "test"
    provider_connector_redirect_policy_tenant = "NkAutomation"
    provider_connector_redirect_policy        = "test2"
    consumer_connector_redirect_policy_tenant = "NkAutomation"
    consumer_connector_redirect_policy        = "test2"
    provider_subnet_ips = ["1.2.3.4/20"]
    consumer_subnet_ips = ["1.2.3.4/20"]
  }
  node_relationship {
    provider_connector_bd_name          = "CBD2"
    provider_connector_bd_schema_id     = "5f16a7c62c00006367812a2f"
    provider_connector_bd_template_name = "Template1"
    consumer_connector_bd_name          = "CBD1"
    consumer_connector_bd_schema_id     = "5f16a7c62c00006367812a2f"
    consumer_connector_bd_template_name = "Template1"
    provider_connector_cluster_interface = "test"
    consumer_connector_cluster_interface = "test"
    provider_connector_redirect_policy_tenant = "NkAutomation"
    provider_connector_redirect_policy        = "test2"
    consumer_connector_redirect_policy_tenant = "NkAutomation"
    consumer_connector_redirect_policy        = "test2"
    provider_subnet_ips = ["1.2.3.4/20"]
    consumer_subnet_ips = ["1.2.3.4/20"]
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Contract Service Graph.
* `site_id` - (Required) SiteID under which you want to deploy Contract Service Graph.
* `template_name` - (Required) Template where Contract Service Graph to be created.
* `contract_name` - (Required) The name of the contract to manage. There should be an existing contract with this name.

* `service_graph_name` - (Required) The name of service graph.
* `service_graph_schema_id` - (Optional) The schema Id in which service graph is created. If not given then `schema_id` will be taken.
* `service_graph_template_name` - (Optional) The Template name in which service graph is created. If not given then `template_name` will be taken.
* `service_graph_site_id` - (Optional) The Site Id for where service graph created. If not given then `site_id` will be taken.

* `node_relationship` - (Required) Service graph node relationship information. You have to define this block for every node of service graph.
* `node_relationship.provider_connector_bd_name` - (Required) bd name for provider connector at template level.
* `node_relationship.provider_connector_bd_schema_id` - (Optional) schema id under which above bd is created. If not given then `schema_id` will be taken.
* `node_relationship.provider_connector_bd_template_name` - (Optional) template name under which above bd is created. If not given then `template_name` will be taken.

* `node_relationship.consumer_connector_bd_name` - (Required) bd name for consumer connector at template level.
* `node_relationship.consumer_connector_bd_schema_id` - (Optional) schema id under which above bd is created. If not given then `schema_id` will be taken.
* `node_relationship.consumer_connector_bd_template_name` - (Optional) template name under which above bd is created. If not given then `template_name` will be taken.

* `node_relationship.provider_connector_cluster_interface` - (Required) cluster interface for provider connector to attach with node at site level. 
* `node_relationship.consumer_connector_cluster_interface` - (Required) cluster interface for consumer connector to attach with node at site level.

* `node_relationship.provider_connector_redirect_policy_tenant` - (Optional) tenant for redirection policy for provider connector at site level. It is required to set redirection policy for provider connector.
* `node_relationship.provider_connector_redirect_policy` - (Optional) redirection policy for provider connector at site level.
* `node_relationship.consumer_connector_redirect_policy_tenant` - (Optional) tenant for redirection policy for consumer connector at site level. It is required to set redirection policy for consumer connector.
* `node_relationship.consumer_connector_redirect_policy` - (Optional) redirection policy for consumer connector at site level.

* `node_relationship.provider_subnet_ips` - (Optional) subnet ips which will be associated with provider connector at site level. It should be in CIDR format.
* `node_relationship.consumer_subnet_ips` - (Optional) subnet ips which will be associated with consumer connector at site level. It should be in CIDR format.



## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Contract service graph can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_contract_service_graph.one {schema_id}/template/{template_name}/contract/{contract_name}/serviceGraph/{service_graph_name}
```