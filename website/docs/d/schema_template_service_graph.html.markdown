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

* `schema_id` - (Required) The schema ID under which you want to deploy Service Graph.
* `template_name` - (Required) The template name under which you want to deploy Service Graph.
* `service_graph_name` - (Required) The name of the Service Graph.

## Attribute Reference ##

* `description` - (Read-Only) Description of Service Graph.
* `service_node` - (Read-Only) List of service nodes attached to Service Graph.
    * `type` - (Read-Only) The name of the service node type.
