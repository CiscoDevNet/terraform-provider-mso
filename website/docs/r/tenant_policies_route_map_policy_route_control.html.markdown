---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_route_control"
sidebar_current: "docs-mso-resource-route_map_policy_route_control"
description: |-
  Manages Route Map Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_route_map_policy_route_control #

Manages Route Map Policies for Route Control on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v5.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy for Route Control

## Example Usage ##

```hcl
resource "mso_tenant_policies_route_map_policy_route_control" "route_map_policy" {
  template_id = mso_template.tenant_template.id
  name        = "example_route_map_policy"
  description = "Example Route Map Policy"
  match_type  = "prefix_and_routing_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template. Changing this forces a new resource.
* `name` - (Required) The name of the Route Map Policy. The name must be between `1` and `64` characters.
* `description` - (Optional) The description of the Route Map Policy. Default: empty string when unset.
* `match_type` - (Optional) The match type of the Route Map Policy. Valid values: `prefix_and_routing_policy`, `routing_policy_only`. Default: `routing_policy_only` when unset.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Route Map Policy.
* `id` - (Read-Only) The unique Terraform identifier of the Route Map Policy in the template.

## Importing ##

An existing MSO Route Map Policy can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its ID/path, using the following command:

```bash
terraform import mso_tenant_policies_route_map_policy_route_control.route_map_policy templateId/{template_id}/RouteMapPolicy/{name}
```
