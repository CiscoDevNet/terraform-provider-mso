---
layout: "mso"
page_title: "MSO: mso_tenant_policies_custom_qos_policy"
sidebar_current: "docs-mso-resource-tenant_policies_custom_qos_policy"
description: |-
  Manages Custom QoS Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_custom_qos_policy #

Manages Custom Quality of Service (QoS) Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.


## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Custom QoS Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_custom_qos_policy"
  description = "Test Custom QoS Policy"
  
  dscp_mappings {
    dscp_from    = "af11"
    dscp_to      = "af12"
    dscp_target  = "af11"
    target_cos   = "background"
    qos_priority = "level1"
  }
  
  dscp_mappings {
    dscp_from    = "af21"
    dscp_to      = "af22"
    dscp_target  = "af21"
    target_cos   = "best_effort"
    qos_priority = "level2"
  }
  
  cos_mappings {
    dot1p_from   = "background"
    dot1p_to     = "best_effort"
    dscp_target  = "af11"
    target_cos   = "background"
    qos_priority = "level1"
  }
  
  cos_mappings {
    dot1p_from   = "excellent_effort"
    dot1p_to     = "critical_applications"
    dscp_target  = "af21"
    target_cos   = "video"
    qos_priority = "level2"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the Custom QoS Policy.
* `description` - (Optional) The description of the Custom QoS Policy. When unset during creation, no description is applied.
* `dscp_mappings` - (Optional) A set of DSCP (Differentiated Services Code Point) mappings. Each mapping defines how DSCP values are mapped to CoS values and priority levels. Multiple DSCP mappings can be configured.
  * `dscp_from` - (Optional) The starting DSCP value of the range. Default: `unspecified`. Valid values: `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0`, `cs1`, `cs2`, `cs3`, `cs4`, `cs5`, `cs6`, `cs7`, `expedited_forwarding`, `voice_admit`, `unspecified`.
  * `dscp_to` - (Optional) The ending DSCP value of the range. Default: `unspecified`. Valid values: `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0`, `cs1`, `cs2`, `cs3`, `cs4`, `cs5`, `cs6`, `cs7`, `expedited_forwarding`, `voice_admit`, `unspecified`.
  * `dscp_target` - (Optional) The target DSCP value for egressing traffic. Default:` unspecified`. Valid values: `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0`, `cs1`, `cs2`, `cs3`, `cs4`, `cs5`, `cs6`, `cs7`, `expedited_forwarding`, `voice_admit`, `unspecified`.
  * `target_cos` - (Optional) The target CoS traffic type for egressing traffic. Default: `unspecified`. Valid values: `background` (maps to CoS 0 - Background traffic), `best_effort` (maps to CoS 1 - Best effort traffic), `excellent_effort` (maps to CoS 2 - Excellent effort traffic), `critical_applications` (maps to CoS 3 - Critical applications traffic), `video` (maps to CoS 4 - Video traffic), `voice` (maps to CoS 5 - Voice traffic), `internetwork_control` (maps to CoS 6 - Internetwork control traffic), `network_control` (maps to CoS 7 - Network control traffic), `unspecified`.
  * `qos_priority` - (Optional) The QoS priority level to which the DSCP values will be mapped. Default: `unspecified`. Valid values: `level1`, `level2`, `level3`, `level4`, `level5`, `level6`, `unspecified`.
* `cos_mappings` - (Optional) A set of CoS (Class of Service) mappings. Each mapping defines how 802.1p CoS values are mapped to DSCP values and priority levels. Multiple CoS mappings can be configured.
  * `dot1p_from` - (Optional) The starting CoS traffic type of the range. Default: unspecified (when unset during creation). Valid values: `background` (maps to CoS 0 - Background traffic), `best_effort` (maps to CoS 1 - Best effort traffic), `excellent_effort` (maps to CoS 2 - Excellent effort traffic), `critical_applications` (maps to CoS 3 - Critical applications traffic), `video` (maps to CoS 4 - Video traffic), `voice` (maps to CoS 5 - Voice traffic), `internetwork_control` (maps to CoS 6 - Internetwork control traffic), `network_control` (maps to CoS 7 - Network control traffic), `unspecified`.
  dot1p_to - (Optional) The ending CoS traffic type of the range. Default: `unspecified`. Valid values: `background` (maps to CoS 0 - Background traffic), `best_effort` (maps to CoS 1 - Best effort traffic), `excellent_effort` (maps to CoS 2 - Excellent effort traffic), `critical_applications` (maps to CoS 3 - Critical applications traffic), `video` (maps to CoS 4 - Video traffic), `voice` (maps to CoS 5 - Voice traffic), `internetwork_control` (maps to CoS 6 - Internetwork control traffic), `network_control` (maps to CoS 7 - Network control traffic), `unspecified`.
  * `dscp_target` - (Optional) The target DSCP value for egressing traffic. Default: `unspecified`. Valid values: `af11`, `af12`, `af13`, `af21`, `af22`, `af23`, `af31`, `af32`, `af33`, `af41`, `af42`, `af43`, `cs0`, `cs1`, `cs2`, `cs3`, `cs4`, `cs5`, `cs6`, `cs7`, `expedited_forwarding`, `voice_admit`, `unspecified`.
  * `target_cos` - (Optional) The target CoS traffic type for egressing traffic. Default: `unspecified`. Valid values: `background` (maps to CoS 0 - Background traffic), `best_effort` (maps to CoS 1 - Best effort traffic), `excellent_effort` (maps to CoS 2 - Excellent effort traffic), `critical_applications` (maps to CoS 3 - Critical applications traffic), `video` (maps to CoS 4 - Video traffic), `voice` (maps to CoS 5 - Voice traffic), `internetwork_control` (maps to CoS 6 - Internetwork control traffic), `network_control` (maps to CoS 7 - Network control traffic), `unspecified`.
  * `qos_priority` - (Optional) The QoS priority level to which the CoS values will be mapped. Default: `unspecified`. Valid values: `level1`, `level2`, `level3`, `level4`, `level5`, `level6`, `unspecified`.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Custom QoS Policy.
* `id` - (Read-Only) The unique terraform identifier of the Custom QoS Policy in the template.

## Importing ##

An existing MSO Custom QoS Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_custom_qos_policy.qos_policy templateId/{template_id}/CustomQoSPolicy/{name}
```
