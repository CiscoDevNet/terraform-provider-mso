---
layout: "mso"
page_title: "MSO: mso_schema_template_bd_dhcp_policy"
sidebar_current: "docs-mso-resource-schema_template_bd_dhcp_policy"
description: |-
  Manages MSO Schema Template Bridge Domain DHCP Policy.
---

# mso_schema_template_bd_dhcp_policy #

Manages MSO Schema Template Bridge Domain DHCP Policy.

## Example Usage ##

```hcl
resource "mso_schema_template_bd_dhcp_policy" "exp" {
  schema_id           = mso_schema.schema.id
  template_name       = mso_schema.schema.template_name
  bd_name             = mso_schema_template_bd.bridge_domain.name
  name                = mso_dhcp_relay_policy.example.name
  version             = 1
  dhcp_option_name    = mso_dhcp_option_policy.example.name
  dhcp_option_version = 1
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Bridge Domain Subnet.
* `template_name` - (Required) Template where Bridge Domain Subnet to be created.
* `bd_name` - (Required) Name of Bridge Domain.
* `name` - (Required) Name of the DHCP Relay Policy.
* `version` - (Optional) Version of the BD DHCP Label.
* `dhcp_option_name` - (Optional) Name of the DHCP Option Policy.
* `dhcp_option_version` - (Optional) Version of the attached DHCP Option Policy

### Note
 `dhcp_option_version` is required if `dhcp_option_name` is set.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Bridge Domain DHCP Policy can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_bd_dhcp_policy.bdsub1 /schemas/{schema_id}/templates/{template_name}/bds/{bd_name}/dhcpLabels/{dhcp_relay_policy_name}
```