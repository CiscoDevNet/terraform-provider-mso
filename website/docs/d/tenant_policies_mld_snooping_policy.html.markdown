---
layout: "mso"
page_title: "MSO: mso_tenant_policies_mld_snooping_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_mld_snooping_policy"
description: |-
  Data source for MLD Snooping Policy.
---

# mso_tenant_policies_mld_snooping_policy #

Data source for Multicast Listener Discovery (MLD) Snooping Policy. This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> MLD Snooping Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_mld_snooping_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the MLD Snooping Policy to retrieve.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the MLD Snooping Policy.
* `id` - (Read-Only) The unique terraform identifier of the MLD Snooping Policy in the template.
* `description` - (Read-Only) The description of the MLD Snooping Policy.
* `admin_state` - (Read-Only) The administrative state of the MLD Snooping Policy.
* `fast_leave_control` - (Read-Only) Enable or disable fast leave processing. When enabled, the switch immediately removes a multicast group when it receives an MLD Done message.
* `querier_control` - (Read-Only) Enable or disable querier control. When enabled, the switch can act as an MLD querier.
* `querier_version` - (Read-Only) The MLD querier version.
* `query_interval` - (Read-Only) The interval in seconds between MLD general query messages.
* `query_response_interval` - (Read-Only) The maximum time in seconds that hosts can wait before responding to an MLD query.
* `last_member_query_interval` - (Read-Only) The interval in seconds between MLD group-specific queries sent in response to an MLD Done message.
* `start_query_interval` - (Read-Only) The interval in seconds between MLD queries sent at startup.
* `start_query_count` - (Read-Only) The number of MLD queries sent at startup.
