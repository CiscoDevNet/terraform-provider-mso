---
layout: "mso"
page_title: "MSO: mso_tenant_policies_igmp_interface_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_igmp_interface_policy"
description: |-
  Data source for IGMP Interface Policy.
---

# mso_tenant_policies_igmp_interface_policy #

Data source for Internet Group Management Protocol (IGMP) Interface Policy. This data source is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> IGMP Interface Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_igmp_interface_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the IGMP Interface Policy to retrieve.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the IGMP Interface Policy.
* `id` - (Read-Only) The unique terraform identifier of the IGMP Interface Policy in the template.
* `description` - (Read-Only) The description of the IGMP Interface Policy.
* `version3_asm`- (Read-Only) Enable or disable IGMP version 3 Any-Source Multicast (ASM).
* `fast_leave` - (Read-Only) Enable or disable fast leave processing. When enabled, the interface immediately removes the group entry upon receiving an IGMP leave message.
* `report_link_local_groups` - (Read-Only) Enable or disable reporting of link-local multicast groups.
* `igmp_version` - (Read-Only) The IGMP protocol version to use on the interface.
* `group_timeout` - (Read-Only) The time in seconds before a multicast group is removed if no IGMP reports are received.
* `query_interval` - (Read-Only) The interval in seconds between IGMP general query messages sent by the querier.
* `query_response_interval` - (Read-Only) The maximum time in seconds that hosts can wait before responding to an IGMP query.
* `last_member_count` - (Read-Only) The number of group-specific queries sent before the router assumes there are no local members.
* `last_member_response_time` - (Read-Only) The maximum time in seconds to wait for a response to a group-specific query before removing the group.
* `startup_query_count` - (Read-Only) The number of queries sent at startup separated by the startup query interval.
* `startup_query_interval` - (Read-Only) The interval in seconds between general queries sent at startup.
* `querier_timeout` - (Read-Only) The time in seconds before a router considers the IGMP querier to be down.
* `robustness_variable` - (Read-Only) Allows tuning for expected packet loss. Higher values provide more robustness but increase recovery time.
* `state_limit_route_map_uuid` - (Read-Only) The NDO UUID of the route map policy for multicast that controls which multicast groups can be joined.
* `report_policy_route_map_uuid` - (Read-Only) The NDO UUID of the route map policy for multicast that filters IGMP reports.
* `static_report_route_map_uuid` - (Read-Only) The NDO UUID of the route map policy for multicast that defines static multicast group memberships.
* `maximum_multicast_entries` - (Read-Only) The maximum number of multicast route entries allowed.
* `reserved_multicast_entries` - (Read-Only) The number of multicast entries reserved and guaranteed for this policy.
