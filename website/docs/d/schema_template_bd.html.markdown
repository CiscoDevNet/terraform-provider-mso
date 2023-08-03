---
layout: "mso"
page_title: "MSO: mso_schema_template_bd"
sidebar_current: "docs-mso-data-source-schema_template_bd"
description: |-
  Data source for MSO Schema Template Bridge Domain (BD).
---

# mso_schema_template_bd #

Data source for MSO Schema Template Bridge Domain (BD).

## Example Usage ##

```hcl

data "mso_schema_template_bd" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  name          = "testBD"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the BD.
* `template_name` - (Required) The template name of the BD.
* `name` - (Required) The name of the BD.

## Attribute Reference ##
* `display_name` - (Read-Only) The name of the BD as displayed on the MSO UI.
* `description` - (Read-Only) The description of the BD.
* `vrf_name` - (Read-Only) The name of the VRF associated with the BD.
* `vrf_schema_id` - (Read-Only) The schema ID of the VRF associated with the BD.
* `vrf_template_name` - (Read-Only) The template name of the VRF associated with the BD.
* `layer2_unknown_unicast` - (Read-Only) The layer 2 unknown unicast type of the BD.
* `intersite_bum_traffic` - (Read-Only) Whether intersite bum traffic is enabled.
* `optimize_wan_bandwidth` - (Read-Only)  Whether wan bandwidth optimization is enabled.
* `layer2_stretch` - (Read-Only) Whether layer-2 stretch is enabled.
* `layer3_multicast` - (Read-Only) Whether layer 3 multicast traffic is enabled.
* `dhcp_policies` - (Read-Only) A list of dhcp_policy configurations for the BD.
    * `name` - (Read-Only) The DHCP policy name of the BD.
    * `version` - (Read-Only) The DHCP policy version of the BD.
    * `dhcp_option_policy_name` - (Read-Only) The DHCP option policy name of the BD.
    * `dhcp_option_policy_version` - (Read-Only) The DHCP option policy version of the BD.
* `unknown_multicast_flooding` - (Read-Only) The unknown multicast flooding settings of the BD.
* `multi_destination_flooding` - (Read-Only) The multi destination flooding settings of the BD.
* `ipv6_unknown_multicast_flooding` - (Read-Only) The IPv6 unknown multicast flooding settings of the BD.
* `arp_flooding` - (Read-Only) The ARP flooding settings of the BD.
* `virtual_mac_address` - (Read-Only) The virtual mac address of the BD.
* `unicast_routing` - (Read-Only) Whether unicast routing is enabled.

### Deprecation warning: do not use 'dhcp_policy' map below in combination with NDO releases 3.2 and higher, use above 'dhcp_policies' block instead.

* `dhcp_policy` - (Read-Only) A map to provide the dhcp_policy configuration.
    * `name` - (Read-Only) The DHCP policy name of the BD.
    * `version` - (Read-Only) The DHCP policy version of the BD.
    * `dhcp_option_policy_name` - (Read-Only) The DHCP option policy name of the BD.
    * `dhcp_option_policy_version` - (Read-Only) The DHCP option policy version of the BD.
