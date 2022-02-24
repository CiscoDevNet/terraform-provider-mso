---
layout: "mso"
page_title: "MSO: mso_dhcp_option_policy_option"
sidebar_current: "docs-mso-resource-dhcp_option_policy_option"
description: |-
  Manages MSO DHCP Option Policy Option
---

# mso_dhcp_option_policy_option

Manages MSO DHCP Option Policy Option

## Example Usage

```hcl
resource "mso_dhcp_option_policy" "example" {
  tenant_id = mso_tenant.example.id
  name = "example"
}

resource "mso_dhcp_option_policy_option" "example"{
    dhcp_option_policy_name = mso_dhcp_option_policy.example.name
    dhcp_option_name = "example"
}

```

## Argument Reference

- `option_name` - (Required) The name of the DHCP option policy option.
- `option_policy_name` - (Required) Policy Name of parent `mso_dhcp_option_policy` resource.
- `option_id` - (Optional) DHCP Option Policy Option Id (Integer).
- `option_data` - (Optional) The DHCP Option Data.

## Attribute Reference

The only Attribute exposed for this resource is `id`. Which is set to the id of dhcp option policy created.

## Importing

An existing MSO Tenant can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_dhcp_option_policy_option.example {dhcp_option_policy_option_id}
```
