---
layout: "mso"
page_title: "MSO: mso_schema_site_service_graph"
sidebar_current: "docs-mso-data-source-schema_site_service_graph"
description: |-
  Data source for MSO Schema Site Level Service Graph
---

# mso_schema_site_service_graph #

Data source for MSO Schema Site Level Service Graph.

## Example Usage ##

```hcl

data "mso_schema_site_service_graph" "example" {
  schema_id          = mso_schema_site.schema_site_1.schema_id
  site_id            = mso_schema_site.schema_site_1.site_id
  template_name      = "template1"
  service_graph_name = "service_graph1"
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID under which you want to deploy Service Graph.
* `template_name` - (Required) The template name under which you want to deploy Service Graph.
* `site_id` - (Required) The site ID under which you want to deploy Service Graph.
* `service_graph_name` - (Required) The name of the Service Graph.


## Attribute Reference ##

* `service_node` - (Read-Only) List of maps to provide Site level Node association.
    * `device_dn` - (Read-Only) Dn of device associated with the service node of the Service Graph.
    * `provider_connector_type` - (Read-Only) Provider connector type of the service node. This parameter is only applicable for cloud sites. This parameter is only applicable for third_party_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites.
    * `consumer_connector_type` - (Read-Only) Consumer connector type of the service node. This parameter is only applicable for cloud sites. This parameter is only applicable for third_party_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites.
    * `provider_interface` - (Read-Only) Interface name of the provider interface of the service node. This parameter is only applicable for network_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites.
    * `consumer_interface` - (Read-Only) Interface name of the consumer interface of the service node. This parameter is only applicable for network_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites.
