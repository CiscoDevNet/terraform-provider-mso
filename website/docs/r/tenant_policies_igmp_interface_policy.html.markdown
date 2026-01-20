---
layout: "mso"
page_title: "MSO: mso_tenant_policies_igmp_interface_policy"
sidebar_current: "docs-mso-resource-tenant_policies_igmp_interface_policy"
description: |-
  Manages IGMP Interface Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_igmp_interface_policy #

Manages Internet Group Management Protocol (IGMP) Interface Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> IGMP Interface Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
  template_id                    = mso_template.template_tenant.id
  name                           = "test_igmp_interface_policy"
  description                    = "Test IGMP Interface Policy"
  version3_asm                   = true
  fast_leave                     = true
  report_link_local_groups       = true
  igmp_version                   = "v3"
  group_timeout                  = 300
  query_interval                 = 125
  query_response_interval        = 10
  last_member_count              = 2
  last_member_response_time      = 1
  startup_query_count            = 2
  startup_query_interval         = 31
  querier_timeout                = 255
  robustness_variable            = 2
  state_limit_route_map_uuid     = mso_tenant_policies_route_map_policy_multicast.state_limit.uuid
  report_policy_route_map_uuid   = mso_tenant_policies_route_map_policy_multicast.report_policy.uuid
  static_report_route_map_uuid   = mso_tenant_policies_route_map_policy_multicast.static_report.uuid
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the IGMP Interface Policy.
* `description` - (Optional) The description of the IGMP Interface Policy. When unset during creation, no description is applied.
* `version3_asm`- (Optional) Enable or disable IGMP version 3 Any-Source Multicast (ASM). Default: `false` (disabled) when unset during creation.
* `fast_leave` - (Optional) Enable or disable fast leave processing. When enabled, the interface immediately removes the group entry upon receiving an IGMP leave message. Default: `false` (disabled) when unset during creation.
* `report_link_local_groups` - (Optional) Enable or disable reporting of link-local multicast groups. Default: `false` (disabled) when unset during creation.
* `igmp_version` - (Optional) The IGMP protocol version to use on the interface. Default: `v2` when unset during creation. Valid values: `v2`, `v3`.
* `group_timeout` - (Optional) The time in seconds before a multicast group is removed if no IGMP reports are received. Default: 260 when unset during creation. Valid range: 3-65535 seconds.
* `query_interval` - (Optional) The interval in seconds between IGMP general query messages sent by the querier. Default: 125 when unset during creation. Valid range: 1-18000 seconds.
* `query_response_interval` - (Optional) The maximum time in seconds that hosts can wait before responding to an IGMP query. Default: 10 when unset during creation. Valid range: 1-25 seconds.
* `last_member_count` - (Optional) The number of group-specific queries sent before the router assumes there are no local members. Default: 2 when unset during creation. Valid range: 1-5.
* `last_member_response_time` - (Optional) The maximum time in seconds to wait for a response to a group-specific query before removing the group. Default: 1 when unset during creation. Valid range: 1-25 seconds.
* `startup_query_count` - (Optional) The number of queries sent at startup separated by the startup query interval. Default: 2 when unset during creation. Valid range: 1-10.
* `startup_query_interval` - (Optional) The interval in seconds between general queries sent at startup. Default: 31 when unset during creation. Valid range: 1-18000 seconds.
* `querier_timeout` - (Optional) The time in seconds before a router considers the IGMP querier to be down. Default: 255 when unset during creation. Valid range: 1-65535 seconds.
* `robustness_variable` - (Optional) Allows tuning for expected packet loss. Higher values provide more robustness but increase recovery time. Default: 2 when unset during creation. Valid range: 1-7.
* `state_limit_route_map_uuid` - (Optional) The NDO UUID of the route map policy for multicast that controls which multicast groups can be joined.
* `report_policy_route_map_uuid` - (Optional) The NDO UUID of the route map policy for multicast that filters IGMP reports.
* `static_report_route_map_uuid` - (Optional) The NDO UUID of the route map policy for multicast that defines static multicast group memberships.
* `maximum_multicast_entries` - (Optional) The maximum number of multicast route entries allowed. Default: 4294967295 when unset during creation. Valid range: 1-4294967295. Note: This parameter is only applicable when state_limit_route_map_uuid is configured.
* `reserved_multicast_entries` - (Optional) The number of multicast entries reserved and guaranteed for this policy. Default: 0 when unset during creation. Valid range: 0-4294967295.

## Attribute Reference ##

* `uuid` - The NDO UUID of the IGMP Interface Policy.
* `id` - The unique terraform identifier of the IGMP Interface Policy in the template.

## Importing ##

An existing MSO IGMP Interface Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_igmp_interface_policy.igmp_policy templateId/{template_id}/IGMPInterfacePolicy/{name}
```
