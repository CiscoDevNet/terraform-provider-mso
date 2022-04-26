---
layout: "mso"
page_title: "MSO: mso_schema_site_service_graph_node"
sidebar_current: "docs-mso-data-source-schema_site_service_graph_node"
description: |-
  Data source for MSO Schema Site Level Service Graph Node
---

# mso_schema_site_service_graph_node #

Data source for MSO Schema Site Level Service Graph Node.

## Example Usage ##

```hcl
data "mso_schema_site_service_graph_node" "test" {
  schema_id          = mso_site.example.schema_id
  template_name      = mso_site.example.template_name
  service_graph_name = "sgtf"
  service_node_type  = "firewall"
  service_node_name  = "tfnode2"
}

```

## Argument Reference ##
* `schema_id`          - (Required) Schema ID holding Service Graph.
* `template_name`      - (Required) Template Name holding Service Graph. 
* `service_graph_name` - (Required) Name of Service Graph.
* `service_node_type`  - (Required) Type of Service Node to be attached to this Graph.
* `service_node_name`  - (Required) Name of the Service Graph Node.


## Attribute Reference ##

* `site_nodes`             - (Optional) List of maps to provide Site level Node association. This maps should be provided if site is associated with template.
* `site_nodes.site_id`     - (Optional) Site-Id Attached with the template. Where Service Graph is created. This parameter is required when site is attached with the Template.
* `site_nodes.tenant_name` - (Optional) Name of Tenant holding the Service Node at site level. This parameter is required when site is attached with the Template.
* `site_nodes.node_name`   - (Optional) Name of Site level Service Node/Device Name. This parameter is required when site is attached with the Template.
