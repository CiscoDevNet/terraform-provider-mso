---
layout: "mso"
page_title: "MSO: mso_dhcp_option_policy_option"
sidebar_current: "docs-mso-resource-dhcp_option_policy_option"
description: |-
  Data source for MSO DHCP Option Policy Option
---

# mso_dhcp_option_policy_option

Data source for MSO DHCP Option Policy

## Example Usage

```hcl
data "mso_dhcp_option_policy_option" "example" {
  option_name = "dhcpOptionPol"
  option_policy_name = mso_dhcp_option_policy.example.name
}
```

## Argument Reference

- `option_name` - (Required) The name of the DHCP Option Policy Option.
- `option_policy_name` - (Required) The name of the DHCP Option Policy

## Attribute Reference

- `option_id` - The ID of the DHCP option (integer).
- `option_data` - The DHCP Option Data.
