---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_route_control"
sidebar_current: "docs-mso-resource-route_map_policy_route_control"
description: |-
  Manages Route Map Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_route_map_policy_route_control #

Manages Route Map Policies for Route Control on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_route_map_policy_route_control" "route_map_policy" {
  template_id = mso_template.template_tenant.id
  name        = "example_route_map_policy"
  description = "Example Route Map Policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the Route Map Policy.
* `description` - (Optional) The description of the Route Map Policy. When unset during creation, no description is applied.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Route Map Policy.
* `id` - (Read-Only) The unique Terraform identifier of the Route Map Policy in the template.

## Importing ##

An existing MSO Route Map Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_route_map_policy_route_control.route_map_policy templateId/{template_id}/RouteMapPolicy/{name}
```
