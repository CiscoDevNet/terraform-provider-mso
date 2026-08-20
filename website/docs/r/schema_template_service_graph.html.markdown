---
layout: "mso"
page_title: "MSO: mso_schema_template_service_graph"
sidebar_current: "docs-mso-resource-schema_template_service_graph"
description: |-
  Manages MSO Schema Template Service Graph
---

# mso_schema_template_service_graph #

Manages MSO Schema Template Service Graph

## Example Usage ##

```hcl

resource "mso_schema_template_service_graph" "test_sg" {
  schema_id          = mso_schema.schema1.id
  template_name      = "Template1"
  service_graph_name = "sgtf"
  description        = "Created by terraform"

  service_node {
    type = "firewall"
  }
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID under which you want to deploy Service Graph.
* `template_name` - (Required) The template name under which you want to deploy Service Graph.
* `service_graph_name` - (Required) Name of the Service Graph.
* `description` - (Optional) Description of Service Graph.
* `service_node` - (Required) List of service nodes attached to the Service Graph. At least one service node is required.
    * `type` - (Required) The name of the service node type. Built-in values are `firewall`, `load-balancer` and `other`. Custom node types created via `mso_service_node_type` are also supported on Nexus Dashboard (ND) 4.2 / NDO 5.2 and earlier.

## Note ##

* The NDO API does not support removing service nodes from an existing Service Graph. To reduce the number of nodes, the Service Graph must be deleted and recreated.

* As of Nexus Dashboard (ND) 4.3 / NDO 5.3, the platform no longer allows creating custom service node types; only the built-in `firewall`, `load-balancer` and `other` types can be used as `service_node.type`. See the `mso_service_node_type` resource documentation for details.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of Service Graph created.

## Importing ##

An existing MSO Schema Template Service Graph can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_service_graph.test_sg {schema_id}/templates/{template_name}/serviceGraphs/{service_graph_name}
```
