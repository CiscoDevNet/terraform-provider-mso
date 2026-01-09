---
layout: "mso"
page_title: "MSO: mso_tenant_policies_dhcp_option_policy"
sidebar_current: "docs-mso-data-source-dhcp_option_policy"
description: |-
  Data source for DHCP Option Policy.
---

# mso_tenant_policies_dhcp_option_policy #

Data source for Dynamic Host Configuration Protocol (DHCP) Option Policies on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> DHCP Option Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_dhcp_option_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the DHCP Option Policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the DHCP Option Policy.
* `id` - (Read-Only) The unique Terraform identifier of the DHCP Option Policy in the template.
* `description` - (Read-Only) The description of the DHCP Option Policy. When unset during creation, no description is applied.
* `options` - (Read-Only) A set of DHCP options. At least one option is required. The entire set of options is replaced during updates.
  * `name` - (Read-Only) The name of the DHCP option. This is a descriptive label for the option.
  * `id` - (Read-Only) The DHCP option ID.
  * `data` - (Read-Only) The value/data for the DHCP option. The format depends on the option type.
