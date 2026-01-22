---
layout: "mso"
page_title: "MSO: mso_tenant_policies_l3out_node_routing_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_l3out_node_routing_policy"
description: |-
  Data source for L3Out Node Routing Policy.
---

# mso_tenant_policies_l3out_node_routing_policy #

Data source for Layer 3 Outside (L3Out) Node Routing Policy.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> L3Out Node Routing Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_l3out_node_routing_policy" "node_policy" {
  template_id = mso_template.template_tenant.id
  name        = "production_node_routing_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the L3Out Node Routing Policy to retrieve.

## Attribute Reference ##

* `uuid` - (Read-Only) The UUID of the L3Out Node Routing Policy.
* `id` - (Read-Only) The unique Terraform identifier of the L3Out Node Routing Policy in the template.
* `description` - (Read-Only) The description of the L3Out Node Routing Policy.
* `as_path_multipath_relax` - (Read-Only) BGP Best Path Control AS path multipath relax setting.
* `bfd_multi_hop_settings` - (Read-Only) A list containing BFD multi-hop configuration. Empty list if not configured. When present, contains a single element with:
  * `admin_state` - (Read-Only) Administrative state.
  * `detection_multiplier` - (Read-Only) The number of consecutive BFD packets that must be missed before the session is declared down.
  * `min_receive_interval` - (Read-Only) The minimum interval in microseconds between received BFD packets.
  * `min_transmit_interval` - (Read-Only) The minimum interval in microseconds between transmitted BFD packets.
* `bgp_node_settings` - (Read-Only) A list containing BGP node configuration. Empty list if not configured. When present, contains a single element with:
  * `graceful_restart_helper` - (Read-Only) BGP graceful restart helper mode.
  * `keep_alive_interval` - (Read-Only) The BGP keepalive interval in seconds.
  * `hold_interval` - (Read-Only) The BGP hold interval in seconds.
  * `stale_interval` - (Read-Only) The BGP stale interval in seconds for graceful restart.
  * `max_as_limit` - (Read-Only) Maximum AS path limit to prevent routing loops.
