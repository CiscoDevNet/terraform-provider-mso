---
layout: "mso"
page_title: "MSO: mso_schema_template_bd"
sidebar_current: "docs-mso-data-source-schema_template_bd"
description: |-
  MSO Schema Template Bridge Domain Data Source.
---

# mso_schema_template_bd #

MSO Schema Template Bridge Domain Data source.

## Example Usage ##

```hcl

data "mso_schema_template_bd" "bridge_domain" {
    schema_id              = data.mso_schema.schema1.id
    template_name          = "Template1"
    name                   = "testBD"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID of Bridge Domain.
* `template_name` - (Required) Template name of Bridge Domain.
* `name` - (Required) Name of Bridge Domain.


## Attribute Reference ##
* `display_name` - (Read-Only) Display Name of the Bridge Domain on the MSO UI.
* `description` - (Read-Only) Description of the Bridge Domain on the MSO UI.
* `vrf_name` - (Read-Only) Name of the VRF attached with Bridge Domain on the MSO UI.
* `vrf_schema_id` - (Read-Only) SchemaID of the VRF.
* `vrf_template_name` - (Read-Only) Template Name of the VRF.
* `layer2_unknown_unicast` - (Read-Only) Type of the layer 2 unknown unicast.
* `intersite_bum_traffic` - (Read-Only) Enable or Disable - boolean flag of the intersite bum traffic.
* `optimize_wan_bandwidth` - (Read-Only) Enable or Disable - boolean flag of the wan bandwidth optimization.
* `layer2_stretch` - (Read-Only) Enable or Disable - boolean flag of the layer-2 stretch.
* `layer3_multicast` - (Read-Only) Enable or Disable - boolean flag of the layer 3 multicast traffic.
* `dhcp_policies` - (Read-Only) Block to provide dhcp_policy configurations. Type: Block.
  * `name` - (Read-Only) DHCP Policy name of the Bridge Domain on the MSO UI.
  * `version` - (Read-Only) DHCP Policy version of the Bridge Domain on the MSO UI.
  * `dhcp_option_policy_name` - (Read-Only) DHCP Option Policy name of the Bridge Domain on the MSO UI.
  * `dhcp_option_policy_version` - (Read-Only) DHCP Option Policy version of the Bridge Domain on the MSO UI.
* `unknown_multicast_flooding` - (Read-Only) Unknown Multicast flooding of the Bridge Domain on the MSO UI.
* `multi_destination_flooding` - (Read-Only) Multi destination flooding of the Bridge Domain on the MSO UI.
* `ipv6_unknown_multicast_flooding` - (Read-Only) IPv6 unknown multicast flooding of the Bridge Domain on the MSO UI.
* `arp_flooding` - (Read-Only) ARP flooding of the Bridge Domain on the MSO UI.
* `virtual_mac_address` - (Read-Only) Virtual Mac Address of the Bridge Domain on the MSO UI.
* `unicast_routing` - (Read-Only) Unicast Routing of the Bridge Domain on the MSO UI.

			

### Deprecation warning: do not use 'dhcp_policy' map below in combination with NDO releases 3.2 and higher, use above 'dhcp_policies' block instead.

* `dhcp_policy` - (Read-Only) Map to provide dhcp_policy configurations. Type: Block.
  * `name` - (Read-Only) DHCP Policy name of the Bridge Domain on the MSO UI.
  * `version` - (Read-Only) DHCP Policy version of the Bridge Domain on the MSO UI.
  * `dhcp_option_policy_name` - (Read-Only) DHCP Option Policy name of the Bridge Domain on the MSO UI.
  * `dhcp_option_policy_version` - (Read-Only) DHCP Option Policy version of the Bridge Domain on the MSO UI.
