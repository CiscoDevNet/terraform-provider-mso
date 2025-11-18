---
layout: "mso"
page_title: "MSO: mso_tenant_policies_bgp_peer_prefix_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_bgp_peer_prefix_policy"
description: |-
  Data source for BGP Peer Prefix Policy.
---

# mso_tenant_policies_bgp_peer_prefix_policy #

Data source for BGP Peer Prefix Policy. This data source is supported in NDO v4.3 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> BGP Peer Prefix Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_bgp_peer_prefix_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the BGP Peer Prefix Policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the BGP Peer Prefix Policy.
* `id` - (Read-Only) The unique terraform identifier of the BGP Peer Prefix Policy in the template.
* `description` - (Read-Only) The description of the BGP Peer Prefix Policy.
* `action` - (Read-Only) The action to take when the maximum number of prefixes is reached.
* `max_number_of_prefixes` - (Read-Only) The maximum number of prefixes allowed for the BGP peer.
* `threshold_percentage` - (Read-Only) The threshold percentage at which a warning is triggered.
* `restart_time` - (Read-Only) The time in seconds to wait before restarting the BGP session after reaching the maximum number of prefixes.
