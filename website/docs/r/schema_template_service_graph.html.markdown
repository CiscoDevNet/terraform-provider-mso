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
  service_node_type  = "firewall"
  description        = "hello"
  site_nodes {
    site_id     = mso_schema_site.schema_site.site_id
    tenant_name = "NkAutomation"
    node_name   = "nk-fw-2"
  }
}

```

## Argument Reference ##
* `schema_id` - (Required) Schema ID where Service Graph to be created.
* `template_name` - (Required) Template Name where Service Graph to be created.
* `service_graph_name` - (Required) Name of Service Graph.
* `service_node_type` - (Required) Type of Service Node attached to this Graph. Allowed values are `firewall`, `load-balancer` , `other`.
* `description` - (Optional) Description of Service Graph.
* `site_nodes` - (Optional) List of maps to provide Site level Node association. This maps should be provided if site is associated with template.
* `site_nodes.site_id` - (Optional) Site-Id Attached with the template. Where Service Graph will be created. This parameter is required when site is attached with the Template.
* `site_nodes.tenant_name` - (Optional) Name of Tenant holding the Service Node. This parameter is required when site is attached with the Template.
* `site_nodes.node_name` - (Optional) Name of Site level Service Node/Device Name. This parameter is required when site is attached with the Template.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the Name of Service Graph created.

## Importing ##

An existing MSO Schema Template Service Graph can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_service_graph.test_sg {schema_id}/template/{template_name}/serviceGraph/{service_graph_name}/nodeIndex/{node_index}
```
