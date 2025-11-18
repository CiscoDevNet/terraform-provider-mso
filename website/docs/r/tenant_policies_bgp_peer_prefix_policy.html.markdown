---
layout: "mso"
page_title: "MSO: mso_tenant_policies_bgp_peer_prefix_policy"
sidebar_current: "docs-mso-resource-tenant_policies_bgp_peer_prefix_policy"
description: |-
  Manages BGP Peer Prefix Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_bgp_peer_prefix_policy #

Manages BGP Peer Prefix Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> BGP Peer Prefix Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
  template_id             = mso_template.template_tenant.id
  name                    = "test_bgp_peer_prefix_policy"
  description             = "Test BGP Peer Prefix Policy"
  action                  = "restart"
  max_number_of_prefixes  = 1000
  threshold_percentage    = 50
  restart_time            = 60
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the BGP Peer Prefix Policy.
* `description` - (Optional) The description of the BGP Peer Prefix Policy.
* `action` - (Optional) The action to take when the maximum number of prefixes is reached. Allowed values are log, reject, restart, shutdown.
* `max_number_of_prefixes` - (Optional) The maximum number of prefixes allowed for the BGP peer. Valid range: 1-300000.
* `threshold_percentage` - (Optional) The threshold percentage at which a warning is triggered. Valid range: 1-100.
* `restart_time` - (Optional) The time in seconds to wait before restarting the BGP session after reaching the maximum number of prefixes. Valid range: 1-65535. This parameter is only applicable when action is set to restart.

## Attribute Reference ##

* `uuid` - The NDO UUID of the BGP Peer Prefix Policy.
* `id` - The unique terraform identifier of the BGP Peer Prefix Policy in the template.

## Importing ##

An existing MSO BGP Peer Prefix Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_tenant_policies_bgp_peer_prefix_policy.bgp_policy templateId/{template_id}/BGPPeerPrefixPolicy/{name}
```
