---
layout: "mso"
page_title: "MSO: mso_fabric_policies_ptp_policy_profile"
sidebar_current: "docs-mso-resource-fabric_policies_ptp_policy_profile"
description: |-
  Manages PTP Policy Profiles on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_ptp_policy_profile #

Manages (Precision Time Protocol) PTP Policy Profiles on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> PTP Policy -> Profiles

## Example Usage ##

```hcl
resource "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
  template_id         = mso_template.fabric_policy_template.id
  ptp_policy_uuid       = mso_fabric_policies_ptp_policy.ptp_policy.uuid
  name                  = "ptp_policy_profile"
  description           = "Example description"
  delay_interval        = -2
  sync_interval         = -3
  announce_timeout      = 3
  announce_interval     = 1
  profile_template      = "aes67"
  override_node_profile = false
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `ptp_policy_uuid` - (Required) The NDO UUID of the PTP Policy.
* `name` - (Required) The name of the PTP Profile.
* `description` - (Optional) The description of the PTP Profile.
* `profile_template` - (Required) The profile template of the PTP Profile. Allowed values are `default`, `aes67`, `smpte` or `telecom`.
* `delay_interval` - (Required) The minimum delay request interval in log base 2 seconds of the PTP Profile. Valid range: -4 to 5.
* `announce_timeout` - (Required) The announce interval timeout count of the PTP Profile. Valid range: 2 to 10.
* `announce_interval` - (Required) The announce interval in log base 2 seconds of the PTP Profile. Valid range: -3 to 4.
* `sync_interval` - (Required) The sync interval in log base 2 seconds of the PTP Profile.Valid range: -4 to 1.
* `override_node_profile` - (Optional) The node profile override of the PTP Profile. This parameter is not applicable when `profile_template` is `telecom`. Allowed Values: `true` or `false`.
* `local_priority` - (Optional) The local priority of the PTP Profile. This parameter is only applicable when `profile_template` is `telecom`. Valid range: 1 to 128.
* `destination_mac_type` - (Optional) The destination MAC for PTP messages of the PTP Profile. This parameter is only applicable when `profile_template` is `telecom`. Allowed values are `forwardable` or `non_forwardable`.
* `mismatched_mac_handling` - (Optional) The mismatched destination MAC handling of the PTP Profile. This parameter is only applicable when `profile_template` is `telecom`. Allowed values are `drop`, `reply_with_config_mac` or `reply_with_received_mac`.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the PTP Profile.
* `id` - (Read-Only) The unique Terraform identifier of the PTP Profile.

## Importing ##

An existing MSO PTP Profile can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_ptp_policy_profile.ptp_policy_profile templateId/{template_id}/ptpPolicyProfile/{name}
```
