---
layout: "mso"
page_title: "MSO: mso_dhcp_relay_policy_provider"
sidebar_current: "docs-mso-resource-dhcp_relay_policy_provider"
description: |-
  Data source for MSO DHCP Relay Policy Provider
---

# mso_dhcp_relay_policy_provider

Data source for MSO DHCP Relay Policy Provider

## Example Usage

```hcl
data "mso_dhcp_relay_policy_provider" "example" {
    dhcp_relay_policy_name = mso_dhcp_relay_policy.example.name
    dhcp_server_address = mso_dhcp_relay_policy_provider.example.dhcp_server_address
    external_epg_ref = mso_dhcp_relay_policy_provider.example.dhcp_server_address
}
```

## Argument Reference
- `dhcp_relay_policy_name` - (Required) The DHCP Relay Policy Name.
- `dhcp_server_address` - (Required) The DHCP Server Address.
- `epg_ref` - (Optional) The reference of the EPG.
- `external_epg_ref` - (Optional) The reference of the External EPG.

## Note

external_epg_ref and epg_ref both are not expected simultaneously, only one of them is required at a time.
