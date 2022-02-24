---
layout: "mso"
page_title: "MSO: mso_dhcp_option_policy"
sidebar_current: "docs-mso-resource-dhcp_option_policy"
description: |-
  Data source for MSO DHCP Option Policy
---

# mso_dhcp_option_policy

Data source for MSO DHCP Option Policy

## Example Usage

```hcl
data "mso_dhcp_option_policy" "example" {
  name = "dhcpOptionPol"
}
```

## Argument Reference

- `name` - (Required) The name of the DHCP option policy.

## Attribute Reference

- `tenant_id` - ID of parent `mso_tenant` resource.
- `description` - The description for this DHCP Option Policy.
- `option` - DHCP Option configuration to be associated to the Policy.
  - `name` - The name of the DHCP option.
  - `id` - The ID of the DHCP option (integer).
  - `data` - The DHCP Option Data.
