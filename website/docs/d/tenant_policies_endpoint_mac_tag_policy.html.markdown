---
layout: "mso"
page_title: "MSO: mso_tenant_policies_endpoint_mac_tag_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_endpoint_mac_tag_policy"
description: |-
  Data source for Endpoint MAC Tag Policy on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_endpoint_mac_tag_policy #

Data source for Endpoint MAC Tag Policy on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Endpoint MAC Tag Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_bd" {
  template_id = mso_template.tf_tenant_template.id
  mac         = "AA:BB:A1:B2:C3:D4"
  bd_uuid     = mso_schema_template_bd.tf_bd.uuid
}
```

```hcl
data "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_vrf" {
  template_id = mso_template.tf_tenant_template.id
  mac         = "AA:BB:A1:B2:C3:D5"
  vrf_uuid    = mso_schema_template_vrf.tf_vrf.uuid
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `mac` - (Required) The MAC address of the Endpoint MAC Tag Policy.
* `bd_uuid` - (Optional) The UUID of the Bridge Domain (BD) associated with the Endpoint MAC Tag Policy. Mutually exclusive with `vrf_uuid`. Either `bd_uuid` or `vrf_uuid` must be specified.
* `vrf_uuid` - (Optional) The UUID of the Virtual Routing and Forwarding (VRF) associated with the Endpoint MAC Tag Policy. Mutually exclusive with `bd_uuid`. Either `bd_uuid` or `vrf_uuid` must be specified.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Endpoint MAC Tag Policy.
* `uuid` - (Read-Only) The NDO UUID of the Endpoint MAC Tag Policy.
* `name` - (Read-Only) The name of the Endpoint MAC Tag Policy as assigned by NDO.
* `tag_annotations` - (Read-Only) A list of annotation key-value pairs for the Endpoint MAC Tag Policy.
  * `key` - (Read-Only) The annotation key.
  * `value` - (Read-Only) The annotation value.
* `policy_tags` - (Read-Only) A list of policy tag key-value pairs for the Endpoint MAC Tag Policy.
  * `key` - (Read-Only) The policy tag key.
  * `value` - (Read-Only) The policy tag value.
