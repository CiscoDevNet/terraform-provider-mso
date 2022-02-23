---
layout: "mso"
page_title: "MSO: mso_dhcp_relay_policy"
sidebar_current: "docs-mso-resource-dhcp_relay_policy"
description: |-
  Data source for MSO DHCP Relay Policy
---

# mso_dhcp_relay_policy

Data source for MSO DHCP Relay Policy

## Example Usage

```hcl
data "mso_dhcp_relay_policy" "example" {
  name = "dhcpRelayPol"
}
```

## Argument Reference

- `name` - (Required) The name of the DHCP Relay policy.

## Attribute Reference

- `tenant_id` - ID of parent `mso_tenant` resource.
- `description` - The description for this DHCP Relay Policy.
- `dhcp_relay_policy_provider` - DHCP Provider configuration to be associated to the Policy.
  - `epg` - The reference of the EPG.
  - `external_epg` - The reference of the External EPG. external_epg and epg both should not be expected simultaneously.
  - `dhcp_server_address` - The DHCP Server Address.
