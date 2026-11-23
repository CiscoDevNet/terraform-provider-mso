---
layout: "mso"
page_title: "MSO Schema Template BD DHCP Policy"
sidebar_current: "docs-mso-resource-schema-template-bd-dhcp-policy"
description: |-
  Manages a DHCP policy association on a Bridge Domain in a schema template.
---

# mso_schema_template_bd_dhcp_policy

Manages one DHCP Relay Policy association on a Bridge Domain. Choose one way to manage DHCP policies for each Bridge Domain: configure the full list in the `dhcp_policies` block on `mso_schema_template_bd`, or use this resource for each policy. Mixing these methods for one Bridge Domain can make Terraform plan to remove policies managed by the other method.

This resource is supported in NDO v4.0 and higher. The DHCP Relay Policy and optional DHCP Option Policy must exist in the tenant policy template. The provider resolves their names to references when associating them with the Bridge Domain.

## Example Usage

```hcl
resource "mso_schema_template_bd_dhcp_policy" "example" {
  schema_id     = mso_schema_template_bd.example.schema_id
  template_name = mso_schema_template_bd.example.template_name
  bd_name       = mso_schema_template_bd.example.name
  name          = mso_tenant_policies_dhcp_relay_policy.example.name
}
```

## Argument Reference

* `schema_id` - (Required) ID of the schema containing the Bridge Domain. Changing this forces a new resource.
* `template_name` - (Required) Name of the template containing the Bridge Domain. Changing this forces a new resource.
* `bd_name` - (Required) Name of the Bridge Domain. Changing this forces a new resource.
* `name` - (Required) Name of the DHCP Relay Policy to associate. Changing this forces a new resource.
* `dhcp_option_policy_name` - (Optional) Name of the DHCP Option Policy to associate.

## Import

Import using `schema_id/template_name/bd_name/policy_name`:

```shell
terraform import mso_schema_template_bd_dhcp_policy.example schema_id/template_name/bridge_domain/dhcp_relay_policy
```
