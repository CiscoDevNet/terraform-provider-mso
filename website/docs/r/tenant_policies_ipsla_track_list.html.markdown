---
layout: "mso"
page_title: "MSO: mso_tenant_policies_ipsla_track_list"
sidebar_current: "docs-mso-resource-tenant_policies_ipsla_track_list"
description: |-
  Manages IPSLA Track List on Cisco Nexus Dashboard Orchestrator (NDO)
---


# mso_tenant_policies_ipsla_track_list #

Manages IPSLA Track List on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.4(1) or higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> IPSLA Track List

## Example Usage ##

```hcl
resource "mso_tenant_policies_ipsla_track_list" "ipsla_track_list" {
  template_id    = mso_template.tenant_policy_template.id
  name           = "ipsla_track_list"
  description    = "Terraform test IPSLA Track List"
  threshold_down = 11
  threshold_up   = 12
  type           = "weight"
  members {
    destination_ip               = "1.1.1.1"
    ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla_monitoring_policy.uuid
    scope_type                   = "bd"
    scope_uuid                   = mso_schema_template_bd.example_bd.uuid
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Tenant Policy template.
* `name` - (Required) The name of the IPSLA Track List.
* `description` - (Optional) The description of the IPSLA Track List.
* `type` - (Optional) The threshold type of the IPSLA Track List. Default: `percentage`. Valid values: `percentage`, `weight`.
* `threshold_up` - (Optional) The IPSLA Track List percentage or weight up threshold. The value must be in the range 0 - 100 when `type=percentage`. The value must be in the range 0 - 255 when `type=weight`. The value must be greater than or equal to `threshold_down`. Default: `0`.
* `threshold_down` - (Optional) The IPSLA Track List percentage or weight down threshold. The value must be in the range 0 - 100 when `type=percentage`. The value must be in the range 0 - 255 when `type=weight`. The value must be less than or equal to `threshold_up`. Default: `0`.
* `members` - (Optional) The list of IPSLA Track List members.
  * `destination_ip` - (Required) The destination IP of the member. Must be a valid IPv4 or IPv6 address.
  * `ipsla_monitoring_policy_uuid` - (Required) The UUID of the IPSLA Monitoring Policy for the member.
  * `scope_type` - (Required) The scope type of the member. Valid values: `bd`, `l3out`.
  * `scope_uuid` - (Required) The UUID of the BD or L3Out used as the scope for the member.
  * `weight` - (Read-Only) The weight of the member.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the IPSLA Track List.
* `id` - (Read-Only) The unique terraform identifier of the IPSLA Track List.

## Importing ##

An existing MSO IPSLA Track List can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_tenant_policies_ipsla_track_list.ipsla_track_list templateId/{template_id}/ipslaTrackLists/{name}
```