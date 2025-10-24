---
layout: "mso"
page_title: "MSO: mso_fabric_policies_synce_interface_policy"
sidebar_current: "docs-mso-data-source-fabric_policies_synce_interface_policy"
description: |-
  Data source for SyncE Interface Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_policies_synce_interface_policy #

Data source for SyncE Interface Policies on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> SyncE Interface Policy

## Example Usage ##

```hcl
data "mso_fabric_policies_synce_interface_policy" "synce_interface_policy" {
  template_id = mso_template.fabric_policy_template.id
  name        = "synce_interface_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the SyncE Interface Policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the SyncE Interface Policy.
* `id` - (Read-Only) The unique Terraform identifier of the SyncE Interface Policy.
* `description` - (Read-Only) The description of the SyncE Interface Policy.
* `admin_state` - (Read-Only) The administrative state of the SyncE Interface Policy.
* `sync_state_msg` - (Read-Only) The sync state message of the SyncE Interface Policy.
* `selection_input` - (Read-Only) The selection input of the SyncE Interface Policy.
* `src_priority` - (Read-Only) The source priority of the SyncE Interface Policy.
* `wait_to_restore` - (Read-Only) The delay before attempting to restore synchronization on a SyncE Interface after a disruption.
