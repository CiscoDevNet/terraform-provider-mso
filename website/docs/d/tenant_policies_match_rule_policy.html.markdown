---
layout: "mso"
page_title: "MSO: mso_tenant_policies_match_rule_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_match_rule_policy"
description: |-
  Reads Match Rule Policies on Cisco Nexus Dashboard Orchestrator (NDO).
---

# mso_tenant_policies_match_rule_policy #

Reads a Match Rule Policy from a tenant policy template on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v5.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Match Rule Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_match_rule_policy" "example" {
  template_id = mso_template.tenant.id
  name        = "example_match_rule_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template.
* `name` - (Required) The name of the Match Rule Policy.

## Attribute Reference ##

* `id` - The Terraform identifier of the Match Rule Policy in the template.
* `uuid` - The NDO UUID of the Match Rule Policy.
* `description` - The description of the Match Rule Policy.
* `match_prefix` - A set of match prefix blocks.
  * `prefix` - The IPv4 or IPv6 CIDR prefix.
  * `description` - The description of the match prefix.
  * `aggregate` - Whether to match a range of prefix lengths.
  * `from_prefix` - The minimum prefix length.
  * `to_prefix` - The maximum prefix length.
* `match_community_term` - A set of community term blocks.
  * `name` - The name of the community term.
  * `description` - The description of the community term.
  * `match_community_factor` - A set of community factor blocks within the term.
    * `community` - The regular, extended, or well-known community to match.
    * `scope` - The community scope.
    * `description` - The description of the community factor.
