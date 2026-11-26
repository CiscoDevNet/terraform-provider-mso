---
layout: "mso"
page_title: "MSO Schema Template BD DHCP Policy"
sidebar_current: "docs-mso-resource-schema-template-bd-dhcp-policy"
description: |-
  Manages a DHCP policy association on a Bridge Domain in a schema template.
---

# mso_schema_template_bd_dhcp_policy

Manages one DHCP Relay Policy association on a Bridge Domain. Choose one way to manage DHCP policies for each Bridge Domain: configure the full list in the `dhcp_policies` block on `mso_schema_template_bd`, or use this resource for each policy. Mixing these methods for one Bridge Domain can make Terraform plan to remove policies managed by the other method.

For steps to switch between the two methods while retaining the existing associations, see [Migrating DHCP Policy Management](https://registry.terraform.io/providers/CiscoDevNet/mso/latest/docs/resources/schema_template_bd#migrating-dhcp-policy-management).

This resource is supported in NDO v4.0 and higher. The DHCP Relay Policy and optional DHCP Option Policy must exist in the tenant policy template. When a policy name is configured, the provider looks up its UUID; when a UUID is configured, the provider sends it directly to NDO.

## Example Usage

```hcl
resource "mso_schema_template_bd_dhcp_policy" "example" {
  schema_id     = mso_schema_template_bd.example.schema_id
  template_name = mso_schema_template_bd.example.template_name
  bd_name       = mso_schema_template_bd.example.name
  name                    = mso_tenant_policies_dhcp_relay_policy.example.name
  dhcp_option_policy_name = mso_tenant_policies_dhcp_option_policy.example.name
}
```

To configure the same association by UUID, use these arguments instead of `name` and `dhcp_option_policy_name`:

```hcl
resource "mso_schema_template_bd_dhcp_policy" "example" {
  schema_id               = mso_schema_template_bd.example.schema_id
  template_name           = mso_schema_template_bd.example.template_name
  bd_name                 = mso_schema_template_bd.example.name
  uuid                    = mso_tenant_policies_dhcp_relay_policy.example.uuid
  dhcp_option_policy_uuid = mso_tenant_policies_dhcp_option_policy.example.uuid
}
```

## Argument Reference

* `schema_id` - (Required) ID of the schema containing the Bridge Domain. Changing this forces a new resource.
* `template_name` - (Required) Name of the template containing the Bridge Domain. Changing this forces a new resource.
* `bd_name` - (Required) Name of the Bridge Domain. Changing this forces a new resource.
* `name` - (Optional) Name of the DHCP Relay Policy to associate. Mutually exclusive with `uuid`. Either `name` or `uuid` must be specified. When set, the provider looks up the policy in the tenant policy template and sends its UUID to NDO. When `uuid` is set, this is populated from the read response. Changing the configured identifier forces a new resource.
* `uuid` - (Optional) UUID of the DHCP Relay Policy to associate. Mutually exclusive with `name`. Either `name` or `uuid` must be specified. When set, the UUID is sent directly to NDO without a provider-side lookup. When `name` is set, this is populated from the read response. Changing the configured identifier forces a new resource.
* `dhcp_option_policy_name` - (Optional) Name of the DHCP Option Policy to associate. Mutually exclusive with `dhcp_option_policy_uuid`. When set, the provider looks up the policy in the tenant policy template and sends its UUID to NDO.
* `dhcp_option_policy_uuid` - (Optional) UUID of the DHCP Option Policy to associate. Mutually exclusive with `dhcp_option_policy_name`. When set, the UUID is sent directly to NDO without a provider-side lookup.

## Import

Import using `schema_id/template_name/bd_name/policy_name`:

```shell
terraform import mso_schema_template_bd_dhcp_policy.example schema_id/template_name/bridge_domain/dhcp_relay_policy
```
