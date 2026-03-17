---
layout: "mso"
page_title: "MSO: mso_tenant_policies_endpoint_mac_tag_policy"
sidebar_current: "docs-mso-resource-tenant_policies_endpoint_mac_tag_policy"
description: |-
  Manages Endpoint MAC Tag Policy on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_endpoint_mac_tag_policy #

Manages Endpoint MAC Tag Policy on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Endpoint MAC Tag Policy

## Example Usage ##

**Endpoint MAC Tag Policy with a Bridge Domain scope:**

```hcl
resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_tag_bd" {
  template_id = mso_template.tf_tenant_template.id
  mac         = "AA:BB:A1:B2:C3:D4"
  bd_uuid     = mso_schema_template_bd.tf_bd.uuid

  tag_annotations {
    key   = "annotation_key_1"
    value = "annotation_value_1"
  }

  policy_tags {
    key   = "policy_key_1"
    value = "policy_value_1"
  }
}
```

**Endpoint MAC Tag Policy with a VRF scope:**

```hcl
resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_tag_vrf" {
  template_id = mso_template.tf_tenant_template.id
  mac         = "AA:BB:A1:B2:C3:D5"
  vrf_uuid    = mso_schema_template_vrf.tf_vrf.uuid

  tag_annotations {
    key   = "annotation_key_1"
    value = "annotation_value_1"
  }

  policy_tags {
    key   = "policy_key_1"
    value = "policy_value_1"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `mac` - (Required) The MAC address of the Endpoint MAC Tag Policy.
* `bd_uuid` - (Optional) The UUID of the Bridge Domain (BD) to associate with Endpoint MAC Tag Policy. Mutually exclusive with `vrf_uuid`. Either `bd_uuid` or `vrf_uuid` must be specified.
* `vrf_uuid` - (Optional) The UUID of the Virtual Routing and Forwarding (VRF) to associate with Endpoint MAC Tag Policy. Mutually exclusive with `bd_uuid`. Either `bd_uuid` or `vrf_uuid` must be specified.
* `tag_annotations` - (Optional) A list of annotation key-value pairs for the Endpoint MAC Tag Policy.
  * `key` - (Required) The annotation key.
  * `value` - (Required) The annotation value.
* `policy_tags` - (Optional) A list of policy tag key-value pairs for the Endpoint MAC Tag Policy.
  * `key` - (Required) The policy tag key.
  * `value` - (Required) The policy tag value.

## Attribute Reference ##

* `id` - (Read-Only) The unique Terraform identifier of the Endpoint MAC Tag Policy.
* `uuid` - (Read-Only) The NDO UUID of the Endpoint MAC Tag Policy.
* `name` - (Read-Only) The name of the Endpoint MAC Tag Policy.

## Importing ##

An existing MSO Endpoint MAC Tag Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html

```bash
terraform import mso_tenant_policies_endpoint_mac_tag_policy.endpoint_mac_tag_bd templateId/{template_id}/EndpointMACTagPolicy/{name}
```
