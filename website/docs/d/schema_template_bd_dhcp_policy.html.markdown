---
layout: "mso"
page_title: "MSO: mso_schema_template_bd_dhcp_policy"
sidebar_current: "docs-mso-data-source-schema_template_bd_dhcp_policy"
description: |-
  Data source for MSO Schema Template Bridge Domain DHCP Policy.
---

# mso_schema_template_bd_dhcp_policy

Data source for MSO Schema Template Bridge Domain DHCP Policy.

## Example Usage

```hcl
data "mso_schema_template_bd_dhcp_policy" "exp" {
  schema_id           = mso_schema.schema.id
  template_name       = mso_schema.schema.template_name
  bd_name             = mso_schema_template_bd.bridge_domain.name
  name                = mso_dhcp_relay_policy.example.name
}
```

## Argument Reference

- `schema_id` - (Required) SchemaID under which you want to deploy Bridge Domain.
- `template_name` - (Required) Template where Bridge Domain to be created.
- `bd_name` - (Required) Name of Bridge Domain.
- `name` - (Required) Name of the DHCP Relay Policy.

## Attribute Reference

- `version` - Version of the BD DHCP Label.
- `dhcp_option_name` - Name of the DHCP Option Policy.
- `dhcp_option_version` - Version of the attached DHCP Option Policy
