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

When an optional setting is omitted, the provider reads and retains its value from Orchestration. Set a clearable string to `""` to clear it, or set a Boolean to `false` to disable it. For a protocol object, `{}` and `enabled = false` both disable the protocol. Collection attributes have their own ownership rules below.

-> **Authentication key limitation:** Orchestration accepts the BFD and OSPF interface group keys but returns only an opaque key reference. The provider retains a previously configured key in Terraform state so that refresh and unrelated updates do not discard it. It cannot verify the stored key against Orchestration or detect a key changed outside Terraform. An imported L3Out has no key value to recover. Treat Terraform state as sensitive because it contains configured keys.

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
  description                  = "L3Out managed with annotations and interface groups"
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

  annotations = {
    owner   = "network"
    purpose = "external-routing"
  }

  interface_groups = {
    edge = {
      description                   = "Edge interfaces"
      interface_routing_policy_uuid = mso_tenant_policies_l3out_interface_routing_policy.example.uuid
      custom_qos_policy_uuid        = mso_tenant_policies_custom_qos_policy.example.uuid
      qos_priority                  = "level6"

      netflow_monitor_uuids = {
        ipv4        = mso_tenant_policies_netflow_monitor.ipv4.uuid
        ipv6        = mso_tenant_policies_netflow_monitor.ipv6.uuid
        ce          = mso_tenant_policies_netflow_monitor.ce.uuid
        unspecified = mso_tenant_policies_netflow_monitor.unspecified.uuid
      }

      bfd = {
        enabled                = true
        authentication_enabled = true
        key_id                 = 20
        key                    = var.bfd_key
      }

      bfd_multi_hop = {
        enabled                = true
        authentication_enabled = true
        key_id                 = 30
        key                    = var.bfd_multi_hop_key
      }

      ospf = {
        enabled             = true
        authentication_type = "md5"
        key_id              = 10
        key                 = var.ospf_key
      }
    }
  }
}
```

The [complete L3Out example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/master/examples/l3out) creates the tenant, site association, VRF, L3 domain, policy templates, QoS and routing policies, NetFlow exporter, record, and all four monitor references.

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
- `annotations` (Map of String, API: `l3outTemplate.l3outs[].tagAnnotations[]` with `tagKey` and `tagValue`) The L3Out annotations indexed by key. When configured, the map owns the entire collection; removing a key removes its annotation, and `{}` clears all annotations. To manage individual annotations, use the [mso_l3out_annotation](/docs/providers/mso/r/l3out_annotation.html) resource and leave this map omitted.
    - **Validation**: Keys must be nonempty and values cannot be null. Different keys may have the same value.
- `interface_groups` (Map of Object) The interface group policies indexed by their unique names. When configured, the map owns the entire collection; removing a name removes that group, and `{}` clears all groups. To manage individual groups, use the [mso_l3out_interface_group_policy](/docs/providers/mso/r/l3out_interface_group_policy.html) resource and leave this map omitted.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[]`.
    - **Validation**: Group names must be nonempty. Use an empty object, such as `edge = {}`, for a named group without explicit settings; `edge = null` is invalid.
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
    - `netflow_monitor_uuids` (Map of String) The referenced NetFlow monitor UUIDs, keyed by traffic type. A configured map owns all references: remove a key to remove that reference, or use `{}` to remove all references.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].netFlowMonitorRefs`.
        - **Valid keys**: `ipv4`, `ipv6`, `ce`, `unspecified`.
        - **Validation**: Each value must be a nonempty string; `""` does not remove a reference.
    - `bfd` (Object) The single-hop BFD settings of the interface group.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd`.
        - `enabled` (Boolean) Whether single-hop BFD is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.enabled` or absence of `l3outTemplate.l3outs[].interfaceGroups[].bfd`.
            - **Default**: `false`.
        - `authentication_enabled` (Boolean) Whether single-hop BFD authentication is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.authEnabled`.
            - **Validation**: True requires `key_id` and `key`; authentication settings require `enabled = true`.
        - `key_id` (Number) The ID of the single-hop BFD authentication key.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.keyID`.
            - **Validation**: Allowed only when authentication is enabled.
        - `key` (String, Sensitive) The single-hop BFD authentication key. Orchestration does not return the secret; the provider retains a previously configured key in state.
            - **API write**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.key.value`.
            - **Validation**: Allowed only when authentication is enabled.
    - `bfd_multi_hop` (Object) The multi-hop BFD settings of the interface group.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop`.
        - `enabled` (Boolean) Whether multi-hop BFD is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.enabled` or absence of `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop`.
            - **Default**: `false`.
        - `authentication_enabled` (Boolean) Whether multi-hop BFD authentication is enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.authEnabled`.
            - **Validation**: True requires `key_id` and `key`; authentication settings require `enabled = true`.
        - `key_id` (Number) The ID of the multi-hop BFD authentication key.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.keyID`.
            - **Validation**: Allowed only when authentication is enabled.
        - `key` (String, Sensitive) The multi-hop BFD authentication key. Orchestration does not return the secret; the provider retains a previously configured key in state.
            - **API write**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.key.value`.
            - **Validation**: Allowed only when authentication is enabled.
    - `ospf` (Object) The OSPF authentication settings of the interface group.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf`.
        - `enabled` (Boolean) Whether OSPF authentication settings are enabled.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.enabled` or absence of `l3outTemplate.l3outs[].interfaceGroups[].ospf`.
            - **Default**: `false`.
        - `authentication_type` (String) The OSPF authentication mode.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.authType`.
            - **Valid values**: `none`, `simple`, `md5`.
            - **Validation**: `simple` and `md5` require `key_id` and `key`; other OSPF settings require `enabled = true`.
        - `key_id` (Number) The ID of the OSPF authentication key.
            - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.keyID`.
            - **Validation**: Allowed only with `simple` or `md5`.
        - `key` (String, Sensitive) The OSPF authentication key. Orchestration does not return the secret; the provider retains a previously configured key in state.
            - **API write**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.key.value`.
            - **Validation**: Allowed only with `simple` or `md5`.

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
- [mso_l3out_annotation resource](/docs/providers/mso/r/l3out_annotation.html)
- [mso_l3out_interface_group_policy resource](/docs/providers/mso/r/l3out_interface_group_policy.html)
- [mso_template resource](/docs/providers/mso/r/template.html)
- [mso_schema_template_vrf resource](/docs/providers/mso/r/schema_template_vrf.html)
- [mso_fabric_policies_l3_domain resource](/docs/providers/mso/r/fabric_policies_l3_domain.html)

### External Documentation

- [Cisco Nexus Dashboard 4.1.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html)
- [Cisco Nexus Dashboard Orchestration Get Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/get-template/)
- [Cisco Nexus Dashboard Orchestration Patch Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/patch-template/)
- [Cisco Nexus Dashboard 4.1.1 API introduction](https://developer.cisco.com/docs/nexus-dashboard/4-1/)
