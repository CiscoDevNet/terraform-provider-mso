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
* `display_name` - (Required) Display Name of the Bridge Domain on the MSO UI.
* `vrf_name` - (Required) Name of VRF attached with Bridge Domain.
* `vrf_schema_id` - (Optional) SchemaID of VRF.
* `vrf_template_name` - (Optional) Template Name of VRF.
* `layer2_unknown_unicast` - (Optional) Type of layer 2 unknown unicast.
* `intersite_bum_traffic` - (Optional) Boolean Flag to enable or disable intersite bum traffic.
* `optimize_wan_bandwidth` - (Optional) Boolean flag to enable or disable the wan bandwidth optimization.
* `layer2_stretch` - (Optional) Boolean flag to enable or disable the layer-2 stretch.
* `layer3_multicast` - (Optional) Boolean flag to enable or disable layer 3 multicast traffic.
* `dhcp_policies` - (Optional) Block to provide dhcp_policy configurations. 
* `dhcp_policies.name` - (Optional) Dhcp_policy name. Required if you specify the dhcp_policy.
* `dhcp_policies.version` - (Optional) Version of dhcp_policy. Required if you specify the dhcp_policy.
* `dhcp_policies.dhcp_option_policy_name` - (Optional) Name of dhcp_option_policy. 
* `dhcp_policies.dhcp_option_policy_version` - (Optional) Version of dhcp_option_policy.

### Deprecation warning: do not use 'dhcp_policy' map below in combination with NDO releases 3.2 and higher, use above 'dhcp_policies' block instead.

* `dhcp_policy` - (Optional) Map to provide dhcp_policy configurations.
* `dhcp_policy.name` - (Optional) dhcp_policy name.
* `dhcp_policy.version` - (Optional) Version of dhcp_policy.
* `dhcp_policy.dhcp_option_policy_name` - (Optional) Name of dhcp_option_policy. 
* `dhcp_policy.dhcp_option_policy_version` - (Optional) Version of dhcp_option_policy. 
