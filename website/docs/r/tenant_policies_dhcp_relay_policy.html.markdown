---
layout: "mso"
page_title: "MSO: mso_tenant_policies_dhcp_relay_policy"
sidebar_current: "docs-mso-resource-tenant_policies_dhcp_relay_policy"
description: |-
  Manages DHCP Relay Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_dhcp_relay_policy #

Manages DHCP Relay Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 or higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> DHCP Relay Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_dhcp_relay_policy" "dhcp_relay_policy" {
  name        = "dhcp_relay_policy"
  template_id = mso_template.tenant_policy_template.id
  description = "example_dhcp_relay_policy"
  dhcp_relay_providers {
    dhcp_server_address  = "1.1.1.1"
    application_epg_uuid = mso_schema_template_anp_epg.anp_epg.uuid
  }
  dhcp_relay_providers {
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
* `dhcp_relay_providers` - (Required) A list of providers for the DHCP relay policy.
  * `dhcp_server_address` - (Required) The DHCP server IP address of the provider.
  * `application_epg_uuid` - (Optional) The NDO UUID of the Application Profile EPG, `application_epg_uuid` is required only when `external_epg_uuid` is not set, as both cannot be used at the same time.
  * `external_epg_uuid` - (Optional) The NDO UUID of the External EPG, `external_epg_uuid` is required only when `application_epg_uuid` is not set, as both cannot be used at the same time.
  * `dhcp_server_vrf_preference` - (Optional) Enabling DHCP Server VRF Preference allows the switch to route DHCP relay packets from the server VRF, regardless of any contract between the client and server EPGs. Consequently, the server VRF requires at least one IP address on all leaf switches where client bridge domains are deployed.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the DHCP relay policy.
* `id` - (Read-Only) The unique terraform identifier of the DHCP relay policy.

## Importing ##

An existing MSO DHCP relay policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```
terraform import mso_tenant_policies_dhcp_relay_policy.dhcp_relay_policy templateId/{template_id}/DHCPRelayPolicy/{name}
```