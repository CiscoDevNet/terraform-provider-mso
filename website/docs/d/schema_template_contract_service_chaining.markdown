---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_service_chaining"
sidebar_current: "docs-mso-data-source-schema_template_contract_service_chaining"
description: |-
  Data source for Schema Template Contract Service Chaining.
---

# mso_schema_template_contract_service_chaining #

Data source for a Schema Template Contract Service Chaining configuration on Cisco Nexus Dashboard Orchestrator (NDO).

## GUI Information ##

* `Location` - Manage -> Manage -> Tenant Template -> Applications -> Schemas -> <Schema Name> -> <Template Name> -> Contracts -> <Contract Name> -> Service Chaining

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
* `node_filter` - (Read-Only) The node filter configured for the service chain.
* `service_nodes` - (Read-Only) A list of service nodes that form the service chain. Each element has the following attributes:
  * `name` - (Read-Only) The name of the service node.
  * `device_type` - (Read-Only) The type of the service device (e.g., firewall, loadBalancer).
  * `device_ref` - (Read-Only) The NDO UUID of the mso_service_device_cluster used for this node.
  * `index` - (Read-Only) The computed order of the node in the service chain.
  * `uuid` - (Read-Only) The NDO UUID of the service node instance within the chain.
  * `consumer_connector` - (Read-Only) A list containing the consumer-side connection block.
    * `interface_name` - (Read-Only) The name of the consumer connector interface.
    * `is_redirect` - (Read-Only) Whether the consumer connector is a redirect.
  * `provider_connector` - (Read-Only) A list containing the provider-side connection block.
    * `interface_name` - (Read-Only) The name of the provider connector interface.
    * `is_redirect` - (Read-Only) Whether the provider connector is a redirect.
