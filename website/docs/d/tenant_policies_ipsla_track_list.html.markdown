---
layout: "mso"
page_title: "MSO: mso_tenant_policies_ipsla_track_list"
sidebar_current: "docs-mso-data-tenant_policies_ipsla_track_list"
description: |-
  Data source for IPSLA Track List.
---


# mso_tenant_policies_ipsla_track_list #

Data source for IPSLA Track List. This data source is supported in NDO v4.4(1) or higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> IPSLA Track List

## Example Usage ##

```hcl
data "mso_tenant_policies_ipsla_track_list" "ipsla_track_list" {
  template_id    = mso_template.tenant_policy_template.id
  name           = "ipsla_track_list"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Tenant Policy template.
* `name` - (Required) The name of the IPSLA Track List.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the IPSLA Track List.
* `id` - (Read-Only) The unique terraform identifier of the IPSLA Track List.
* `description` - (Read-Only) The description of the IPSLA Track List.
* `type` - (Read-Only) The threshold type of the IPSLA Track List.
* `threshold_up` - (Read-Only) The IPSLA Track List percentage or weight up threshold.
* `threshold_down` - (Read-Only) The IPSLA Track List percentage or weight down threshold.
* `members` - (Read-Only) The list of IPSLA Track List members.
  * `destination_ip` - (Read-Only) The destination IP of the member.
  * `ipsla_monitoring_policy_uuid` - (Read-Only) The UUID of the IPSLA Monitoring Policy for the member.
  * `scope_type` - (Read-Only) The scope type of the member.
  * `scope_uuid` - (Read-Only) The UUID of the BD or L3Out used as the scope for the member.
  * `weight` - (Read-Only) The weight of the member.
