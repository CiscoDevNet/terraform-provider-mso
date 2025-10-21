---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_service_chaining"
sidebar_current: "docs-mso-data-source-schema_template_contract_service_chaining"
description: |-
  Data source for Schema Template Contract Service Chaining.
---

# mso_schema_template_contract_service_chaining #

Data source for a Schema Template Contract Service Chaining configuration on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.2(3) and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Applications -> Schemas -> <Schema Name> -> <Template Name> -> Contracts -> <Contract Name> -> Service Chaining

## Example Usage ##

```hcl
data "mso_schema_template_contract_service_chaining" "chain" {
  schema_id     = "a1b2c3d4-e5f6-7890-1234-567890abcdef"
  template_name = "Template1"
  contract_name = "WebAppContract"
}
```

## Argument Reference ##

* `schema_id` - (Required) The ID of the schema where the contract resides.
* `template_name` - (Required) The name of the template where the contract resides.
* `contract_name` - (Required) The name of the contract to look up the service chain for.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the service chain.
* `name` - (Read-Only) The name of the service chain, derived from the contract name.
* `node_filter` - (Read-Only) Specifies the name of a filter used to selectively redirect a subset of the contract-permitted traffic through the service chain.
* `service_nodes` - (Read-Only) A list of the service nodes that constitute the service chain, presented in their processing order. Each element details the configuration of a single service node.
  * `name` - (Read-Only) The name of the service node.
  * `device_type` - (Read-Only) The type of the service device (e.g., firewall, loadBalancer).
  * `device_ref` - (Read-Only) The NDO UUID of the mso_service_device_cluster used for this node.
  * `index` - (Read-Only) The computed order of the node in the service chain.
  * `uuid` - (Read-Only) The NDO UUID of the service node instance within the chain.
  * `consumer_connector` - (Read-Only) A list containing the consumer-side connection block.
    * `interface_name` - (Read-Only) The name of the consumer connector interface.
    * `is_redirect` - (Read-Only) When is_redirect is set to true, the consumer_connector specifies the interface that receives traffic diverted by a policy, rather than traffic flowing directly through the service device.
  * `provider_connector` - (Read-Only) A list containing the provider-side connection block.
    * `interface_name` - (Read-Only) The name of the provider connector interface.
    * `is_redirect` - (Read-Only) When is_redirect is set to true, the provider_connector specifies the interface used to send traffic back into the network fabric after it has been processed by the service device.
