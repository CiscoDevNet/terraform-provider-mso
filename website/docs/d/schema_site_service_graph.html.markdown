---
layout: "mso"
page_title: "MSO: mso_schema_site_service_graph"
sidebar_current: "docs-mso-resource-schema_site_service_graph_node"
description: |-
  Manages MSO Schema Site Level Service Graph Node
---

# mso_schema_site_service_graph #

Manages MSO Schema Site Level Service Graph Node.

## Example Usage ##

```hcl

resource "mso_schema_site_service_graph" "test_sg" {
  schema_id          = mso_schema_site.schema_site_1.schema_id
  site_id            = mso_schema_site.schema_site_1.site_id
  template_name      = mso_schema_template_service_graph.test_sg.template_name
  service_graph_name = mso_schema_template_service_graph.test_sg.service_graph_name
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID under which you want to deploy Service Graph.
* `template_name` - (Required) The template name under which you want to deploy Service Graph.
* `site_id` - (Required) The site ID under which you want to deploy Service Graph.
* `service_graph_name` - (Required) The name of the Service Graph.


## Attribute Reference ##

* `service_node` - (Required) List of maps to provide Site level Node association.
    * `device_dn` - (Required) Dn of device associated with the service node of the Service Graph.
