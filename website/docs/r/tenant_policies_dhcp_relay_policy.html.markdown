---
layout: "mso"
page_title: "MSO: mso_tenant_policies_dhcp_relay_policy"
sidebar_current: "docs-mso-resource-tenant_policies_dhcp_relay_policy"
description: |-
  Manages DHCP Relay Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_dhcp_relay_policy #

Manages DHCP Relay Policies on Cisco Nexus Dashboard Orchestrator (NDO)

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> DHCP Relay Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_dhcp_relay_policy" "dhcp_relay_policy" {
  name        = "dhcp_relay_policy"
  template_id = mso_template.tenant_policy_template.id
  description = "example_dhcp_relay_policy"
  providers {
    dhcp_server_address  = "1.1.1.1"
    application_epg_uuid = mso_schema_template_anp_epg.anp_epg.uuid
  }
  providers {
    dhcp_server_address        = "2.2.2.2"
    external_epg_uuid          = mso_schema_template_external_epg.ext_epg.uuid
    dhcp_server_vrf_preference = true
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the DHCP relay policy.
* `description` - (Optional) The description of the DHCP relay policy.
* `providers` - (Required) A list of providers for the DHCP relay policy.
  * `dhcp_server_address` - (Required) The DHCP server IP address of the provider.
  * `application_epg_uuid` - (Optional) The UUID of the Application Profile EPG.
  * `external_epg_uuid` - (Optional) The UUID of the External EPG.
  * `dhcp_server_vrf_preference` - (Optional) Indicates whether the server VRF is used.

## Attribute Reference ##

* `id` - The unique identifier of the DHCP relay policy in the template.
* `uuid` - The UUID of the DHCP relay policy.

## Importing ##

An existing MSO DHCP relay policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```
terraform import mso_tenant_policies_dhcp_relay_policy.dhcp_relay_policy templateId/{template_id}/DHCPRelayPolicy/{name}
```