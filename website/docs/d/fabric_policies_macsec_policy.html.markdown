---
layout: "mso"
page_title: "MSO: mso_fabric_policies_macsec_policy"
sidebar_current: "docs-mso-data-source-fabric_policies_macsec_policy"
description: |-
  Data source for MACsec Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_policies_macsec_policy #

Data source for MACsec Policies on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> MACsec

## Example Usage ##

```hcl
data "mso_fabric_policies_macsec_policy" "macsec_policy" {
  template_id = mso_template.fabric_policy_template.id
  name        = "macsec_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the MACsec Policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the MACsec Policy.
* `id` - (Read-Only) The unique Terraform identifier of the MACsec Policy.
* `description` - (Read-Only) The description of the MACsec Policy.
* `admin_state` - (Read-Only) The administrative state of the MACsec Policy.
* `interface_type` - (Read-Only) The type of the interfaces the MACsec Policy will be applied to.
* `cipher_suite` - (Read-Only) The cipher suite of the MACsec Policy to be used for encryption.
* `window_size` - (Read-Only) The window size of the MACsec Policy. It defines the maximum number of frames that can be received out of order before a replay attack is detected.
* `security_policy` - (Read-Only) The security policy to allow traffic on the link for the MACsec Policy.
* `sak_expire_time` - (Read-Only) The expiry time for the Security Association Key (SAK) for the MACsec Policy.
* `confidentiality_offset` - (Read-Only) The confidentiality offset for the MACsec Policy.
* `key_server_priority` - (Read-Only) The key server priority for the MACsec Policy.
* `macsec_keys` - (Read-Only) The list of MACsec Keys.
  * `macsec_key.key_name` - (Read-Only) The name of the MACsec Key.
  * `macsec_key.psk` - (Read-Only) The Pre-Shared Key (PSK) for the MACsec Key.
  * `macsec_key.start_time` - (Read-Only) The start time for the MACsec Key.
  * `macsec_key.end_time` - (Read-Only) The end time for the MACsec Key.
