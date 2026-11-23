---
layout: "mso"
page_title: "MSO: mso_l3out Data Source"
subcategory: "L3Out"
sidebar_current: "docs-mso-data-source-l3out"
description: |-
  Reads an L3Out from an NDO L3Out template.
---

# mso_l3out (Data Source)

Reads an L3Out by name from a Nexus Dashboard Orchestration (NDO) L3Out template. The template and L3Out must already exist.

- **Supported versions**: Nexus Dashboard 3.1+ (NDO 4.3+).

## API Information

- **Read**: `GET /mso/api/v1/templates/{template_id}` and select the L3Out by name.

## GUI Information

In Nexus Dashboard 4.1+ (NDO 5.1+), choose **Manage > Orchestration > Tenant Templates > L3Out**. On earlier supported versions, choose **Configure > Tenant Templates > L3Out**. Open the template and select the L3Out.

## Example Usage

```terraform
data "mso_l3out" "example" {
  template_id = mso_template.l3out.id
  name        = mso_l3out.example.name
}
```

See the [full L3Out example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/develop/examples/l3out) for details on how the L3Out relates to other Terraform resources and data sources.

The full example also creates route-map contexts and separate match-rule and set-rule policies, which require Nexus Dashboard 4.1+ (NDO 5.1+).

## Schema

### Required

- `template_id` (String) The UUID of the L3Out template to read.
    - **API**: template `templateId`, request path `{template_id}`.
- `name` (String) The name of the L3Out to select within that template.
    - **API**: `l3outTemplate.l3outs[].name`.

### Optional

None.

### Read-Only

- `id` (String, provider generated) The stable identifier composed of the template UUID and L3Out UUID.
- `uuid` (String) The UUID assigned to the L3Out by NDO.
    - **API**: `l3outTemplate.l3outs[].uuid`.
- `description` (String) The description of the L3Out.
    - **API**: `l3outTemplate.l3outs[].description`.
- `vrf_uuid` (String) The UUID of the VRF associated with the L3Out.
    - **API**: `l3outTemplate.l3outs[].vrfRef`.
- `l3_domain` (String) The name of the L3 domain associated with the L3Out.
    - **API**: `l3outTemplate.l3outs[].l3domain`.
- `target_dscp` (String) The target DSCP for traffic leaving the L3Out.
    - **API**: `l3outTemplate.l3outs[].targetDscp`.
- `pim_enabled` (Boolean) Whether PIM is enabled on the L3Out.
    - **API**: `l3outTemplate.l3outs[].pim`.
- `import_route_control_enabled` (Boolean) Whether import route control is enabled on the L3Out.
    - **API**: `l3outTemplate.l3outs[].importRouteControl`.
- `route_maps` (Object) Route-map policy references configured on the L3Out.
    - `outbound_uuid` (String) Outbound route-map policy UUID, if configured.
        - **API**: `l3outTemplate.l3outs[].exportRouteMapRef`.
    - `interleak_uuid` (String) Interleak route-map policy UUID, if configured.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.interleakRef`.
    - `static_route_redistribution_uuid` (String) Static-route redistribution route-map policy UUID, if configured.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.staticRouteRedistRef`.
    - `connected_route_redistribution_uuid` (String) Connected-route redistribution route-map policy UUID, if configured.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.connectedRouteRedistRef`.
    - `attached_host_route_redistribution_uuid` (String) Attached-host-route redistribution route-map policy UUID, if configured.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.attachedHostRouteRedistRef`.
- `originate_default_route` (String) The default-route origination mode; `""` means no leak policy exists.
    - **API**: `l3outTemplate.l3outs[].defaultRouteLeak.originateDefaultRoute`.
- `bgp` (Object) The BGP protocol settings of the L3Out.
    - `enabled` (Boolean) Whether BGP is enabled.
        - **API**: `l3outTemplate.l3outs[].routingProtocol`.
    - `route_dampening_ipv4_uuid` (String) IPv4 route-dampening route-map policy UUID, if configured.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.routeDampeningV4Ref`.
    - `route_dampening_ipv6_uuid` (String) IPv6 route-dampening route-map policy UUID, if configured.
        - **API**: `l3outTemplate.l3outs[].advancedRouteMapRefs.routeDampeningV6Ref`.
- `ospf` (Object) The OSPF routing and area settings of the L3Out.
    - `enabled` (Boolean) Whether OSPF is enabled.
        - **API**: `l3outTemplate.l3outs[].routingProtocol`.
    - `area_id` (String) The OSPF area ID.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.id`.
    - `area_type` (String) The OSPF area type.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.areaType`.
    - `cost` (Number) The cost assigned to the OSPF area.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.cost`.
    - `originate_default_route_always` (Boolean) Whether OSPF advertises the default route even when it is absent from the routing table.
        - **API**: `l3outTemplate.l3outs[].defaultRouteLeak.always`.
    - `send_redistributed_lsa` (Boolean) Whether OSPF sends redistributed LSAs.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.control.redistribute`.
    - `originate_summary_lsa` (Boolean) Whether OSPF originates summary LSAs.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.control.originate`.
    - `suppress_forwarding_address_in_translated_lsa` (Boolean) Whether OSPF suppresses forwarding addresses in translated LSAs.
        - **API**: `l3outTemplate.l3outs[].ospfAreaConfig.control.suppressFA`.
- `annotations` (Map of String) The L3Out annotations keyed by annotation key.
    - **API**: `l3outTemplate.l3outs[].tagAnnotations[]`, using `tagKey` and `tagValue`.

## References

### Related Terraform Objects

- [mso_l3out_annotation data source](/docs/providers/mso/d/l3out_annotation.html)
- [mso_l3out resource](/docs/providers/mso/r/l3out.html)

### External Documentation

- [Cisco Nexus Dashboard 4.3.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-431/external-connectivity-l3out.html)
- [NDO API reference](https://developer.cisco.com/docs/nexus-dashboard/latest/api-reference-orchestration-api-orchestration-api-overview/)
