---
layout: "mso"
page_title: "MSO Schema Template BD DHCP Policy Data Source"
sidebar_current: "docs-mso-datasource-schema-template-bd-dhcp-policy"
description: |-
  Gets a DHCP policy association from a Bridge Domain in a schema template.
---

# mso_schema_template_bd_dhcp_policy

Use this data source to read a DHCP Relay Policy association from a Bridge Domain.

This data source is supported in NDO v4.0 and higher.

## Example Usage

```hcl
data "mso_schema_template_bd_dhcp_policy" "example" {
  schema_id     = mso_schema_template_bd.example.schema_id
  template_name = mso_schema_template_bd.example.template_name
  bd_name       = mso_schema_template_bd.example.name
  name          = mso_tenant_policies_dhcp_relay_policy.example.name
}
```

## Argument Reference

* `schema_id` - (Required) ID of the schema containing the Bridge Domain.
* `template_name` - (Required) Name of the template containing the Bridge Domain.
* `bd_name` - (Required) Name of the Bridge Domain.
* `name` - (Optional) Name of the DHCP Relay Policy association. Mutually exclusive with `uuid`. Either `name` or `uuid` must be specified.
* `uuid` - (Optional) UUID of the DHCP Relay Policy association. Mutually exclusive with `name`. Either `name` or `uuid` must be specified.

## Attribute Reference

* `name` - Resolved name of the associated DHCP Relay Policy.
* `uuid` - UUID of the associated DHCP Relay Policy.
* `dhcp_option_policy_name` - Name of the associated DHCP Option Policy.
* `dhcp_option_policy_uuid` - UUID of the associated DHCP Option Policy.

To look up the association by UUID, replace the `name` argument with `uuid`:

```hcl
data "mso_schema_template_bd_dhcp_policy" "example" {
  schema_id     = mso_schema_template_bd.example.schema_id
  template_name = mso_schema_template_bd.example.template_name
  bd_name       = mso_schema_template_bd.example.name
  uuid          = mso_tenant_policies_dhcp_relay_policy.example.uuid
}
```
