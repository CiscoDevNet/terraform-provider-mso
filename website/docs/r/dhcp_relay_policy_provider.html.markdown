---
layout: "mso"
page_title: "MSO: mso_dhcp_relay_policy_provider"
sidebar_current: "docs-mso-resource-dhcp_relay_policy_provider"
description: |-
  Manages MSO DHCP Relay Policy Provider
---

# mso_dhcp_relay_policy_provider

Manages MSO DHCP Relay Policy Provider

## Example Usage

```hcl
resource "mso_dhcp_relay_policy_provider" "example" {
    dhcp_relay_policy_name = mso_dhcp_relay_policy.example.name
    dhcp_server_address = "1.2.3.4"
    external_epg_ref = mso_schema_template_external_epg.test.id
}

```

## Argument Reference
- `dhcp_relay_policy_name` - (Required) The DHCP Relay Policy Name.
- `dhcp_server_address` - (Required) The DHCP Server Address.
- `epg_ref` - (Optional) The reference of the EPG.
- `external_epg_ref` - (Optional) The reference of the External EPG.


## Attribute Reference

The only Attribute exposed for this resource is `id`. Which is set to the id of DHCP Relay Policy Provider in {dhcp_relay_policy_name}/{epg_ref/external_epg_ref}/{dhcp_server_address} format.

## Note

external_epg_ref and epg_ref both are not expected simultaneously, only one of them is required at a time.

## Importing

An existing MSO Tenant can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_dhcp_relay_policy_provider.example {dhcp_relay_policy_name}/{epg_ref/external_epg_ref}/{dhcp_server_address}
```