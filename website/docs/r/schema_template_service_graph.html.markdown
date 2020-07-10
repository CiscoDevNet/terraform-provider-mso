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
  schema_id          = "5f06a4c40f0000b63dbbd647"
  template_name      = "Template1"
  service_graph_name = "sgtf"
  service_node_type  = "firewall"
  description        = "hello"
  site_nodes {
    site_id     = "5f05c69f1900002234d0537e"
    tenant_name = "NkAutomation"
    node_name   = "nk-fw-2"
  }

}

```

## Argument Reference ##
* `schema_id` - (Required) Schema ID where Service Graph to be created.
* `template_name` - (Required) Template Name where Service Graph to be created.
* `service_graph_name` - (Required) Name of Service Graph.
* `service_node_type` - (Required) Type of Service Node attached to this Graph.
* `description` - (Optional) Description of Service Graph.
* `site_nodes` - (Optional) List of maps to provide Site level Node association. This maps should be provided if site is associated with template.
* `site_nodes.site_id` - (Optional) Site-Id Attached with the template. Where Service Graph will be created. This parameter is required when site is attached with the Template.
* `site_nodes.tenant_name` - (Optional) Name of Tenant holding the Service Node. This parameter is required when site is attached with the Template.
* `site_nodes.node_name` - (Optional) Name of Site level Service Node/Device Name. This parameter is required when site is attached with the Template.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the Name of Service Graph created.
