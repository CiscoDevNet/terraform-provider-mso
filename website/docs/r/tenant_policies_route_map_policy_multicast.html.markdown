---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_multicast"
sidebar_current: "docs-mso-resource-tenant_policies_route_map_policy_multicast"
description: |-
  Manages Route Map Policies for Multicast on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_tenant_policies_route_map_policy_multicast #

Manages Route Map Policies for Multicast on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.4(1) or higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy for Multicast

## Example Usage ##

```hcl
resource "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
  template_id = mso_template.tenant_template.id
  name        = "route_map_policy_multicast"
  description = "Example description"
  route_map_multicast_entries {
    order                   = 1
    group_ip                = "226.2.2.2/8"
    source_ip               = "1.1.1.1/1"
    rendezvous_point_ip     = "1.1.1.2"
    action                  = "permit"
  }
  route_map_multicast_entries {
    order     = 2
    group_ip  = "230.3.3.3/32"
    action    = "deny"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Tenant Policy template.
* `name` - (Required) The name of the Route Map Policy for Multicast.
* `description` - (Optional) The description of the Route Map Policy for Multicast.
* `route_map_multicast_entries` - (Optional) The list of Route Map entries for Multicast. Each entry is a rule that defines an action based on one or more matching criteria.
  * `route_map_multicast_entries.order` - (Required) The order in which the rule for an entry is evaluated.
  * `route_map_multicast_entries.group_ip` - (Optional) The Group IP address. The Group IP range must be between `224.0.0.0` and `239.255.255.255` with a netmask between `/8` and `/32`. The subnet mask must be provided.
  * `route_map_multicast_entries.source_ip` - (Optional) The Source IP address.
  * `route_map_multicast_entries.rendezvous_point_ip` - (Optional) The Rendezvous Point IP address.
  * `route_map_multicast_entries.action` - (Optional) The action defined for an entry. Allowed values are `permit`, `deny`.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Route Map Policy for Multicast.
* `id` - (Read-Only) The unique terraform identifier of the Route Map Policy for Multicast.

## Importing ##

An existing MSO Route Map Policy for Multicast can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast templateId/{template_id}/RouteMapPolicyMulticast/{name}
```
