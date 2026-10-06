---
layout: "mso"
page_title: "MSO: mso_l3out Data Source"
subcategory: "L3Out"
sidebar_current: "docs-mso-data-source-l3out"
description: |-
  Reads an L3Out from a Nexus Dashboard Orchestration L3Out template.
---

# mso_l3out (Data Source)

Reads an IP-based L3Out by name from a Nexus Dashboard Orchestration L3Out template. The template and L3Out must already exist. Dedicated L3Out templates were introduced in Nexus Dashboard Orchestrator (NDO) 4.1(1); this page targets Nexus Dashboard 4.1+ (Orchestration 5.1+).

-> Orchestration does not return BFD or OSPF interface group authentication keys. This data source exposes the authentication settings and key IDs, but not the key values.

## API Information

- **APIs**: Nexus Dashboard Orchestration API (template endpoints).
- **Read**: `GET /mso/api/v1/templates/{template_id}`; select an entry from `l3outTemplate.l3outs` by name.
- **Target versions**: Nexus Dashboard 4.1+ (Orchestration 5.1+).
- **Terraform ID**: `{template_id}/{l3out_uuid}`.

## GUI Information

In Nexus Dashboard, choose **Manage > Orchestration > Tenant Templates > L3Out**, open the template, and select the L3Out. See Cisco's [L3Out configuration guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html).

## Example Usage

```terraform
data "mso_l3out" "example" {
  template_id = mso_template.l3out.id
  name        = mso_l3out.example.name
}

output "observed_l3out_uuid" {
  value = data.mso_l3out.example.uuid
}
```

The [complete L3Out example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/master/examples/l3out) creates the referenced template and L3Out.

## Schema

### Required

- `template_id` (String, API: template `templateId`, request path `{template_id}`) The UUID of the L3Out template to read.
- `name` (String, API: `l3outTemplate.l3outs[].name`) The name of the L3Out to select within that template.

### Optional

None.

### Read-Only

- `id` (String, provider generated) The stable identifier composed of the template UUID and L3Out UUID.
- `uuid` (String, API: `l3outTemplate.l3outs[].uuid`) The UUID assigned to the L3Out by Orchestration.
- `description` (String, API: `l3outTemplate.l3outs[].description`) The description of the L3Out.
- `vrf_uuid` (String, API: `l3outTemplate.l3outs[].vrfRef`) The UUID of the VRF associated with the L3Out.
- `l3_domain` (String, API: `l3outTemplate.l3outs[].l3domain`) The name of the L3 domain associated with the L3Out.
- `target_dscp` (String, API: `l3outTemplate.l3outs[].targetDscp`) The target DSCP returned by Orchestration, translated to the provider's schema spelling.
    - **Valid values**: `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0` through `cs7`, `expedited_forwarding`, `voice_admit`, `unspecified`.
    - **Orchestration default**: `unspecified` when the field is absent on a new L3Out.
- `pim_enabled` (Boolean, API: `l3outTemplate.l3outs[].pim`) Whether PIM is enabled on the L3Out.
    - **Orchestration default**: `false` when the field is absent on a new L3Out.
- `import_route_control_enabled` (Boolean, API: `l3outTemplate.l3outs[].importRouteControl`) Whether import route control is enabled on the L3Out.
    - **Orchestration default**: `false` when the field is absent on a new L3Out.
- `originate_default_route` (String, API: `l3outTemplate.l3outs[].defaultRouteLeak.originateDefaultRoute`) The default-route origination mode; `""` means no leak policy exists.
    - **Valid values**: `""`, `only`, `in_addition` (API value `inAddition`).
- `bgp` (Object, API: `l3outTemplate.l3outs[].routingProtocol`) The BGP protocol settings of the L3Out.
    - **API values**: Orchestration stores the combined protocol selection as `none`, `bgp`, `ospf`, or `bgpOspf`.
    - `enabled` (Boolean, API: `routingProtocol` contains `bgp`) Whether BGP is enabled.
- `ospf` (Object, API: `l3outTemplate.l3outs[].routingProtocol` and `ospfAreaConfig`) The OSPF routing and area settings of the L3Out.
    - `enabled` (Boolean, API: `routingProtocol` contains `ospf`) Whether OSPF is enabled.
    - `area_id` (String, API: `ospfAreaConfig.id`) The OSPF area ID.
    - `area_type` (String, API: `ospfAreaConfig.areaType`) The OSPF area type.
        - **Valid values**: `regular`, `stub`, `nssa`.
    - `cost` (Number, API: `ospfAreaConfig.cost`) The cost assigned to the OSPF area.
    - `originate_default_route_always` (Boolean, API: `defaultRouteLeak.always`) Whether OSPF always originates the default route.
    - `send_redistributed_lsa` (Boolean, API: `ospfAreaConfig.control.redistribute`) Whether OSPF sends redistributed LSAs.
    - `originate_summary_lsa` (Boolean, API: `ospfAreaConfig.control.originate`) Whether OSPF originates summary LSAs.
    - `suppress_forwarding_address_in_translated_lsa` (Boolean, API: `ospfAreaConfig.control.suppressFA`) Whether OSPF suppresses forwarding addresses in translated LSAs.
- `annotations` (Map of String, API: `l3outTemplate.l3outs[].tagAnnotations[]` with `tagKey` and `tagValue`) The L3Out annotations indexed by key; the map is empty when none exist.
- `interface_groups` (Map of Object) The interface group policies indexed by name.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[]`.
    - `description` (String) The description of the interface group policy.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].description`.
    - `interface_routing_policy_uuid` (String) The UUID of the referenced tenant interface routing policy.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].interfaceRoutingPolicyRef`.
    - `custom_qos_policy_uuid` (String) The UUID of the referenced tenant custom QoS policy.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].qosRef`.
    - `qos_priority` (String) The QoS priority of the interface group.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].qosPriority`.
        - **Valid values**: `level1` through `level6`, `unspecified`.
        - **Orchestration default**: `unspecified` when the field is absent on a new interface group.
    - `netflow_monitor_uuids` (Map of String) The referenced NetFlow monitor UUIDs, keyed by traffic type.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].netFlowMonitorRefs`.
        - **Valid keys**: `ipv4`, `ipv6`, `ce`, `unspecified`.
    - `bfd` (Object) The single-hop BFD settings of the interface group.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd`.
        - `enabled` (Boolean) Whether single-hop BFD is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.enabled` or absence of `l3outTemplate.l3outs[].interfaceGroups[].bfd`.
        - `authentication_enabled` (Boolean) Whether single-hop BFD authentication is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.authEnabled`.
        - `key_id` (Number) The ID of the single-hop BFD authentication key.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.keyID`.
    - `bfd_multi_hop` (Object) The multi-hop BFD settings of the interface group.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop`.
        - `enabled` (Boolean) Whether multi-hop BFD is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.enabled` or absence of `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop`.
        - `authentication_enabled` (Boolean) Whether multi-hop BFD authentication is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.authEnabled`.
        - `key_id` (Number) The ID of the multi-hop BFD authentication key.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.keyID`.
    - `ospf` (Object) The OSPF authentication settings of the interface group.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf`.
        - `enabled` (Boolean) Whether OSPF authentication settings are enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.enabled` or absence of `l3outTemplate.l3outs[].interfaceGroups[].ospf`.
        - `authentication_type` (String) The OSPF authentication mode.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.authType`.
            - **Valid values**: `none`, `simple`, `md5`.
        - `key_id` (Number) The ID of the OSPF authentication key.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.keyID`.

## References

### Related Terraform Objects

- [mso_l3out resource](/docs/providers/mso/r/l3out.html)
- [mso_l3out_annotation data source](/docs/providers/mso/d/l3out_annotation.html)
- [mso_l3out_interface_group_policy data source](/docs/providers/mso/d/l3out_interface_group_policy.html)

### External Documentation

- [Cisco Nexus Dashboard 4.1.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html)
- [Cisco Nexus Dashboard Orchestration Get Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/get-template/)
- [Cisco Nexus Dashboard 4.1.1 API introduction](https://developer.cisco.com/docs/nexus-dashboard/4-1/)
