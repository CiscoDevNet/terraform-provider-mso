---
layout: "mso"
page_title: "MSO: mso_tenant_policies_dhcp_relay_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_dhcp_relay_policy"
description: |-
  Data source for DHCP Relay Policy.
---

# mso_tenant_policies_dhcp_relay_policy #

Data source for DHCP Relay Policy.

## Example Usage ##

```hcl
data "mso_tenant_policies_dhcp_relay_policy" "dhcp_relay_policy" {
  template_id = mso_template.tenant_policy_template.id
  name        = "dhcp_relay_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the DHCP relay policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The UUID of the DHCP relay policy.
* `description` - (Read-Only) The description of the DHCP relay policy.
* `providers` - (Read-Only) A list of providers for the DHCP relay policy.
  * `dhcp_server_address` - (Read-Only) The DHCP server IP address of the provider.
  * `application_epg` - (Read-Only) The UUID of the Application Profile EPG.
  * `external_epg` - (Read-Only) The UUID of the External EPG.
  * `dhcp_server_vrf_preference` - (Read-Only) Indicates whether the server VRF is used.
