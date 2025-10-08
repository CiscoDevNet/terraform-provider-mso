---
layout: "mso"
page_title: "MSO: mso_schema_template_contract_service_chaining"
sidebar_current: "docs-mso-resource-schema_template_contract_service_chaining"
description: |-
  Manages Schema Template Contract Service Chaining on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_schema_template_contract_service_chaining #

Manages Schema Template Contract Service Chaining on Cisco Nexus Dashboard Orchestrator (NDO).

## GUI Information ##

* `Location` - Manage -> Manage -> Tenant Template -> Applications -> Schemas -> <Schema Name> -> <Template Name> -> Contracts -> <Contract Name> -> Service Chaining

## Example Usage ##

```hcl
# This example creates a full stack of dependencies for service chaining,
# including templates, a schema, VRF, BDs, a contract, and service device clusters.

resource "mso_tenant" "tenant" {
  name         = "ServiceChainTenant"
  display_name = "ServiceChainTenant"
}

resource "mso_template" "device_template" {
  template_name = "DeviceTemplateForSC"
  template_type = "service_device"
  tenant_id     = mso_tenant.tenant.id
}

resource "mso_schema" "schema" {
  name = "SchemaForServiceChaining"
  template {
    name          = "Template1"
    display_name  = "Template1"
    tenant_id     = mso_tenant.tenant.id
    template_type = "aci_multi_site"
  }
}

resource "mso_schema_template_vrf" "vrf" {
  schema_id    = mso_schema.schema.id
  template     = "Template1"
  name         = "SC_VRF"
  display_name = "SC_VRF"
}

resource "mso_schema_template_bd" "bd1" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  name          = "SC_BD1"
  vrf_name      = mso_schema_template_vrf.vrf.name
}

resource "mso_schema_template_bd" "bd2" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  name          = "SC_BD2"
  vrf_name      = mso_schema_template_vrf.vrf.name
}

resource "mso_service_device_cluster" "fw_device" {
  template_id = mso_template.device_template.id
  name        = "FirewallCluster"
  device_mode = "layer3"
  device_type = "firewall"

  interface_properties {
    name    = "fw_interface"
    bd_uuid = mso_schema_template_bd.bd1.uuid
  }
}

resource "mso_service_device_cluster" "lb_device" {
  template_id = mso_template.device_template.id
  name        = "LoadBalancerCluster"
  device_mode = "layer3"
  device_type = "loadBalancer"

  interface_properties {
    name    = "lb_prov_if"
    bd_uuid = mso_schema_template_bd.bd1.uuid
  }
  interface_properties {
    name    = "lb_cons_if"
    bd_uuid = mso_schema_template_bd.bd2.uuid
  }
}

resource "mso_schema_template_contract" "contract" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  contract_name = "WebAppContract"
  display_name  = "WebAppContract"
  scope         = "context"
}

# Main resource for the Service Chain
resource "mso_schema_template_contract_service_chaining" "chain" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  contract_name = mso_schema_template_contract.contract.contract_name
  node_filter   = "allow-all"

  service_nodes {
    name        = "firewall-node"
    device_type = "firewall"
    device_ref  = mso_service_device_cluster.fw_device.uuid

    consumer_connector {
      interface_name = "fw_interface"
    }
    provider_connector {
      interface_name = "fw_interface"
    }
  }

  service_nodes {
    name        = "loadbalancer-node"
    device_type = "loadBalancer"
    device_ref  = mso_service_device_cluster.lb_device.uuid

    consumer_connector {
      interface_name = "lb_cons_if"
    }
    provider_connector {
      interface_name = "lb_prov_if"
    }
  }
}
```

## Argument Reference ##

* `schema_id` - (Required) The ID of the schema where the contract resides.
* `template_name` - (Required) The name of the template where the contract resides.
* `contract_name` - (Required) The name of the contract to which this service chain will be applied.
* `node_filter` - (Optional) The node filter for the service chain. Defaults to allow-all.
* `service_nodes` - (Required) A list of service nodes that form the service chain. The order of the nodes in this list defines the order in the service chain.
  * `name` - (Required) A unique name for the service node within the chain.
  * `device_type` - (Required) The type of the service device. Allowed values are firewall, loadBalancer, and other.
  * `device_ref` - (Required) The NDO UUID of the mso_service_device_cluster to be used for this node.
  * `consumer_connector` - (Required) A block that defines the consumer-side connection for the service node.
    * `interface_name` - (Required) The name of the interface on the service device cluster that will act as the consumer connector.
    * `is_redirect` - (Optional) Specifies if the connector is a redirect. Defaults to false.
  * `provider_connector` - (Required) A block that defines the provider-side connection for the service node.
    * `interface_name` - (Required) The name of the interface on the service device cluster that will act as the provider connector.
    * `is_redirect` - (Optional) Specifies if the connector is a redirect. Defaults to false.

## Attribute Reference ##

* `name` - The name of the service chain, which is derived from the contract_name.
* `id` - The unique Terraform identifier of the service chain.
* `service_nodes` - In addition to the arguments configured, the following attributes are exported for each service node
  * `index` - The computed order of the node in the service chain, starting from 0.
  * `uuid` - The NDO UUID of the service node instance within the chain.

## Importing ##

An existing MSO Schema Template Contract Service Chaining can be [imported][docs-import] into this resource via its ID/path, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_schema_template_contract_service_chaining.chain schemas/{schema_id}/templates/{template_name}/contracts/{contract_name}/serviceChaining
```
