---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_route_control"
sidebar_current: "docs-mso-data-source-route_map_policy_route_control"
description: |-
  Data source for Route Map Policy.
---

# mso_tenant_policies_route_map_policy_route_control #

Data source for Route Map Policies for Route Control on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_route_map_policy_route_control" "route_map_policy" {
  template_id = mso_template.template_tenant.id
  name        = "example_route_map_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the Route Map Policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Route Map Policy.
* `id` - (Read-Only) The unique Terraform identifier of the Route Map Policy in the template.
* `description` - (Read-Only) The description of the Route Map Policy.
