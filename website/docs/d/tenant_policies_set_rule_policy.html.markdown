---
layout: "mso"
page_title: "MSO: mso_tenant_policies_set_rule_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_set_rule_policy"
description: |-
  Reads Set Rule Policies on Cisco Nexus Dashboard Orchestrator (NDO).
---

# mso_tenant_policies_set_rule_policy #

Reads a Set Rule Policy from a tenant policy template on Cisco Nexus Dashboard Orchestrator (NDO). This data source is supported in NDO v5.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> Set Rule Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_set_rule_policy" "example" {
  template_id = mso_template.tenant.id
  name        = "example_set_rule_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template.
* `name` - (Required) The name of the Set Rule Policy.

## Attribute Reference ##

* `id` - The Terraform identifier of the Set Rule Policy in the template.
* `uuid` - The NDO UUID of the Set Rule Policy.
* `description` - The description of the Set Rule Policy.
* `set_route_tag` - A single route tag action block.
  * `value` - The route tag value.
* `set_weight` - A single weight action block.
  * `value` - The weight value.
* `set_preference` - A single preference action block.
  * `value` - The preference value.
* `set_metric` - The metric value. Returns `0` when the action is not configured.
* `set_next_hop` - The next-hop IPv4 or IPv6 address.
* `set_metric_type` - The metric type.
* `set_next_hop_propagate` - Whether next-hop propagation is enabled.
* `set_multi_path` - Whether multipath is enabled.
* `set_community` - A single community action block.
  * `community` - The community value.
  * `criteria` - The community criteria.
* `set_additional_community` - The set of community strings to append.
* `set_dampening` - A single route dampening action block.
  * `half_life` - The half-life time in minutes.
  * `reuse_limit` - The reuse limit.
  * `suppress_limit` - The suppress limit.
  * `max_suppress_time` - The maximum suppress time in minutes.
* `set_as_path` - A single AS-path action block.
  * `prepend` - A single block for prepending configured AS numbers.
    * `path_asn` - A set of AS number blocks.
      * `asn` - The AS number. The same ASN can appear with different `order` values.
      * `order` - The order determining the prepend sequence.
  * `prepend_last_as` - A single block for prepending the last AS number.
    * `asn_count` - The number of times to prepend the last AS.
