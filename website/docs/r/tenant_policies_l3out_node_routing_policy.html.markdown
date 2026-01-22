---
layout: "mso"
page_title: "MSO: mso_tenant_policies_l3out_node_routing_policy"
sidebar_current: "docs-mso-resource-tenant_policies_l3out_node_routing_policy"
description: |-
  Manages L3Out Node Routing Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_l3out_node_routing_policy #

Manages Layer 3 Outside (L3Out) Node Routing Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> L3Out Node Routing Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
  template_id             = mso_template.template_tenant.id
  name                    = "production_node_routing_policy"
  description             = "Production L3Out Node Routing Policy"
  as_path_multipath_relax = true
  
  bfd_multi_hop_settings {
    admin_state           = "enabled"
    detection_multiplier  = 3
    min_receive_interval  = 250
    min_transmit_interval = 250
  }
  
  bgp_node_settings {
    graceful_restart_helper = true
    keep_alive_interval     = 60
    hold_interval           = 180
    stale_interval          = 300
    max_as_limit            = 0
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the L3Out Node Routing Policy.
* `description` - (Optional) The description of the L3Out Node Routing Policy. When unset during creation, no description is applied.
* `as_path_multipath_relax` - (Optional) BGP Best Path Control - enables AS path multipath relaxation to allow load balancing across paths with different AS paths. When unset during creation, this setting is not configured.
* `bfd_multi_hop_settings` - (Optional) BFD multi-hop configuration block. Omitting this block will remove BFD multi-hop settings if they exist.
* `admin_state` - (Optional) Administrative state. Default: `enabled` when unset during creation. Valid values: `enabled`, `disabled`.
* `detection_multiplier` - (Optional) The number of consecutive BFD packets that must be missed before the session is declared down. Default: 3 when unset during creation. Valid range: 1-50.
* `min_receive_interval` - (Optional) The minimum interval in microseconds between received BFD packets. Default: 250 when unset during creation. Valid range: 250-999 microseconds.
* `min_transmit_interval` - (Optional) The minimum interval in microseconds between transmitted BFD packets. Default: 250 when unset during creation. Valid range: 250-999 microseconds.
* `bgp_node_settings` - (Optional) BGP node configuration block. Omitting this block will remove BGP node settings if they exist.
* `graceful_restart_helper` - (Optional) Enable or disable BGP graceful restart helper mode, allowing the router to assist peers during restart. Default: `true` (enabled) when unset during creation.
* `keep_alive_interval` - (Optional) The BGP keepalive interval in seconds. Keepalive messages maintain the BGP session. Default: 60 when unset during creation. Valid range: 0-3600 seconds.
* `hold_interval` - (Optional) The BGP hold interval in seconds. If no message is received within this time, the session is terminated. Default: 180 when unset during creation. Must be 0 or between 3-3600 seconds.
* `stale_interval` - (Optional) The BGP stale interval in seconds for graceful restart. Routes are marked stale after this period. Default: 300 when unset during creation. Valid range: 1-3600 seconds.
* `max_as_limit` - (Optional) Maximum AS path limit to prevent routing loops. A value of 0 means no limit. Default: 0 (no limit) when unset during creation. Valid range: 0-2000.

## Attribute Reference ##

* `uuid` - The NDO UUID of the L3Out Node Routing Policy.
* `id` - The unique Terraform identifier of the L3Out Node Routing Policy in the template.

## Importing ##

An existing MSO L3Out Node Routing Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html

```bash
terraform import mso_tenant_policies_l3out_node_routing_policy.node_policy templateId/{template_id}/L3OutNodeRoutingPolicy/{name}
```
