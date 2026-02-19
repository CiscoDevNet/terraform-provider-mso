---
layout: "mso"
page_title: "MSO: mso_fabric_policies_ptp_policy_profile"
sidebar_current: "docs-mso-data-source-fabric_policies_ptp_policy"
description: |-
  Data source for PTP Profiles on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_ptp_policy_profile #

Data source for (Precision Time Protocol) PTP Policy Profiles on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> PTP Policy -> Profiles

## Example Usage ##

```hcl
data "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
  template_id = mso_template.fabric_policy_template.id
  name        = "ptp_profile"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the PTP Profile.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the PTP Profile.
* `ptp_policy_uuid` - (Read-Only) The NDO UUID of the PTP Policy.
* `id` - (Read-Only) The unique Terraform identifier of the PTP Profile.
* `description` - (Read-Only) The description of the PTP Profile.
* `delay_interval` - (Read-Only) The delay interval of the PTP Profile.
* `profile_template` - (Read-Only) The profile template of the PTP Profile.
* `sync_interval` - (Read-Only) The sync interval of the PTP Profile.
* `override_node_profile` - (Read-Only) The node profile override of the PTP Profile.
* `announce_timeout` - (Read-Only) The announce timeout of the PTP Profile.
* `announce_interval` - (Read-Only) The announce interval of the PTP Profile.
* `local_priority` - (Read-Only) The local priority of the PTP Profile.
* `destination_mac_type` - (Read-Only) The destination MAC for PTP messages of the PTP Profile.
* `mismatched_mac_handling` - (Read-Only) The mismatched destination MAC handling of the PTP Profile.
