---
layout: "mso"
page_title: "MSO: mso_dhcp_option_policy"
sidebar_current: "docs-mso-resource-dhcp_option_policy"
description: |-
  Manages MSO DHCP Option Policy
---

# mso_dhcp_option_policy

Manages MSO DHCP Option Policy

## Example Usage

```hcl
resource "mso_dhcp_option_policy" "example" {
  tenant_id = mso_tenant.example.id
  name = "dhcpOptionPol"
  description = "from Terraform"
  option {
    name = "op1"
    data = "d1"
    id = "1"
  }
  option {
    name = "op2"
    data = "d2"
    id = "2"
  }
}

```

## Argument Reference

- `name` - (Required) The name of the DHCP option policy.
- `tenant_id` - (Required) ID of parent `mso_tenant` resource.
- `description` - (Optional) The description for this DHCP Option Policy.
- `option` - (Optional) DHCP Option configuration to be associated to the Policy.
  - `name` - (Required) The name of the DHCP option.
  - `id` - (Required) The ID of the DHCP option (integer).
  - `data` - (Optional) The DHCP Option Data.

## Attribute Reference

The only Attribute exposed for this resource is `id`. Which is set to the id of tenant created.

## Importing

An existing MSO Tenant can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_dhcp_option_policy.example {dhcp_option_policy_id}
```
