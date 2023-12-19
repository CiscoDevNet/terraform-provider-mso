---
layout: "mso"
page_title: "MSO: mso_schema_site_service_graph"
sidebar_current: "docs-mso-resource-schema_site_service_graph"
description: |-
  Manages MSO Schema Site Level Service Graph
---

# mso_schema_site_service_graph #

Manages MSO Schema Site Level Service Graph.

## Example Usage ##

```hcl

resource "mso_schema_site_service_graph" "example" {
  schema_id          = mso_schema_site.schema_site_1.schema_id
  site_id            = mso_schema_site.schema_site_1.site_id
  template_name      = "template1"
  service_graph_name = "service_graph1"
  service_node {
    device_dn = data.aci_l4_l7_device.l4_l7_device_1.id
    provider_connector_type = "redir"
    consumer_connector_type = "redir"
  }
  service_node {
    device_dn = data.aci_l4_l7_device.l4_l7_device_2.id
    firewall_provider_connector_type = "snat_dnat"
    consumer_connector_type          = "redir"
    consumer_interface               = tolist(data.aci_cloud_l4_l7_third_party_device.third_party_load_balancer.interface_selectors)[0].name
    provider_interface               = tolist(data.aci_cloud_l4_l7_third_party_device.third_party_load_balancer.interface_selectors)[0].name
  }
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID under which you want to deploy Service Graph.
* `template_name` - (Required) The template name under which you want to deploy Service Graph.
* `site_id` - (Required) The site ID under which you want to deploy Service Graph.
* `service_graph_name` - (Required) The name of the Service Graph.
* `service_node` - (Required) List of service nodes attached to the Site Service Graph. Maintaining the order of the service nodes is essential.
    * `device_dn` - (Required) Dn of device associated with the service node of the Service Graph.
    * `provider_connector_type` - (Optional) Provider connector type of the service node. This parameter is only applicable for cloud sites. This parameter is only applicable for third_party_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites. Allowed values are `none`, `redir`, `snat`, `dnat` or `snat_dnat`.

        -> `snat`, `dnat` or `snat_dnat` are only supported for template_service_graph.service_node.type `firewall`.

    * `consumer_connector_type` - (Optional) Consumer connector type of the service node. This parameter is only applicable for third_party_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites. Allowed values are `redir` and `none`.
    * `provider_interface` - (Optional) Interface name of the provider interface of the service node. This parameter is only applicable for network_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites.
    * `consumer_interface` - (Optional) Interface name of the consumer interface of the service node. This parameter is only applicable for network_load_balancer and third-party firewall service nodes, when the template is attached to cloud sites.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of Service Graph created.

## Note ##
- This resource is supported only for NDO 4.1.1i and above.

- Deletion of site Service Graph is not supported by the API. Site Service Graph will be removed when site is disassociated from the template or when Service Graph is removed at the template level.

## Importing ##

An existing MSO Schema Site Service Graph can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_service_graph.example "{schema_id}/sites/{site_id}/template/{template_name}/serviceGraphs/{service_graph_name}"
```

