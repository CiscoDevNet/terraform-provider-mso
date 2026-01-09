---
layout: "mso"
page_title: "MSO: mso_tenant_policies_option_policy"
sidebar_current: "docs-mso-resource-dhcp_option_policy"
description: |-
  Manages DHCP Option Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_dhcp_option_policy #

Manages Dynamic Host Configuration Protocol (DHCP) Option Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> DHCP Option Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_dhcp_option_policy"
  description = "Test DHCP Option Policy"
  
  options {
    name = "example_server"
    id   = 1
    data = "8.8.8.8"
  }
  
  options {
    name = "domain_name"
    id   = 2
    data = "example.com"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the DHCP Option Policy.
* `description` - (Optional) The description of the DHCP Option Policy. When unset during creation, no description is applied.
* `options` - (Required) A set of DHCP options. At least one option is required. The entire set of options is replaced during updates.
  * `name` - (Required) The name of the DHCP option. This is a descriptive label for the option.
  * `id` - (Optional) The DHCP option ID. Valid range: 0-255.
  * `data` - (Optional) The value/data for the DHCP option. The format depends on the option type.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the DHCP Option Policy.
* `id` - (Read-Only) The unique Terraform identifier of the DHCP Option Policy in the template.

## Importing ##

An existing MSO DHCP Option Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_dhcp_option_policy.dhcp_policy templateId/{template_id}/DHCPOptionPolicy/{name}
```
