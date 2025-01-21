---
layout: "mso"
page_title: "MSO: mso_schema_template_bd"
sidebar_current: "docs-mso-resource-schema_template_bd"
description: |-
  Manages MSO Schema Template Bridge Domain.
---

# mso_schema_template_bd #

Manages MSO Schema Template Bridge Domain.

## Example Usage ##

```hcl

resource "mso_schema_template_bd" "bridge_domain" {
    schema_id              = mso_schema.schema1.id
    template_name          = "Template1"
    name                   = "testBD"
    display_name           = "test"
    vrf_name               = mso_schema_template_vrf.vrf1.name
    layer2_unknown_unicast = "proxy" 
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Bridge Domain.
* `template_name` - (Required) Template where Bridge Domain to be created.
* `name` - (Required) Name of the Bridge Domain.
* `display_name` - (Required) Display Name of the Bridge Domain on the MSO UI.
* `description` - (Optional) Description of the Bridge Domain on the MSO UI.
* `vrf_name` - (Required) Name of the VRF to attach with the Bridge Domain. VRF must exist in order to do this.
* `vrf_schema_id` - (Optional) SchemaID of the VRF. schema_id of the Bridge Domain will be used if not provided. This parameter is required when the VRF is in a schema which is different from the schema of the BD.
* `vrf_template_name` - (Optional) Template Name of the VRF. template_name of the Bridge Domain will be used if not provided. This parameter is required when the VRF is in a schema which is different from the schema of the BD.
* `layer2_unknown_unicast` - (Optional) Type of the layer 2 unknown unicast. Allowed values are `flood` and `proxy`. Default to `flood`.
* `intersite_bum_traffic` - (Optional) Boolean Flag to enable or disable intersite bum traffic. Default to false.
* `optimize_wan_bandwidth` - (Optional) Boolean flag to enable or disable the wan bandwidth optimization. Default to false.
* `layer2_stretch` - (Optional) Boolean flag to enable or disable the layer-2 stretch. Default to false. Should enable this flag if you want to create subnets under this Bridge Domain.
* `layer3_multicast` - (Optional) Boolean flag to enable or disable layer 3 multicast traffic. Default to false.
* `arp_flooding` - (Optional) ARP Flooding status. Default to false.
* `virtual_mac_address` - (Optional) Virtual MAC Address.
* `ep_move_detection_mode` - (Optional) Trigger an endpoint move based on incoming Gratuitous Address Resolution Protocol (GARP) packets. Allowed values are `garp` (enabled) and `none` (disabled). Default to `none`.
* `unicast_routing` - (Optional) Unicast Routing status. Default to false.
* `ipv6_unknown_multicast_flooding` - (Optional) IPv6 Unknown Multicast Flooding behavior. Allowed values are `flood` and `optimized_flooding`. Default to `flood`.
* `multi_destination_flooding` - (Optional) Multi-destination flooding behavior. Allowed values are `flood_in_bd`, `drop` and `flood_in_encap`. Default to `flood_in_bd`.
* `unknown_multicast_flooding` - (Optional) Unknown Multicast Flooding behavior. Allowed values are `flood` and `optimized_flooding`. Default to `flood`.
* `dhcp_policies` - (Optional) Block to provide dhcp_policy configurations. Type: Block.
  * `name` - (Optional) DHCP Policy name of the Bridge Domain on the MSO UI. Required if you specify the dhcp_policy.
  * `version` - (Optional) DHCP Policy version of the Bridge Domain on the MSO UI. Required if you specify the dhcp_policy.
  * `dhcp_option_policy_name` - (Optional) DHCP Option Policy name of the Bridge Domain on the MSO UI.
  * `dhcp_option_policy_version` - (Optional) DHCP Option Policy version of the Bridge Domain on the MSO UI. Required if you specify the `dhcp_option_policy_name`.

### Deprecation warning: do not use 'dhcp_policy' map below in combination with NDO releases 3.2 and higher, use above 'dhcp_policies' block instead.

* `dhcp_policy` - (Optional) Map to provide dhcp_policy configurations. Type: Block.
  * `name` - (Optional) DHCP Policy name of the Bridge Domain on the MSO UI. Required if you specify the dhcp_policy.
  * `version` - (Optional) DHCP Policy version of the Bridge Domain on the MSO UI. Required if you specify the dhcp_policy.
  * `dhcp_option_policy_name` - (Optional) DHCP Option Policy name of the Bridge Domain on the MSO UI.
  * `dhcp_option_policy_version` - (Optional) DHCP Option Policy version of the Bridge Domain on the MSO UI. Required if you specify the `dhcp_option_policy_name`.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Bridge Domain can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_bd.bridge_domain {schema_id}/template/{template_name}/bd/{name}
```