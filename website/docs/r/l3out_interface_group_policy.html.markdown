---
layout: "mso"
page_title: "MSO: mso_l3out_interface_group_policy"
subcategory: "L3Out"
sidebar_current: "docs-mso-resource-l3out_interface_group_policy"
description: |-
  Manages one interface group policy on an L3Out in Nexus Dashboard Orchestration.
---

# mso_l3out_interface_group_policy (Resource)

Manages one named interface group policy on an existing L3Out. The L3Out template, L3Out, and any referenced tenant routing, custom QoS, or NetFlow policies must already exist.

When an optional setting is omitted, the provider reads and retains its value from Orchestration. Set a clearable string to `""` to clear it. For BFD and OSPF objects, `{}` and `enabled = false` both disable the object. The NetFlow reference map has its own ownership rule below.

-> **Authentication key limitation:** Orchestration accepts the BFD and OSPF keys but returns only an opaque key reference. The provider retains a previously configured key in Terraform state so that refresh and unrelated updates do not discard it. It cannot verify the stored key against Orchestration or detect a key changed outside Terraform. An imported interface group has no key value to recover. Treat Terraform state as sensitive because it contains configured keys.

## API Information

- **APIs**: Nexus Dashboard Orchestration API (template endpoints).
- **Create**: `GET /mso/api/v1/templates/{template_id}`, then `PATCH /mso/api/v1/templates/{template_id}` to add an `interfaceGroups` entry.
- **Read**: `GET /mso/api/v1/templates/{template_id}` and select the L3Out by UUID and group by name.
- **Update**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` with changes to the selected group.
- **Delete**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` to remove the selected group.
- **Target versions**: Nexus Dashboard 4.1+ (Orchestration 5.1+).
- **Terraform ID**: `{template_id}/{l3out_uuid}/{name}`.
- **Managed resource identity**: `template_id`, `l3out_uuid`, and `name`.

The provider sends template PATCH requests with `validate=false` and refreshes the selected object from the template GET response.

## GUI Information

In Nexus Dashboard, choose **Manage > Orchestration > Tenant Templates > L3Out** and open the template. In the L3Out view, click **+Create Node/Interface Group Policy**, then choose **Interface**. Cisco describes the [interface group workflow](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html).

## Example Usage

```terraform
resource "mso_l3out_interface_group_policy" "edge" {
  template_id                   = mso_template.l3out.id
  l3out_uuid                    = mso_l3out.example.uuid
  name                          = "edge"
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
```

The [complete standalone interface group example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/master/examples/l3out_interface_group_policy) creates the L3Out, the referenced tenant policies, a NetFlow exporter and record, all four monitors, and BFD, multi-hop BFD, and OSPF authentication.

## Schema

### Required

- `template_id` (String) The UUID of the existing L3Out template.
    - **API**: template `templateId`, request path `{template_id}`.
    - **Force New**: `true`.
- `l3out_uuid` (String) The UUID of the parent L3Out.
    - **API**: `l3outTemplate.l3outs[].uuid`.
    - **Force New**: `true`.
- `name` (String) The name of the interface group policy within its parent L3Out.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].name`.
    - **Force New**: `true`.

### Optional

- `description` (String) The description of the interface group policy.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].description`.
- `interface_routing_policy_uuid` (String) The UUID of the referenced `mso_tenant_policies_l3out_interface_routing_policy`.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].interfaceRoutingPolicyRef`.
- `custom_qos_policy_uuid` (String) The UUID of the referenced `mso_tenant_policies_custom_qos_policy`.
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
        - **Validation**: True requires `key_id` and `key`. Authentication settings require `enabled = true`.
    - `key_id` (Number) The ID of the single-hop BFD authentication key.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.keyID`.
        - **Validation**: Allowed only when authentication is enabled.
    - `key` (String, Sensitive) The single-hop BFD authentication key. Orchestration does not return the secret, so the provider preserves a previously configured key in Terraform state.
        - **API write**: `l3outTemplate.l3outs[].interfaceGroups[].bfd.key.value`.
        - **Validation**: Allowed only when authentication is enabled.
- `bfd_multi_hop` (Object) The multi-hop BFD settings of the interface group.
    - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop`.
    - `enabled` (Boolean) Whether multi-hop BFD is enabled.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.enabled` or absence of `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop`.
        - **Default**: `false`.
    - `authentication_enabled` (Boolean) Whether multi-hop BFD authentication is enabled.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.authEnabled`.
        - **Validation**: True requires `key_id` and `key`. Authentication settings require `enabled = true`.
    - `key_id` (Number) The ID of the multi-hop BFD authentication key.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].bfdMultiHop.keyID`.
        - **Validation**: Allowed only when authentication is enabled.
    - `key` (String, Sensitive) The multi-hop BFD authentication key. Orchestration does not return the secret, so the provider preserves a previously configured key in Terraform state.
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
        - **Validation**: `simple` and `md5` require `key_id` and `key`. Other settings require `enabled = true`.
    - `key_id` (Number) The ID of the OSPF authentication key.
        - **API**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.keyID`.
        - **Validation**: Allowed only with `simple` or `md5`.
    - `key` (String, Sensitive) The OSPF authentication key. Orchestration does not return the secret, so the provider preserves a previously configured key in Terraform state.
        - **API write**: `l3outTemplate.l3outs[].interfaceGroups[].ospf.key.value`.
        - **Validation**: Allowed only with `simple` or `md5`.

### Read-Only

- `id` (String, provider generated) The stable Terraform ID composed of the template UUID, L3Out UUID, and group name.

## Importing

```shell
terraform import mso_l3out_interface_group_policy.edge '<template_id>/<l3out_uuid>/<name>'
```

Terraform 1.12 and later can also import with the managed resource identity:

```terraform
import {
  to = mso_l3out_interface_group_policy.edge
  identity = {
    template_id = "<template_id>"
    l3out_uuid  = "<l3out_uuid>"
    name        = "edge"
  }
}
```

## References

### Related Terraform Objects

- [mso_l3out resource](/docs/providers/mso/r/l3out.html)
- [mso_l3out_interface_group_policy data source](/docs/providers/mso/d/l3out_interface_group_policy.html)
- [mso_tenant_policies_custom_qos_policy resource](/docs/providers/mso/r/tenant_policies_custom_qos_policy.html)
- [mso_tenant_policies_l3out_interface_routing_policy resource](/docs/providers/mso/r/tenant_policies_l3out_interface_routing_policy.html)
- [mso_tenant_policies_netflow_monitor resource](/docs/providers/mso/r/tenant_policies_netflow_monitor.html)

### External Documentation

- [Cisco Nexus Dashboard 4.1.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html)
- [Cisco Nexus Dashboard Orchestration Get Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/get-template/)
- [Cisco Nexus Dashboard Orchestration Patch Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/patch-template/)
- [Cisco Nexus Dashboard 4.1.1 API introduction](https://developer.cisco.com/docs/nexus-dashboard/4-1/)
