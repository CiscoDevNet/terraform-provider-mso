---
layout: "mso"
page_title: "MSO: mso_l3out_interface_group_policy Data Source"
subcategory: "L3Out"
sidebar_current: "docs-mso-data-source-l3out_interface_group_policy"
description: |-
  Reads one interface group policy from an L3Out in Nexus Dashboard Orchestration.
---

# mso_l3out_interface_group_policy (Data Source)

Reads one interface group policy by name from an L3Out in a Nexus Dashboard Orchestration L3Out template. The template, L3Out, and group must already exist.

-> Orchestration does not return BFD or OSPF authentication keys. This data source exposes the authentication settings and key IDs, but not the key values.

## API Information

- **APIs**: Nexus Dashboard Orchestration API (template endpoints).
- **Read**: `GET /mso/api/v1/templates/{template_id}`; select the L3Out by UUID and group by name.
- **Target versions**: Nexus Dashboard 4.1+ (Orchestration 5.1+).
- **Terraform ID**: `{template_id}/{l3out_uuid}/{name}`.

## GUI Information

In Nexus Dashboard, choose **Manage > Orchestration > Tenant Templates > L3Out** and open the template. Interface groups appear under the L3Out's **Node/Interface Group Policy** objects; choose **Interface** for that policy type. Cisco describes the [interface group workflow](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html).

## Example Usage

```terraform
data "mso_l3out_interface_group_policy" "edge" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  name        = mso_l3out_interface_group_policy.edge.name
}

output "edge_ipv4_monitor_uuid" {
  value = data.mso_l3out_interface_group_policy.edge.netflow_monitor_uuids["ipv4"]
}
```

The [complete standalone interface group example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/master/examples/l3out_interface_group_policy) includes the parent L3Out and all referenced tenant policies.

## Schema

### Required

- `template_id` (String) The UUID of the L3Out template to read.
    - **API**: template `templateId`, request path `{template_id}`.
- `l3out_uuid` (String) The UUID of the parent L3Out.
    - **API**: `l3outTemplate.l3outs[].uuid`.
- `name` (String) The name of the interface group policy to select within the parent L3Out.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].name`.

### Optional

None.

### Read-Only

- `id` (String, provider generated) The stable identifier composed of the template UUID, L3Out UUID, and group name.
- `description` (String) The description of the interface group policy.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].description`.
- `interface_routing_policy_uuid` (String) The UUID of the referenced tenant interface routing policy.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].interfaceRoutingPolicyRef`.
- `custom_qos_policy_uuid` (String) The UUID of the referenced tenant custom QoS policy.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].qosRef`.
- `qos_priority` (String) The QoS priority of the interface group.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].qosPriority`.
    - **Valid values**: `level1` through `level6`, `unspecified`.
- `netflow_monitor_uuids` (Map of String) The referenced NetFlow monitor UUIDs, indexed by traffic type.
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

- [mso_l3out_interface_group_policy resource](/docs/providers/mso/r/l3out_interface_group_policy.html)
- [mso_l3out data source](/docs/providers/mso/d/l3out.html)

### External Documentation

- [Cisco Nexus Dashboard 4.1.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html)
- [Cisco Nexus Dashboard Orchestration Get Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/get-template/)
- [Cisco Nexus Dashboard 4.1.1 API introduction](https://developer.cisco.com/docs/nexus-dashboard/4-1/)
