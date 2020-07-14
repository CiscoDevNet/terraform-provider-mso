---
layout: "mso"
page_title: "MSO: mso_schema_template_service_graph"
sidebar_current: "docs-mso-data-source-schema_template_service_graph"
description: |-
  Data Source for MSO Schema Template Service Graph
---

# mso_schema_template_service_graph #

Data Source for MSO Schema Template Service Graph

## Example Usage ##

```hcl
data "mso_schema_template_service_graph" "test_sg" {
  schema_id          = "5f06a4c40f0000b63dbbd647"
  template_name      = "Template1"
  service_graph_name = "sgtf"
  node_index         = 1

}

```

## Argument Reference ##
* `schema_id` - (Required) Schema ID where Service Graph is created.
* `template_name` - (Required) Template Name where Service Graph is created.
* `service_graph_name` - (Required) Name of Service Graph.
* `node_index` - (Required) Integer node index of service nodes.

## Attribute Reference ##

* `service_node_type` - (Optional) Type of Service Node attached to this Graph at index provided by `node_index`.
* `site_nodes` - (Optional) List of maps to Hold Site level Node association. 
* `site_nodes.site_id` - (Optional) Site-Id Attached with the template. Where Service Graph is created. 
* `site_nodes.tenant_name` - (Optional) Name of Tenant holding the Service Node. 
* `site_nodes.node_name` - (Optional) Name of Site level Service Node/Device Name.

