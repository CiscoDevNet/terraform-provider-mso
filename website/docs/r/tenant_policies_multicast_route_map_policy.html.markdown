---
layout: "mso"
page_title: "MSO: mso_tenant_policies_multicast_route_map_policy"
sidebar_current: "docs-mso-resource-tenant_policies_multicast_route_map_policy"
description: |-
  Manages Multicast Route Map Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_tenant_policies_multicast_route_map_policy #

Manages Route Map Policies for Multicast on Cisco Nexus Dashboard Orchestrator (NDO).

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Route Map Policy for Multicast

## Example Usage ##

```hcl
resource "mso_tenant_policies_multicast_route_map_policy" "multicast_route_map_policy" {
  template_id = mso_template.tenant_template.id
  name        = "multicast_route_map_policy"
  description = "Example description"
  multicast_route_map_entries {
    order     = 1
    group_ip  = "226.2.2.2/8"
    source_ip = "1.1.1.1/1"
    rp_ip     = "1.1.1.2"
    action    = "permit"
  }
  multicast_route_map_entries {
    order     = 2
    group_ip  = "230.3.3.3/32"
    action    = "deny"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Tenant Policy template.
* `name` - (Required) The name of the Multicast Route Map Policy.
* `description` - (Optional) The description of the Multicast Route Map Policy.
* `multicast_route_map_entries` - (Optional) The list of Multicast Route Map entries. Each entry is a rule that defines an action based on one or more matching criteria.
  * `multicast_route_map_entries.order` - (Required) The order in which the rule for an entry is evaluated.
  * `multicast_route_map_entries.group_ip` - (Optional) The Group IP address. The Group IP range must be between `224.0.0.0` and `239.255.255.255` with a netmask between `/8` and `/32`. The subnet mask must be provided.
  * `multicast_route_map_entries.source_ip` - (Optional) The Source IP address.
  * `multicast_route_map_entries.rp_ip` - (Optional) The Rendezvous Point IP address.
  * `multicast_route_map_entries.action` - (Optional) The action defined for a entry. Allowed values are `permit`, `deny`.

## Attribute Reference ##

* `uuid` - The UUID of the Multicast Route Map Policy.

## Importing ##

An existing MSO Multicast Route Map Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_tenant_policies_multicast_route_map_policy.multicast_route_map_policy templateId/{template_id}/MulticastRouteMapPolicy/{name}
```
