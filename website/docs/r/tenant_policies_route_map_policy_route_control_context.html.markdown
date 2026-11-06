---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_route_control_context"
sidebar_current: "docs-mso-resource-route_map_policy_route_control_context"
description: |-
  Manages Route Map Policy Route Control Contexts on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_route_map_policy_route_control_context #

Manages Route Control Contexts within a Route Map Policy on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v5.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy for Route Control -> Route Control Context

## Example Usage ##

```hcl
resource "mso_tenant_policies_route_map_policy_route_control_context" "context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.route_map_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.route_map_policy.uuid
  name                  = "example_context"
  description           = "Example Route Control Context"
  order                 = 1
  action                = "permit"
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.set_rule_policy.uuid
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.match_rule_policy.uuid]
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template. Changing this forces a new resource.
* `route_map_policy_uuid` - (Required) The UUID of the parent Route Map Policy. Changing this forces a new resource.
* `name` - (Required) The name of the Route Control Context. The name must be between `1` and `64` characters.
* `description` - (Optional) The description of the Route Control Context. Default: empty string when unset.
* `order` - (Optional) The order of the Route Control Context. Range: 0-9. Default: `0` when unset.
* `action` - (Optional) The action of the Route Control Context. Valid values: `permit`, `deny`. Default: `permit` when unset.
* `set_rule_uuid` - (Optional) The UUID of the Set Rule Policy to associate with this context. Omit or set to `null` to remove the reference.
* `match_rule_uuids` - (Optional) A set of Match Rule Policy UUIDs to associate with this context. Omit or set to `[]` to remove the references.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Route Control Context.

## Importing ##

An existing MSO Route Map Policy Route Control Context can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its ID/path, using the following command:

```bash
terraform import mso_tenant_policies_route_map_policy_route_control_context.context templateId/{template_id}/RouteMapPolicy/{route_map_policy_uuid}/context/{context_name}
```
