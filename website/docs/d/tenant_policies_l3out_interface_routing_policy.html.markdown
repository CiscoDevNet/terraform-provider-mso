---
layout: "mso"
page_title: "MSO: mso_tenant_policies_l3out_interface_routing_policy"
sidebar_current: "docs-mso-data-source-l3out_interface_routing_policy"
description: |-
  Data source for L3Out Interface Routing Policy.
---

# mso_tenant_policies_l3out_interface_routing_policy #

Data source for Layer 3 Outside (L3Out) Interface Routing Policy. This data source is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> L3Out Interface Routing Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
  template_id = mso_template.template_tenant.id
  name        = "production_routing_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the L3Out Interface Routing Policy to retrieve.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the L3Out Interface Routing Policy.
* `id` - (Read-Only) The unique Terraform identifier of the L3Out Interface Routing Policy in the template.
* `description` - (Read-Only) The description of the L3Out Interface Routing Policy.
* `bfd_multi_hop_settings` - (Read-Only) BFD multi-hop configuration block.
  * `admin_state` - (Read-Only) Administrative state.
  * `detection_multiplier` - (Read-Only) The number of consecutive BFD packets that must be missed before the session is declared down.
  * `min_receive_interval` - (Read-Only) The minimum interval in microseconds between received BFD packets.
  * `min_transmit_interval` - (Read-Only) The minimum interval in microseconds between transmitted BFD packets.
* `bfd_settings` - (Read-Only) BFD configuration block.
  * `admin_state` - (Read-Only) Administrative state.
  * `detection_multiplier` - (Read-Only) The number of consecutive BFD packets that must be missed before the session is declared down.
  * `min_receive_interval` - (Read-Only) The minimum interval in microseconds between received BFD packets. 
  * `min_transmit_interval` - (Read-Only) The minimum interval in microseconds between transmitted BFD packets. 
  * `echo_receive_interval` - (Read-Only) The minimum interval in microseconds between received BFD echo packets. 
  * `echo_admin_state` - (Read-Only) Echo administrative state.
  * `interface_control` - (Read-Only) Interface control. 
* `ospf_interface_settings` - (Read-Only) OSPF interface configuration block.
  * `network_type` - (Read-Only) The OSPF network type.
  * `priority` - (Read-Only) The OSPF router priority for designated router election.
  * `cost_of_interface` - (Read-Only) The OSPF cost metric for the interface.
  * `hello_interval` - (Read-Only) The interval in seconds between OSPF hello packets.
  * `dead_interval` - (Read-Only) The interval in seconds before a neighbor is considered down.
  * `retransmit_interval` - (Read-Only) The interval in seconds between LSA retransmissions.
  * `transmit_delay` - (Read-Only) The estimated time in seconds to transmit an LSA.
  * `advertise_subnet` - (Read-Only) Enable or disable subnet advertisement. 
  * `bfd` - (Read-Only) Enable or disable BFD for OSPF. 
  * `mtu_ignore` - (Read-Only) Enable or disable MTU mismatch detection. 
  * `passive_participation` - (Read-Only) Enable or disable passive OSPF participation (no hello packets sent).
