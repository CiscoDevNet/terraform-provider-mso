---
layout: "mso"
page_title: "MSO: mso_fabric_policies_ptp_policy"
sidebar_current: "docs-mso-data-source-fabric_policies_ptp_policy"
description: |-
  Data source for the PTP Policy on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_ptp_policy #

Data source for the (Precision Time Protocol) PTP Policy on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> PTP Policy

## Example Usage ##

```hcl
data "mso_fabric_policies_ptp_policy" "ptp_policy" {
  template_id = mso_template.fabric_policy_template.id
  name        = "ptp_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the PTP Policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the PTP Policy.
* `id` - (Read-Only) The unique Terraform identifier of the PTP Policy.
* `description` - (Read-Only) The description of the PTP Policy.
* `admin_state` - (Read-Only) The administrative state of the PTP Policy.
* `global_priority1` - (Read-Only) The global priority1 of the PTP Policy.
* `global_priority2` - (Read-Only) The global priority2 of the PTP Policy.
* `global_domain` - (Read-Only) The global domain of the PTP Policy.
* `fabric_profile_template` - (Read-Only) The fabric profile template of the PTP Policy.
* `fabric_announce_interval` - (Read-Only) The fabric announce interval of the PTP Policy.
* `fabric_sync_interval` - (Read-Only) The fabric sync interval of the PTP Policy.
* `fabric_delay_interval` - (Read-Only) The fabric delay interval of the PTP Policy.
* `fabric_announce_timeout` - (Read-Only) The fabric announce timeout of the PTP Policy.
