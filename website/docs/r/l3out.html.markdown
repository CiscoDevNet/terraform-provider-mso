---
layout: "mso"
page_title: "MSO: mso_l3out"
subcategory: "L3Out"
sidebar_current: "docs-mso-resource-l3out"
description: |-
  Manages an L3Out in a Nexus Dashboard Orchestration L3Out template.
---

# mso_l3out (Resource)

Manages one IP-based L3Out in a Nexus Dashboard Orchestration L3Out template. Create an `mso_template` with `template_type = "l3out"` and an associated VRF first. Dedicated L3Out templates were introduced in Nexus Dashboard Orchestrator (NDO) 4.1(1); this page targets Nexus Dashboard 4.1+ (Orchestration 5.1+). The `mso_schema_template_l3out` resource manages the older schema-template representation.

When an optional setting is omitted, the provider reads and retains its value from Orchestration. Set a clearable string to `""` to clear it, or set a Boolean to `false` to disable it. For a protocol object, `{}` and `enabled = false` both disable the protocol.

## API Information

- **APIs**: Nexus Dashboard Orchestration API (template endpoints).
- **Create**: `PATCH /mso/api/v1/templates/{template_id}` with an `add` operation under `l3outTemplate.l3outs`.
- **Read**: `GET /mso/api/v1/templates/{template_id}` and select the L3Out by UUID or name.
- **Update**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` with changes to the selected L3Out.
- **Delete**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` with a `remove` operation.
- **Target versions**: Nexus Dashboard 4.1+ (Orchestration 5.1+).
- **Terraform ID**: `{template_id}/{l3out_uuid}`.
- **Managed resource identity**: `template_id` and `uuid`.

The provider sends template PATCH requests with `validate=false` and refreshes the selected object from the template GET response.

## GUI Information

In Nexus Dashboard, choose **Manage > Orchestration > Tenant Templates > L3Out**, open the template, then choose **Create Object > L3Out**. Select a VRF and, if needed, an L3 domain. Cisco documents the workflow in its [L3Out configuration guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html).

## Example Usage

```terraform
resource "mso_l3out" "example" {
  template_id                  = mso_template.l3out.id
  name                         = "example_l3out"
  description                  = "L3Out with BGP and OSPF"
  vrf_uuid                     = mso_schema_template_vrf.example.uuid
  l3_domain                    = mso_fabric_policies_l3_domain.example.name
  target_dscp                  = "unspecified"
  pim_enabled                  = true
  import_route_control_enabled = true
  originate_default_route      = "in_addition"

  bgp = {
    enabled = true
  }

  ospf = {
    enabled                                       = true
    area_id                                       = "0.0.0.10"
    area_type                                     = "regular"
    cost                                          = 1
    originate_default_route_always                = false
    send_redistributed_lsa                        = true
    originate_summary_lsa                         = true
    suppress_forwarding_address_in_translated_lsa = false
  }
}
```

The [complete L3Out example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/master/examples/l3out) creates the tenant, site association, VRF, L3 domain, and L3Out template.

## Schema

### Required

- `template_id` (String, API: template `templateId`, request path `{template_id}`) The UUID of the existing L3Out template.
    - **Force New**: `true`.
- `name` (String, API: `l3outTemplate.l3outs[].name`) The name of the L3Out.
- `vrf_uuid` (String, API: `l3outTemplate.l3outs[].vrfRef`) The UUID of the VRF associated with the L3Out.

### Optional

- `description` (String, API: `l3outTemplate.l3outs[].description`) The description of the L3Out.
- `l3_domain` (String, API: `l3outTemplate.l3outs[].l3domain`) The name of the L3 domain associated with the L3Out.
- `target_dscp` (String, API: `l3outTemplate.l3outs[].targetDscp`) The target DSCP for traffic leaving the L3Out.
    - **Valid values**: `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0` through `cs7`, `expedited_forwarding`, `voice_admit`, `unspecified`.
    - **Orchestration default**: `unspecified` when the field is absent on a new L3Out.
- `pim_enabled` (Boolean, API: `l3outTemplate.l3outs[].pim`) Whether PIM is enabled on the L3Out.
    - **Orchestration default**: `false` when the field is absent on a new L3Out.
- `import_route_control_enabled` (Boolean, API: `l3outTemplate.l3outs[].importRouteControl`) Whether import route control is enabled on the L3Out.
    - **Orchestration default**: `false` when the field is absent on a new L3Out.
- `originate_default_route` (String, API: `l3outTemplate.l3outs[].defaultRouteLeak.originateDefaultRoute`) The default-route origination mode for BGP or OSPF.
    - **Valid values**: `""`, `only`, `in_addition` (API value `inAddition`).
    - **Validation**: Must be explicitly set to `only` or `in_addition` when `ospf.originate_default_route_always` is true.
- `bgp` (Object, API: `l3outTemplate.l3outs[].routingProtocol`) The BGP protocol settings for the L3Out.
    - **API values**: Orchestration stores the combined protocol selection as `none`, `bgp`, `ospf`, or `bgpOspf`.
    - `enabled` (Boolean, API: `routingProtocol` contains `bgp`) Whether BGP is enabled.
        - **Provider default**: `false` when the `bgp` object is configured without `enabled`.
- `ospf` (Object, API: `l3outTemplate.l3outs[].routingProtocol` and `ospfAreaConfig`) The OSPF routing and area settings for the L3Out.
    - `enabled` (Boolean, API: `routingProtocol` contains `ospf`) Whether OSPF is enabled.
        - **Provider default**: `false` when the `ospf` object is configured without `enabled`.
        - **Validation**: `area_id` and `area_type` are required when true; other OSPF settings require true.
    - `area_id` (String, API: `ospfAreaConfig.id`) The OSPF area ID.
        - **Validation**: Required when `enabled` is true.
    - `area_type` (String, API: `ospfAreaConfig.areaType`) The OSPF area type.
        - **Valid values**: `regular`, `stub`, `nssa`.
        - **Validation**: Required when `enabled` is true.
    - `cost` (Number, API: `ospfAreaConfig.cost`) The cost assigned to the OSPF area.
        - **Default**: `1` when OSPF is enabled and no configured or prior value is available.
    - `originate_default_route_always` (Boolean, API: `defaultRouteLeak.always`) Whether OSPF advertises the default route even when it is absent from the routing table.
        - **Default**: `false` when OSPF is enabled and no configured or prior value is available.
        - **Validation**: True requires `originate_default_route` explicitly set to `only` or `in_addition`.
    - `send_redistributed_lsa` (Boolean, API: `ospfAreaConfig.control.redistribute`) Whether OSPF sends redistributed LSAs.
        - **Default**: `false` when no configured or prior value is available.
    - `originate_summary_lsa` (Boolean, API: `ospfAreaConfig.control.originate`) Whether OSPF originates summary LSAs.
        - **Default**: `false` when no configured or prior value is available.
    - `suppress_forwarding_address_in_translated_lsa` (Boolean, API: `ospfAreaConfig.control.suppressFA`) Whether OSPF suppresses the forwarding address in translated LSAs.
        - **Default**: `false` when no configured or prior value is available.
### Read-Only

- `uuid` (String, API: `l3outTemplate.l3outs[].uuid`) The UUID assigned to the L3Out by Orchestration.
- `id` (String, provider generated) The stable Terraform ID composed of the template UUID and L3Out UUID.

## Importing

Import by L3Out template UUID and L3Out UUID:

```shell
terraform import mso_l3out.example '<template_id>/<l3out_uuid>'
```

Terraform 1.12 and later can also import with the managed resource identity:

```terraform
import {
  to = mso_l3out.example
  identity = {
    template_id = "<template_id>"
    uuid        = "<l3out_uuid>"
  }
}
```

## References

### Related Terraform Objects

- [mso_l3out data source](/docs/providers/mso/d/l3out.html)
- [mso_template resource](/docs/providers/mso/r/template.html)
- [mso_schema_template_vrf resource](/docs/providers/mso/r/schema_template_vrf.html)
- [mso_fabric_policies_l3_domain resource](/docs/providers/mso/r/fabric_policies_l3_domain.html)

### External Documentation

- [Cisco Nexus Dashboard 4.1.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html)
- [Cisco Nexus Dashboard Orchestration Get Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/get-template/)
- [Cisco Nexus Dashboard Orchestration Patch Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/patch-template/)
- [Cisco Nexus Dashboard 4.1.1 API introduction](https://developer.cisco.com/docs/nexus-dashboard/4-1/)
