---
layout: "mso"
page_title: "MSO: mso_tenant_policies_l3out_interface_routing_policy"
sidebar_current: "docs-mso-resource-l3out_interface_routing_policy"
description: |-
  Manages L3Out Interface Routing Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_l3out_interface_routing_policy #

Manages Layer 3 Outside (L3Out) Interface Routing Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> L3Out Interface Routing Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy" {
  template_id = mso_template.template_tenant.id
  name        = "production_routing_policy"
  description = "Production L3Out Interface Routing Policy"
  
  bfd_multi_hop_settings {
    admin_state           = "enabled"
    detection_multiplier  = 3
    min_receive_interval  = 250
    min_transmit_interval = 250
  }
  
  bfd_settings {
    admin_state           = "enabled"
    detection_multiplier  = 3
    min_receive_interval  = 50
    min_transmit_interval = 50
    echo_receive_interval = 50
    echo_admin_state      = "enabled"
    interface_control     = false
  }
  
  ospf_interface_settings {
    network_type          = "point_to_point"
    priority              = 100
    cost_of_interface     = 10
    hello_interval        = 10
    dead_interval         = 40
    retransmit_interval   = 5
    transmit_delay        = 1
    advertise_subnet      = true
    bfd                   = true
    mtu_ignore            = false
    passive_participation = false
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the L3Out Interface Routing Policy.
* `description` - (Optional) The description of the L3Out Interface Routing Policy. When unset during creation, no description is applied.
* `bfd_multi_hop_settings` - (Optional) BFD multi-hop configuration block. Omitting this block will remove BFD multi-hop settings if they exist.
  * `admin_state` - (Optional) Administrative state. Default: enabled when unset during creation. Valid values: enabled, disabled.
  * `detection_multiplier` - (Optional) The number of consecutive BFD packets that must be missed before the session is declared down. Default: 3 when unset during creation. Valid range: 1-50.
  * `min_receive_interval` - (Optional) The minimum interval in microseconds between received BFD packets. Default: 250 when unset during creation. Valid range: 250-999 microseconds.
  * `min_transmit_interval` - (Optional) The minimum interval in microseconds between transmitted BFD packets. Default: 250 when unset during creation. Valid range: 250-999 microseconds.
* `bfd_settings` - (Optional) BFD configuration block. Omitting this block will remove BFD settings if they exist.
  * `admin_state` - (Optional) Administrative state. Default: enabled when unset during creation. Valid values: enabled, disabled.
  * `detection_multiplier` - (Optional) The number of consecutive BFD packets that must be missed before the session is declared down. Default: 3 when unset during creation. Valid range: 1-50.
  * `min_receive_interval` - (Optional) The minimum interval in microseconds between received BFD packets. Default: 50 when unset during creation. Valid range: 50-999 microseconds.
  * `min_transmit_interval` - (Optional) The minimum interval in microseconds between transmitted BFD packets. Default: 50 when unset during creation. Valid range: 50-999 microseconds.
  * `echo_receive_interval` - (Optional) The minimum interval in microseconds between received BFD echo packets. Default: 50 when unset during creation. Valid range: 50-999 microseconds.
  * `echo_admin_state` - (Optional) Echo administrative state. Default: enabled when unset during creation. Valid values: enabled, disabled.
  * `interface_control` - (Optional) Interface control. Default: false (disabled) when unset during creation.
* `ospf_interface_settings` - (Optional) OSPF interface configuration block. Omitting this block will remove OSPF settings if they exist.
  * `network_type` - (Optional) The OSPF network type. Default: broadcast when unset during creation. Valid values: broadcast, point_to_point.
  * `priority` - (Optional) The OSPF router priority for designated router election. Default: 1 when unset during creation. Valid range: 0-255.
  * `cost_of_interface` - (Optional) The OSPF cost metric for the interface. Default: 0 (auto-calculated) when unset during creation. Valid range: 0-65535.
  * `hello_interval` - (Optional) The interval in seconds between OSPF hello packets. Default: 10 when unset during creation. Valid range: 1-65535 seconds.
  * `dead_interval` - (Optional) The interval in seconds before a neighbor is considered down. Default: 40 when unset during creation. Valid range: 1-65535 seconds.
  * `retransmit_interval` - (Optional) The interval in seconds between LSA retransmissions. Default: 5 when unset during creation. Valid range: 1-65535 seconds.
  * `transmit_delay` - (Optional) The estimated time in seconds to transmit an LSA. Default: 1 when unset during creation. Valid range: 1-450 seconds.
  * `advertise_subnet` - (Optional) Enable or disable subnet advertisement. Default: false (disabled) when unset during creation.
  * `bfd` - (Optional) Enable or disable BFD for OSPF. Default: false (disabled) when unset during creation.
  * `mtu_ignore` - (Optional) Enable or disable MTU mismatch detection. Default: false (disabled) when unset during creation.
  * `passive_participation` - (Optional) Enable or disable passive OSPF participation (no hello packets sent). Default: false (disabled) when unset during creation.

## Attribute Reference ##

* `uuid` - The NDO UUID of the L3Out Interface Routing Policy.
* `id` - The unique Terraform identifier of the L3Out Interface Routing Policy in the template.

## Importing ##

An existing MSO L3Out Interface Routing Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_l3out_interface_routing_policy.routing_policy templateId/{template_id}/L3OutInterfaceRoutingPolicy/{name}
```
