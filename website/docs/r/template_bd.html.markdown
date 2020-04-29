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
    schema_id              = "5ea809672c00003bc40a2799"
    template_name          = "Template1"
    name                   = "testBD"
    display_name           = "test"
    vrf_name               = "demo"
    layer2_unknown_unicast = "proxy" 
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Bridge Domain.
* `template_name` - (Required) Template where Bridge Domain to be created.
* `name` - (Required) Name of Bridge Domain.
* `display_name` - (Required) Display Name of the Bridge Domain on the MSO UI.
* `vrf_name` - (Required) Name of VRF to attach with Bridge Domain. VRF must exist.
* `vrf_schema_id` - (Optional) SchemaID of VRF. schema_id of Bridge Domain will be used if not provided. Should use this parameter when VRF is in different schema than BD.
* `vrf_template_name` - (Optional) Template Name of VRF. template_name of Bridge Domain will be used if not provided. Should use this parameter when VRF is in different schema than BD.
* `layer2_unknown_unicast` - (Optional) Type of layer 2 unknown unicast. Allowed values are `flood` and `proxy`. Default to `flood`.
* `intersite_bum_traffic` - (Optional) Boolean Flag to enable or disable intersite bum traffic. Default to false.
* `optimize_wan_bandwidth` - (Optional) Boolean flag to enable or disable the wan bandwidth optimization. Default to false.
* `layer2_stretch` - (Optional) Boolean flag to enable or disable the layer-2 stretch. Default to false. Should enable this flag if you want to create subnets under this Bridge Domain.
* `layer3_multicast` - (Optional) Boolean flag to enable or disable layer 3 multicast traffic. Default to false.
* `dhcp_policy` - (Optional) Map to provide dhcp_policy configurations. 
* `dhcp_policy.name` - (Optional) dhcp_policy name. Required if you specify the dhcp_policy.
* `dhcp_policy.version` - (Optional) Version of dhcp_policy. Required if you specify the dhcp_policy.
* `dhcp_policy.dhcp_option_policy_name` - (Optional) Name of dhcp_option_policy. 
* `dhcp_policy.dhcp_option_policy_version` - (Optional) Version of dhcp_option_policy. Require if you specify the `dhcp_policy.dhcp_option_policy_name`.


## Attribute Reference ##

No attributes are exported.