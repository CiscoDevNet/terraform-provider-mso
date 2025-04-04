---
layout: "mso"
page_title: "MSO: mso_tenant_policies_route_map_policy_multicast"
sidebar_current: "docs-mso-data-source-tenant_policies_route_map_policy_multicast"
description: |-
  Data source for Route Map Policy for Multicast.
---



# mso_tenant_policies_route_map_policy_multicast #

Data source for Route Map Policy for Multicast. This data source is supported in NDO v4.4(1) or higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy for Multicast

## Example Usage ##

```hcl
data "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
  template_id = mso_template.tenant_template.id
  name        = "route_map_policy_multicast"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Tenant Policy template.
* `name` - (Required) The name of the Route Map Policy for Multicast.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Route Map Policy for Multicast.
* `id` - The unique terraform identifier of the Route Map Policy for Multicast.
* `description` - (Read-Only) The description of the Route Map Policy for Multicast.
* `route_map_entries_multicast` - (Read-Only) The list of Route Map entries for Multicast.
  * `route_map_entries_multicast.order` - (Read-Only) The order in which the rule for an entry is evaluated.
  * `route_map_entries_multicast.group_ip` - (Read-Only) The Group IP address.
  * `route_map_entries_multicast.source_ip` - (Read-Only) The Source IP address.
  * `route_map_entries_multicast.rendezvous_point_ip` - (Read-Only) The Rendezvous Point IP address.
  * `route_map_entries_multicast.action` - (Read-Only) The action defined for a entry.
