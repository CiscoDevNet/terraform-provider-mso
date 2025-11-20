---
layout: "mso"
page_title: "MSO: mso_custom_qos_policy"
sidebar_current: "docs-mso-data-source-custom_qos_policy"
description: |-
  Data source for Custom QoS Policy.
---

# mso_custom_qos_policy #

Data source for Custom Quality of Service (QoS) Policy. This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Custom QoS Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_custom_qos_policy" "qos_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_custom_qos_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the Custom QoS Policy to retrieve.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Custom QoS Policy.
* `id` - (Read-Only) The unique terraform identifier of the Custom QoS Policy in the template.
* `description` - (Read-Only) The description of the Custom QoS Policy.
* `dscp_mappings` - (Read-Only) A list of DSCP (Differentiated Services Code Point) mappings.
  * `dscp_from` - (Read-Only) The starting DSCP value of the range.
  * `dscp_to` - (Read-Only) The ending DSCP value of the range.
  * `dscp_target` - (Read-Only) The target DSCP value for egressing traffic.
  * `target_cos` - (Read-Only) The target CoS traffic type for egressing traffic.
  * `qos_priority` - (Read-Only) The QoS priority level.
* `cos_mappings` - (Read-Only) A list of CoS (Class of Service) mappings.
  * `dot1p_from` - (Read-Only) The starting CoS traffic type of the range.
  * `dot1p_to` - (Read-Only) The ending CoS traffic type of the range.
  * `dscp_target` - (Read-Only) The target DSCP value for egressing traffic.
  * `target_cos` - (Read-Only) The target CoS traffic type for egressing traffic.
  * `qos_priority` - (Read-Only) The QoS priority level.
