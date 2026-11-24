---
layout: "mso"
page_title: "MSO: mso_l3out"
subcategory: "L3Out"
sidebar_current: "docs-mso-resource-l3out"
description: |-
  Manages an L3Out in an NDO L3Out template.
---

# mso_l3out (Resource)

Manages one L3Out in a Nexus Dashboard Orchestration (NDO) L3Out template. Create an `mso_template` with `template_type = "l3out"` and an associated VRF first. The `mso_schema_template_l3out` resource manages the older schema-template representation.

- **Supported versions**: Nexus Dashboard 4.1+ (NDO 5.1+).

## API Information

- **Create**: `PATCH /mso/api/v1/templates/{template_id}` with an `add` operation under `l3outTemplate.l3outs`.
- **Read**: `GET /mso/api/v1/templates/{template_id}` and select the L3Out by UUID or name.
- **Update**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` with changes to the selected L3Out.
- **Delete**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` with a `remove` operation.

## GUI Information

In Nexus Dashboard, choose **Manage > Orchestration > Tenant Templates > L3Out** and open the template. Choose **Create Object > L3Out** to create an L3Out, or select an existing L3Out to view or edit its settings.

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

  route_maps = {
    outbound_uuid                           = mso_tenant_policies_route_map_policy_route_control.outbound.uuid
    interleak_uuid                          = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control.static_redistribution.uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control.connected_redistribution.uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control.attached_host_redistribution.uuid
  }

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_ipv4.uuid
    route_dampening_ipv6_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_ipv6.uuid
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

See the [full L3Out example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/develop/examples/l3out) for details on how the L3Out relates to other Terraform resources and data sources.

## Schema

### Required

- `template_id` (String) The UUID of the existing L3Out template.
    - **API**: template `templateId`, request path `{template_id}`.
    - **Force New**: `true`.
- `name` (String) The name of the L3Out.
    - **API**: `l3outTemplate.l3outs[].name`.
- `vrf_uuid` (String) The UUID of the VRF associated with the L3Out.
    - **API**: `l3outTemplate.l3outs[].vrfRef`.

### Optional

- `description` (String) The description of the L3Out. Set to `""` to clear the description.
    - **API**: `l3outTemplate.l3outs[].description`.
- `l3_domain` (String) The name of the L3 domain associated with the L3Out. Set to `""` to remove the L3 domain association.
    - **API**: `l3outTemplate.l3outs[].l3domain`.
- `target_dscp` (String) The target DSCP for traffic leaving the L3Out.
    - **API**: `l3outTemplate.l3outs[].targetDscp`.
    - **Valid values**: `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0` through `cs7`, `expedited_forwarding`, `voice_admit`, `unspecified`.
    - **NDO default**: `unspecified` when the field is absent on a new L3Out.
- `pim_enabled` (Boolean) Whether PIM is enabled on the L3Out.
    - **API**: `l3outTemplate.l3outs[].pim`.
    - **NDO default**: `false` when the field is absent on a new L3Out.
- `import_route_control_enabled` (Boolean) Whether import route control is enabled on the L3Out.
    - **API**: `l3outTemplate.l3outs[].importRouteControl`.
    - **NDO default**: `false` when the field is absent on a new L3Out.
- `originate_default_route` (String) The default-route origination mode for BGP or OSPF. Set to `""` to disable default-route origination.
    - **API**: `l3outTemplate.l3outs[].defaultRouteLeak.originateDefaultRoute`.
    - **Valid values**: `only`, `in_addition`.
    - **Validation**: Must be set to `only` or `in_addition` when `ospf.originate_default_route_always = true` is configured.
- `route_maps` (Object) References to route-map policies for L3Out route control. Each UUID is optional and computed; omit a field to preserve its current value.
    - `outbound_uuid` (String) Outbound route-map policy UUID. Set to `""` to remove the outbound route-map policy association.
        - **API**: `l3outTemplate.l3outs[].exportRouteMapRef`.
    - `interleak_uuid` (String) Interleak route-map policy UUID. Set to `""` to remove the interleak route-map policy association.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.interleakRef`.
    - `static_route_redistribution_uuid` (String) Static-route redistribution route-map policy UUID. Set to `""` to remove the static-route redistribution route-map policy association.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.staticRouteRedistRef`.
    - `connected_route_redistribution_uuid` (String) Connected-route redistribution route-map policy UUID. Set to `""` to remove the connected-route redistribution route-map policy association.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.connectedRouteRedistRef`.
    - `attached_host_route_redistribution_uuid` (String) Attached-host-route redistribution route-map policy UUID. Set to `""` to remove the attached-host-route redistribution route-map policy association.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.attachedHostRouteRedistRef`.
- `bgp` (Object) The BGP protocol settings for the L3Out.
    - **Validation**: Non-empty dampening UUIDs require `bgp.enabled = true`.
    - `enabled` (Boolean) Whether BGP is enabled.
        - **API**: `l3outTemplate.l3outs[].routingProtocol`.
        - **Provider default**: `false` when the `bgp` object is configured without `enabled`.
    - `route_dampening_ipv4_uuid` (String) IPv4 route-dampening route-map policy UUID. Set to `""` to remove the IPv4 route-dampening route-map policy association.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.routeDampeningV4Ref`.
    - `route_dampening_ipv6_uuid` (String) IPv6 route-dampening route-map policy UUID. Set to `""` to remove the IPv6 route-dampening route-map policy association.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.routeDampeningV6Ref`.
- `ospf` (Object) The OSPF routing and area settings for the L3Out.
    - `enabled` (Boolean) Whether OSPF is enabled.
        - **API**: `l3outTemplate.l3outs[].routingProtocol`.
        - **Provider default**: `false` when the `ospf` object is configured without `enabled`.
        - **Validation**: `ospf.area_id` and `ospf.area_type` are required when `ospf.enabled` is true; other OSPF settings require `ospf.enabled` to be true.
    - `area_id` (String) The OSPF area ID.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.id`.
        - **Validation**: Required when `ospf.enabled` is true.
    - `area_type` (String) The OSPF area type.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.areaType`.
        - **Valid values**: `regular`, `stub`, `nssa`.
        - **Validation**: Required when `ospf.enabled` is true.
    - `cost` (Number) The cost assigned to the OSPF area.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.cost`.
        - **NDO default**: `1` when OSPF is enabled and no configured or prior value is available.
    - `originate_default_route_always` (Boolean) Whether OSPF advertises the default route even when it is absent from the routing table.
        - **API**: `l3outTemplate.l3outs[].defaultRouteLeak.always`.
        - **Provider default**: `false` when OSPF is enabled and no configured or prior value is available.
        - **Validation**: True requires `originate_default_route` explicitly set to `only` or `in_addition`.
    - `send_redistributed_lsa` (Boolean) Whether OSPF sends redistributed LSAs.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.control.redistribute`.
        - **NDO default**: `false` when no configured or prior value is available.
    - `originate_summary_lsa` (Boolean) Whether OSPF originates summary LSAs.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.control.originate`.
        - **NDO default**: `false` when no configured or prior value is available.
    - `suppress_forwarding_address_in_translated_lsa` (Boolean) Whether OSPF suppresses the forwarding address in translated LSAs.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.control.suppressFA`.
        - **NDO default**: `false` when no configured or prior value is available.

### Read-Only

- `uuid` (String) The UUID assigned to the L3Out by NDO.
    - **API**: `l3outTemplate.l3outs[].uuid`.
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
- [mso_tenant_policies_route_map_policy_route_control resource](/docs/providers/mso/r/tenant_policies_route_map_policy_route_control.html)
- [mso_tenant_policies_route_map_policy_route_control_context resource](/docs/providers/mso/r/tenant_policies_route_map_policy_route_control_context.html)

### External Documentation

- [Cisco Nexus Dashboard 4.3.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-431/external-connectivity-l3out.html)
- [NDO API reference](https://developer.cisco.com/docs/nexus-dashboard/latest/api-reference-orchestration-api-orchestration-api-overview/)
