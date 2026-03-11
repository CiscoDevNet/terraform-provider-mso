---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_route_control_context"
sidebar_current: "docs-mso-resource-route_map_policy_route_control_context"
description: |-
  Manages Route Map Policy Route Control Contexts on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_route_map_policy_route_control_context #

Manages Route Control Contexts within a Route Map Policy on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy -> Route Control Context

## Example Usage ##

```hcl
resource "mso_tenant_policies_route_map_policy_route_control_context" "context" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.route_map_policy.id
  name        = "example_context"
  description = "Example Route Control Context"
  order       = 1
  action      = "permit"
}
```

## Argument Reference ##

* `parent_id` - (Required) The Terraform ID of the parent Route Map Policy resource.
* `name` - (Required) The name of the Route Control Context.
* `description` - (Optional) The description of the Route Control Context.
* `order` - (Optional) The order of the Route Control Context. Valid range: 0-9. Default: `0`.
* `action` - (Optional) The action of the Route Control Context. Allowed values: `permit`, `deny`. Default: `permit`.
* `set_rule_uuid` - (Optional) The UUID of the Set Rule Policy to associate with this context.
* `match_rules` - (Optional) A set of Match Rule Policy UUIDs to associate with this context.
  * `uuid` - (Required) The UUID of the Match Rule Policy.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Route Control Context.

## Importing ##

An existing MSO Route Map Policy Route Control Context can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_route_map_policy_route_control_context.context templateId/{template_id}/RouteMapPolicy/{policy_name}/context/{context_name}
```
