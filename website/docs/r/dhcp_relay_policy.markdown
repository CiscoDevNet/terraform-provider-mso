---
layout: "mso"
page_title: "MSO: mso_dhcp_relay_policy"
sidebar_current: "docs-mso-resource-dhcp_relay_policy"
description: |-
  Manages MSO DHCP relay Policy
---

# mso_dhcp_relay_policy

Manages MSO DHCP relay Policy

## Example Usage

```hcl
resource "mso_dhcp_relay_policy" "example" {
  tenant_id = mso_tenant.example.id
  name = "dhcpRelayPol"
  description = "from Terraform"
  dhcp_relay_policy_provider {
    epg = mso_schema_template_anp_epg.example.id
    dhcp_server_address = "1.2.3.4"
  }
  dhcp_relay_policy_provider {
    external_epg = mso_schema_template_external_epg.example.id
    dhcp_server_address = "1.2.3.4"
  }
}

```

## Argument Reference

- `name` - (Required) The name of the DHCP relay policy.
- `tenant_id` - (Required) ID of parent `mso_tenant` resource.
- `description` - (Optional) The description for this DHCP Relay Policy.
- `dhcp_relay_policy_provider` - (Optional) DHCP Provider configuration to be associated to the Policy.
  - `epg` - (Required) The reference of the EPG.
  - `external_epg` - (Required) The reference of the External EPG. external_epg and epg both are not expected simultaneously, only one of them is required at a time.
  - `dhcp_server_address` - (Required) The DHCP Server Address.

## Attribute Reference

The only Attribute exposed for this resource is `id`. Which is set to the id of tenant created.

## Importing

An existing MSO Tenant can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_dhcp_relay_policy.example {dhcp_relay_policy_id}
```
