---
layout: "mso"
page_title: "MSO: mso_fabric_policies_macsec_policy"
sidebar_current: "docs-mso-resource-fabric_policies_macsec_policy"
description: |-
  Manages MACsec Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_policies_macsec_policy #

Manages MACsec Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> MACsec

## Example Usage ##

```hcl
resource "mso_fabric_policies_macsec_policy" "macsec_policy" {
  template_id            = mso_template.fabric_policy_template.id
  name                   = "macsec_policy"
  description            = "Example description"
  admin_state            = "enabled"
  interface_type         = "access"
  cipher_suite           = "256GcmAes"
  window_size            = 128
  security_policy        = "shouldSecure"
  sak_expire_time        = 60
  confidentiality_offset = "offset30"
  key_server_priority    = 8
  macsec_keys {
    key_name             = "abc123"
    psk                  = "AA111111111111111111111111111111111111111111111111111111111111aa"
    start_time           = "now"
    end_time             = "2027-09-23 00:00:00"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the MACsec Policy.
* `description` - (Optional) The description of the MACsec Policy.
* `admin_state` - (Optional) The administrative state of the MACsec Policy. Allowed values are `enabled` or `disabled`. Defaults to `enabled` when unset during creation.
* `interface_type` - (Optional) The type of the interfaces the MACsec Policy will be applied to. Allowed values are `fabric` or `access`.
* `cipher_suite` - (Optional) The cipher suite of the MACsec Policy to be used for encryption. Allowed values are `128GcmAes`, `128GcmAesXpn`, `256GcmAes` or `256GcmAesXpn`. Defaults to `256GcmAesXpn` when unset during creation.
* `window_size` - (Optional) The window size of the MACsec Policy. It defines the maximum number of frames that can be received out of order before a replay attack is detected. Valid range: 0-4294967295. Defaults to 0 for `fabric` type or to 64 for `access` type when unset during creation.
* `security_policy` - (Optional) The security policy to allow traffic on the link for the MACsec Policy. Allowed values are `shouldSecure` or `mustSecure`. Defaults to `shouldSecure` when unset during creation.
* `sak_expire_time` - (Optional) The expiry time for the Security Association Key (SAK) for the MACsec Policy. Allowed value is 0 or valid range: 60-2592000. Defaults to 0 when unset during creation.
* `confidentiality_offset` - (Optional) The confidentiality offset for the MACsec Policy. This parameter is only configurable for `access` type. Allowed values are `offset0`, `offset30` or `offset50`. Defaults to `offset0` when unset during creation.
* `key_server_priority` - (Optional) The key server priority for the MACsec Policy. This parameter is only configurable for `access` type. Valid range: 0-255. Defaults to 16 when unset during creation.
* `macsec_keys` - (Optional) The list of MACsec Keys.
  * `macsec_keys.key_name` - (Required) The name of the MACsec Key. Key Name should contain hexadecimal characters [0-9a-fA-F].
  * `macsec_keys.psk` - (Required) The Pre-Shared Key (PSK) for the MACsec Key. PSK should contain hexadecimal characters [0-9a-fA-F]. PSK should be 64 characters long if cipher suite is `256GcmAes` or `256GcmAesXpn`. PSK should be 32 characters long if cipher suite is `128GcmAes` or `128GcmAesXpn`.
  * `macsec_keys.start_time` - (Optional) The start time for the MACsec Key. Allowed values are of the following format `YYYY-MM-DD HH:MM:SS` or `now`. The start time for each Key should be unique.
  * `macsec_keys.end_time` - (Optional) The end time for the MACsec Key. Allowed values are of the following format `YYYY-MM-DD HH:MM:SS` or `infinite`.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the MACsec Policy.
* `id` - (Read-Only) The unique Terraform identifier of the MACsec Policy.

## Importing ##

An existing MSO MACsec Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_macsec_policy.macsec_policy templateId/{template_id}/macsecPolicy/{name}
```
