---
layout: "mso"
page_title: "MSO: mso_tenant_policies_multicast_route_map_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_multicast_route_map_policy"
description: |-
  Data source for Multicast Route Policy.
---



# mso_tenant_policies_multicast_route_map_policy #

Data source for Multicast Route Policy.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy for Multicast

## Example Usage ##

```hcl
data "mso_tenant_policies_multicast_route_map_policy" "multicast_route_map_policy" {
  template_id = mso_template.tenant_template.id
  name        = "multicast_route_map_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the Multicast Route Policy policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The UUID of the Multicast Route Policy policy.
* `id` - (Read-Only) The unique identifier of the Multicast Route Policy policy in the template.
* `description` - (Read-Only) The description of the Multicast Route Policy policy.
* `multicast_route_map_entries` - (Read-Only) The list of Multicast Route Map entries.
  * `multicast_route_map_entries.order` - (Read-Only) The order in which the rule for an entry is evaluated.
  * `multicast_route_map_entries.group_ip` - (Read-Only) The Group IP address.
  * `multicast_route_map_entries.source_ip` - (Read-Only) The Source IP address.
  * `multicast_route_map_entries.rp_ip` - (Read-Only) The Rendezvous Point IP address.
  * `multicast_route_map_entries.action` - (Read-Only) The action defined for a entry.
