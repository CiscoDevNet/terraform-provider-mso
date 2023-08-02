---
layout: "mso"
page_title: "MSO: mso_schema_template_service_graph"
sidebar_current: "docs-mso-data-source-schema_template_service_graph"
description: |-
  Data Source for MSO Schema Template Service Graph.
---

# mso_schema_template_service_graph #

Data Source for MSO Schema Template Service Graph.

## Example Usage ##

```hcl

data "mso_schema_template_service_graph" "example" {
  schema_id          = data.mso_schema.schema1.id
  template_name      = "Template1"
  service_graph_name = "sgtf"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Service Graph.
* `template_name` - (Required) The template name of the Service Graph.
* `service_graph_name` - (Required) The name of the Service Graph.

## Attribute Reference ##

* `site_nodes` - (Read-Only) A list of site nodes for the Service Graph.
    * `node_index` - (Read-Only) The index of the Service Node.
    * `service_node_type` - (Read-Only) The type of the Service Node.
    * `site_id` - (Read-Only) The site ID of the Service Node.
    * `tenant_name` - (Read-Only) The tenant name of the Service Node.
    * `node_name` - (Read-Only) The name of the site level Service Node.
