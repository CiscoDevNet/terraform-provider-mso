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
  service_node {
    type = "firewall"
  }
  description        = "Created by terraform"
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID under which you want to deploy Service Graph.
* `template_name` - (Required) The template name under which you want to deploy Service Graph.
* `service_graph_name` - (Required) Name of the Service Graph.
* `service_node_type` - (Optional) **Deprecated**. Type of Service Node attached to this Graph. Allowed values are `firewall`, `load-balancer` and `other`.
* `description` - (Optional) Description of Service Graph.
* `service_node` - (Required) List of service nodes attached to Service Graph.
    * `service_node.type` - (Required) Type of Service Node attached to the Service Graph. Allowed values are `firewall`, `load-balancer` and `other`.


## NOTE ##
The `site_nodes` parameters are removed from Template level Service Graph resource.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of Service Graph created.

## Importing ##

An existing MSO Schema Template Service Graph can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_service_graph.test_sg {schema_id}/template/{template_name}/serviceGraph/{service_graph_name}/nodeIndex/{node_index}
```
