---
layout: "mso"
page_title: "MSO: mso_schema_site_service_graph_node"
sidebar_current: "docs-mso-resource-schema_site_service_graph_node"
description: |-
  Manages MSO Schema Site Level Service Graph Node
---

# mso_schema_site_service_graph_node #

Manages MSO Schema Site Level Service Graph Node.

## Example Usage ##

```hcl

resource "mso_schema_site_service_graph_node" "test_sg" {
  schema_id          = mso_schema.schema1.id
  template_name      = "Template1"
  service_graph_name = "sgtf"
  service_node_type  = "firewall"
  site_nodes {
    site_id     = mso_site.site1.id
    tenant_name = "NkAutomation"
    node_name   = "nk-fw-2"
  }
}

```

## Argument Reference ##
* `schema_id` - (Required) Schema ID holding Service Graph.
* `template_name` - (Required) Template Name holding Service Graph. 
* `service_graph_name` - (Required) Name of Service Graph.
* `service_node_type` - (Required) Type of Service Node to be attached to this Graph.
* `site_nodes` - (Optional) List of maps to provide Site level Node association. This maps should be provided if site is associated with template.
* `site_nodes.site_id` - (Optional) Site-Id Attached with the template. Where Service Graph is created. This parameter is required when site is attached with the Template.
* `site_nodes.tenant_name` - (Optional) Name of Tenant holding the Service Node at site level. This parameter is required when site is attached with the Template.
* `site_nodes.node_name` - (Optional) Name of Site level Service Node/Device Name. This parameter is required when site is attached with the Template.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the node name of Service Node created.
