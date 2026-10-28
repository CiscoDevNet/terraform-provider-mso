---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_route_control_context"
sidebar_current: "docs-mso-data-source-route_map_policy_route_control_context"
description: |-
  Data source for Route Map Policy Route Control Context.
---

# mso_tenant_policies_route_map_policy_route_control_context #

Data source for Route Control Contexts within a Route Map Policy on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v5.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy for Route Control -> Route Control Context

## Example Usage ##

```hcl
data "mso_tenant_policies_route_map_policy_route_control_context" "context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.route_map_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.route_map_policy.uuid
  name                  = "example_context"
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template.
* `route_map_policy_uuid` - (Required) The UUID of the parent Route Map Policy.
* `name` - (Required) The name of the Route Control Context.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Route Control Context.
* `description` - (Read-Only) The description of the Route Control Context.
* `action` - (Read-Only) The action of the Route Control Context.
* `order` - (Read-Only) The order of the Route Control Context.
* `set_rule_uuid` - (Read-Only) The UUID of the Set Rule Policy associated with this context.
* `match_rule_uuids` - (Read-Only) A set of Match Rule Policy UUIDs associated with this context.
