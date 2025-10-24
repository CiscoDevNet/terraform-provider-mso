---
layout: "mso"
page_title: "MSO: mso_fabric_policies_synce_interface_policy"
sidebar_current: "docs-mso-resource-fabric_policies_synce_interface_policy"
description: |-
  Manages SyncE Interface Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_policies_synce_interface_policy #

Manages SyncE Interface Policys on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> SyncE Interface Policy

## Example Usage ##

```hcl
resource "mso_fabric_policies_synce_interface_policy" "synce_interface_policy" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "synce_interface_policy"
  description     = "Example description"
	admin_state     = "enabled"
  sync_state_msg  = "enabled"
  selection_input = "enabled"
  src_priority    = 120
  wait_to_restore = 6
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the SyncE Interface Policy.
* `description` - (Optional) The description of the SyncE Interface Policy.
* `admin_state` - (Optional) The administrative state of the SyncE Interface Policy. Allowed values are `enabled` and `disabled`. Default to `disabled`.
* `sync_state_msg` - (Optional) The sync state message of the SyncE Interface Policy. Allowed values are `enabled` and `disabled`. Default to `disabled`.
* `selection_input` - (Optional) The selection input of the SyncE Interface Policy. Allowed values are `enabled` and `disabled`. Default to `disabled`.
* `src_priority` - (Optional) The source priority of the SyncE Interface Policy. Valid range: 1-254. Default to 100
* `wait_to_restore` - (Optional) The delay before attempting to restore synchronization on a SyncE Interface after a disruption. Valid range: 0-12. Default to 5

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the SyncE Interface Policy.
* `id` - (Read-Only) The unique Terraform identifier of the SyncE Interface Policy.

## Importing ##

An existing MSO SyncE Interface Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_synce_interface_policy.synce_interface_policy templateId/{template_id}/SyncEInterfacePolicy/{name}
```
