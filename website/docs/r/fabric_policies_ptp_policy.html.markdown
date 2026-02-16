---
layout: "mso"
page_title: "MSO: mso_fabric_policies_ptp_policy"
sidebar_current: "docs-mso-resource-fabric_policies_ptp_policy"
description: |-
  Manages PTP Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_ptp_policy #

Manages PTP Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> PTP Policy

## Example Usage ##

```hcl
resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
  template_id            = mso_template.fabric_policy_template.id
  name                   = "ptp_policy"
  description            = "Example description"
  admin_state            = "enabled"
  TODO
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the PTP Policy.
* `description` - (Optional) The description of the PTP Policy.
* `admin_state` - (Required) The administrative state of the PTP Policy. Allowed values are `enabled` or `disabled`.
* `global_priority1` - (Required) The global priority1 of the PTP Policy. Valid range: 1-255.
* `global_priority2` - (Required) The global priority2 of the PTP Policy. Valid range: 1-255.
* `global_domain` - (Required) The global domain of the PTP Policy. Valid range: 0-128.
* `fabric_profile_template` - (Required) The fabric profile template of the PTP Policy. Allowed values are `default`, `aes67`, or `smpte`.
* `fabric_announce_interval` - (Required) The fabric announce interval in log base 2 seconds of the PTP Policy. Valid range: -3 to 4.
* `fabric_sync_interval` - (Required) The fabric sync interval in log base 2 seconds of the PTP Policy. Valid range: -4 to 1.
* `fabric_delay_interval` - (Required) The fabric delay interval in log base 2 seconds of the PTP Policy. Valid range: -4 to 5.
* `fabric_announce_timeout` - (Required) The fabric announce interval timeout count of the PTP Policy. Valid range: 2 to 10.
* `ptp_profiles` - (Optional) The list of PTP Profiles.
  * `name` - (Required) The name of the PTP Profile.
  * `description` - (Optional) The description of the PTP Profile.
  * `delay_interval` - (Required) The minimum delay request interval in log base 2 seconds of the PTP Profile. Valid range: -4 to 5.
  * `profile_template` - (Required) The profile template of the PTP Profile. Allowed values are `default`, `aes67`, `smpte` or `telecom`.
  * `sync_interval` - (Required) The sync interval in log base 2 seconds of the PTP Profile.Valid range: -4 to 1.
  * `override_node_profile` - (Optional) The node profile override of the PTP Profile. Allowed Values: `true` or `false`.
  * `announce_timeout` - (Required) The announce interval timeout count of the PTP Profile. Valid range: 2 to 10.
  * `announce_interval` - (Required) The announce interval in log base 2 seconds of the PTP Profile. Valid range: -3 to 4.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the PTP Policy.
* `id` - (Read-Only) The unique Terraform identifier of the PTP Policy.

## Importing ##

An existing MSO PTP Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_ptp_policy.ptp_policy templateId/{template_id}/ptpPolicy/{name}
```
