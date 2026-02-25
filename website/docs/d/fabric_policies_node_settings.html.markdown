---
layout: "mso"
page_title: "MSO: mso_fabric_policies_node_settings"
sidebar_current: "docs-mso-data-source-fabric_policies_node_settings"
description: |-
  Data source for Node Settings on Cisco Nexus Dashboard Orchestrator (NDO)
---



# mso_fabric_policies_node_settings #

Data source for Node Settings on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> Node Settings

## Example Usage ##

```hcl
data "mso_fabric_policies_node_settings" "node_settings" {
  template_id = mso_template.fabric_policy_template.id
  name        = "node_settings"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the Node Settings.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Node Settings.
* `id` - (Read-Only) The unique Terraform identifier of the Node Settings.
* `description` - (Read-Only) The description of the Node Settings.
* `synce` - (Read-Only) The Synchronous Ethernet (SyncE) configuration map of the Node Settings.
  * `admin_state` - (Read-Only) The SyncE administrative state of the Node Settings.
  * `quality_level` - (Read-Only) The SyncE quality level of the Node Settings.
* `ptp` - (Read-Only) The Precision Time Protocol (PTP) configuration map of the Node Settings.
  * `node_domain` - (Read-Only) The PTP domain of the Node Settings.
  * `priority_2` - (Read-Only) The PTP priority 2 of the Node Settings.
