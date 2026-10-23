---
layout: "mso"
page_title: "MSO: mso_tenant_policies_match_rule_policy"
sidebar_current: "docs-mso-resource-tenant_policies_match_rule_policy"
description: |-
  Manages Match Rule Policies on Cisco Nexus Dashboard Orchestrator (NDO).
---

# mso_tenant_policies_match_rule_policy #

Manages Match Rule Policies in a tenant policy template on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v5.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Match Rule Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_match_rule_policy" "example" {
  template_id = mso_template.tenant.id
  name        = "example_match_rule_policy"
  description = "Match selected routes"

  match_prefix {
    prefix      = "10.10.10.0/24"
    description = "IPv4 routes"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 32
  }

  match_prefix {
    prefix      = "2001:db8::/64"
    description = "IPv6 routes"
    aggregate   = false
    from_prefix = 0
    to_prefix   = 0
  }

  match_community_term {
    name        = "export_community"
    description = "Match export communities"

    match_community_factor {
      community   = "no-export"
      scope       = "transitive"
      description = "Match no-export"
    }

    match_community_factor {
      community   = "no-advertise"
      scope       = "non-transitive"
      description = "Match no-advertise"
    }
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template. Changing this forces a new resource.
* `name` - (Required) The name of the Match Rule Policy. Length: 1-64 characters.
* `description` - (Optional) The description of the Match Rule Policy. Default: empty string when unset.
* `match_prefix` - (Optional) A set of match prefix blocks. Omit to remove the prefixes.
  * `prefix` - (Required) The IPv4 or IPv6 CIDR prefix. Prefixes must be unique within the policy.
  * `description` - (Optional) The description of the match prefix. Default: empty string when unset.
  * `aggregate` - (Optional) Whether to match a range of prefix lengths. Default: `false` when unset. When `false`, both `from_prefix` and `to_prefix` must be `0`.
  * `from_prefix` - (Optional) The minimum prefix length. Range: 0-32 for IPv4, 0-128 for IPv6. Default: `0` when unset. When `aggregate` is `true`, must be `0` or greater than the prefix mask, and less than or equal to `to_prefix`.
  * `to_prefix` - (Optional) The maximum prefix length. Range: 0-32 for IPv4, 0-128 for IPv6. Default: `0` when unset.
* `match_community_term` - (Optional) A set of community term blocks. Omit to remove the community terms.
  * `name` - (Required) The name of the community term. Length: 1-64 characters. Names must be unique within the policy.
  * `description` - (Optional) The description of the community term. Default: empty string when unset.
  * `match_community_factor` - (Optional) A set of community factor blocks. Omit to remove the community factors.
    * `community` - (Required) The regular, extended, or well-known community to match. Values must be nonempty and unique within the term.
    * `scope` - (Optional) The scope of the community factor. Valid values: `transitive`, `non-transitive`. Default: `transitive` when unset.
    * `description` - (Optional) The description of the community factor. Default: empty string when unset.

## Attribute Reference ##

* `uuid` - The NDO UUID of the Match Rule Policy.
* `id` - The Terraform identifier of the Match Rule Policy in the template.

## Importing ##

An existing Match Rule Policy can be [imported](https://developer.hashicorp.com/terraform/cli/import) using its template ID and name:

```bash
terraform import mso_tenant_policies_match_rule_policy.example templateId/{template_id}/MatchRulePolicy/{name}
```
