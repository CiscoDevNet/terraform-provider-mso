---
layout: "mso"
page_title: "MSO: mso_tenant_policies_mld_snooping_policy"
sidebar_current: "docs-mso-resource-tenant_policies_mld_snooping_policy"
description: |-
  Manages MLD Snooping Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_mld_snooping_policy #

Manages Multicast Listener Discovery (MLD) Snooping Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> MLD Snooping Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
  template_id                = mso_template.template_tenant.id
  name                       = "test_mld_snooping_policy"
  description                = "Test MLD Snooping Policy"
  admin_state                = "enabled"
  fast_leave_control         = true
  querier_control            = true
  querier_version            = "v2"
  query_interval             = 125
  query_response_interval    = 10
  last_member_query_interval = 1
  start_query_interval       = 31
  start_query_count          = 2
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the MLD Snooping Policy.
* `description` - (Optional) The description of the MLD Snooping Policy.
* `admin_state` - (Optional) The administrative state of the MLD Snooping Policy. Default: `disabled` when unset during creation. Valid values: `enabled`, `disabled`.
* `fast_leave_control` - (Optional) Enable or disable fast leave processing. When enabled, the switch immediately removes a multicast group when it receives an MLD Done message. Default: `false` when unset during creation.
* `querier_control` - (Optional) Enable or disable querier control. When enabled, the switch can act as an MLD querier. Default: `false` when unset during creation.
* `querier_version` - (Optional) The MLD querier version. Default: `v2` when unset during creation. Valid values: `v1`, `v2`.
* `query_interval` - (Optional) The interval in seconds between MLD general query messages. Default: 125 when unset during creation. Valid range: 1-18000 seconds.
* `query_response_interval` - (Optional) The maximum time in seconds that hosts can wait before responding to an MLD query. Default: 10 when unset during creation. Valid range: 1-25 seconds.
* `last_member_query_interval` - (Optional) The interval in seconds between MLD group-specific queries sent in response to an MLD Done message. Default: 1 when unset during creation. Valid range: 1-25 seconds.
* `start_query_interval` - (Optional) The interval in seconds between MLD queries sent at startup. Default: 31 when unset during creation. Valid range: 1-18000 seconds.
* `start_query_count` - (Optional) The number of MLD queries sent at startup. Default: 2 when unset during creation. Valid range: 1-10.

## Attribute Reference ##

* `uuid` - The NDO UUID of the MLD Snooping Policy.
* `id` - The unique terraform identifier of the MLD Snooping Policy in the template.

## Importing ##

An existing MSO MLD Snooping Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_mld_snooping_policy.mld_policy templateId/{template_id}/MLDSnoopingPolicy/{name}
```
