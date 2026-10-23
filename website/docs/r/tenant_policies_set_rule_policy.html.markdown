---
layout: "mso"
page_title: "MSO: mso_tenant_policies_set_rule_policy"
sidebar_current: "docs-mso-resource-tenant_policies_set_rule_policy"
description: |-
  Manages Set Rule Policies on Cisco Nexus Dashboard Orchestrator (NDO).
---

# mso_tenant_policies_set_rule_policy #

Manages Set Rule Policies in a tenant policy template on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v5.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Set Rule Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_set_rule_policy" "example" {
  template_id = mso_template.tenant.id
  name        = "example_set_rule_policy"
  description = "Set attributes on matching routes"

  set_route_tag {
    value = 500
  }

  set_weight {
    value = 50
  }

  set_preference {
    value = 200
  }

  set_metric = 100

  set_next_hop    = "192.0.2.1"
  set_metric_type = "ospf-type2"

  set_next_hop_propagate = false
  set_multi_path         = false

  set_community {
    community = "no-export"
    criteria  = "replace"
  }

  set_additional_community = ["no-advertise", "extended:color:35"]

  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }

  set_as_path {
    prepend {
      path_asn {
        asn   = 10
        order = 1
      }

      path_asn {
        asn   = 20
        order = 2
      }
    }

    prepend_last_as {
      asn_count = 10
    }
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template. Changing this forces a new resource.
* `name` - (Required) The name of the Set Rule Policy. Length: 1-64 characters.
* `description` - (Optional) The description of the Set Rule Policy. Default: empty string when unset.
* `set_route_tag` - (Optional) A single block containing required integer `value`. Range: 0-999999999. Omit the block to remove the action.
* `set_weight` - (Optional) A single block containing required integer `value`. Range: 0-9999. Omit the block to remove the action.
* `set_preference` - (Optional) A single block containing required integer `value`. Range: 0-999999999. Omit the block to remove the action.
* `set_metric` - (Optional) The metric value. Range: 1-999999999. Omit to remove the action.
* `set_next_hop` - (Optional) The next-hop IPv4 or IPv6 address. Omit or set to `null` to remove the action.
* `set_metric_type` - (Optional) The metric type. Valid values: `none`, `ospf-type1`, `ospf-type2`. Omit or set to `null` to remove the action.
* `set_next_hop_propagate` - (Optional) Enable next-hop propagation. Default: `false` when unset. Cannot be enabled with a configured `set_route_tag` block. Omit or set to `false` to remove the action.
* `set_multi_path` - (Optional) Enable multipath. Default: `false` when unset. Requires `set_next_hop_propagate = true`. Omit or set to `false` to remove the action.
* `set_community` - (Optional) A single community action block. Omit the block to remove the action.
  * `criteria` - (Required) The community criteria. Valid values: `none`, `append`, `replace`.
  * `community` - (Optional) The community value. Required when `criteria` is `append` or `replace`.
* `set_additional_community` - (Optional) A set of nonempty community strings to append. Omit or set to `[]` to remove the additional communities.
* `set_dampening` - (Optional) A single route dampening block. Omit the block to remove the action.
  * `half_life` - (Required) The half-life time in minutes. Range: 1-60.
  * `reuse_limit` - (Required) The reuse limit. Range: 1-20000.
  * `suppress_limit` - (Required) The suppress limit. Must be larger than `reuse_limit`. Range: 1-20000.
  * `max_suppress_time` - (Required) The maximum suppress time in minutes. Range: 1-255.
* `set_as_path` - (Optional) A single AS-path action block. At least one of `prepend` or `prepend_last_as` must be configured. Omit the block to remove the AS-path actions.
  * `prepend` - (Optional) A single AS-path prepend block. An empty block is supported. Omit the block to remove the action.
    * `path_asn` - (Optional) A set of AS numbers to prepend. Omit to clear the configured AS numbers.
      * `asn` - (Required) The AS number. Range: 1-4294967295. The same ASN can appear with different `order` values.
      * `order` - (Required) The prepend sequence order. Range: 0-31. Each order must be unique.
  * `prepend_last_as` - (Optional) A single block for prepending the last AS. Omit the block to remove the action.
    * `asn_count` - (Required) The number of times to prepend the last AS. Range: 1-10.

## Attribute Reference ##

* `uuid` - The NDO UUID of the Set Rule Policy.
* `id` - The Terraform identifier of the Set Rule Policy in the template.

## Importing ##

An existing Set Rule Policy can be [imported](https://developer.hashicorp.com/terraform/cli/import) using its template ID and name:

```bash
terraform import mso_tenant_policies_set_rule_policy.example templateId/{template_id}/SetRulePolicy/{name}
```
