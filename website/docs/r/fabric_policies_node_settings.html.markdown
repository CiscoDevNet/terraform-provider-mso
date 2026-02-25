---
layout: "mso"
page_title: "MSO: mso_fabric_policies_node_settings"
sidebar_current: "docs-mso-resource-fabric_policies_node_settings"
description: |-
  Manages Node Settings on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_node_settings #

Manages Node Settings on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3(1) or higher.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> Node Settings

## Example Usage ##

```hcl
resource "mso_fabric_policies_node_settings" "node_settings" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "tf_node_settings"
  description     = "Terraform Node Settings Policy"
  synce = {
    admin_state   = "enabled"
    quality_level = "option_2_generation_1"
  }
  ptp = {
    node_domain   = 25
    priority_2    = 99
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the Fabric Policy template.
* `name` - (Required) The name of the Node Settings.
* `description` - (Optional) The description of the Node Settings.
* `synce` - (Optional) The Synchronous Ethernet (SyncE) configuration map of the Node Settings.
  * `admin_state` - (Required) The SyncE administrative state of the Node Settings. Allowed values are `enabled` or `disabled`.
  * `quality_level` - (Required) The SyncE quality level of the Node Settings. Allowed values are `option_1`, `option_2_generation_1` or `option_2_generation_2`.
* `ptp` - (Optional) The Precision Time Protocol (PTP) configuration map of the Node Settings.
  * `node_domain` - (Required) The PTP domain of the Node Settings. Valid range: 24-43.
  * `priority_2` - (Required) The PTP priority 2 of the Node Settings. Valid range: 0-255.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Node Settings.
* `id` - (Read-Only) The unique Terraform identifier of the Node Settings.

## Importing ##

An existing MSO Node Settings can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_fabric_policies_node_settings.node_settings templateId/{template_id}/nodeSettings/{name}
```
