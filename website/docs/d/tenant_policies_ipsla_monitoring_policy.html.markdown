---
layout: "mso"
page_title: "MSO: mso_tenant_policies_ipsla_monitoring_policy"
sidebar_current: "docs-mso-data-source-tenant_policies_ipsla_monitoring_policy"
description: |-
  Data source for IPSLA Monitoring Policy.
---



# mso_tenant_policies_ipsla_monitoring_policy #

Data source for Internet Protocol Service Level Agreement (IPSLA) Monitoring Policy.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> IPSLA Monitoring Policy

## Example Usage ##

```hcl
data "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
  template_id        = mso_template.template_tenant.id
  name               = "test_ipsla_policy"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the IPSLA monitoring policy.

## Attribute Reference ##

* `uuid` - (Read-Only) The UUID of the IPSLA monitoring policy.
* `id` - (Read-Only) The unique identifier of the IPSLA monitoring policy in the template.
* `description` - (Read-Only) The description of the IPSLA monitoring policy.
* `sla_type` - (Read-Only) The type of Service Level Agreement (SLA).
* `destination_port` - (Read-Only) The destination port for the IPSLA.
* `http_version` - (Read-Only) The HTTP version used for IPSLA.
* `http_uri` - (Read-Only) The URI used for HTTP IPSLA.
* `sla_frequency` - (Read-Only) The frequency of IPSLA monitoring in seconds.
* `detect_multiplier` - (Read-Only) The detection multiplier for IPSLA.
* `request_data_size` - (Read-Only) The size of the request data in bytes.
* `type_of_service` - (Read-Only) The IPv4 Type of Service.
* `operation_timeout` - (Read-Only) The operation timeout for IPSLA in milliseconds.
* `threshold` - (Read-Only) The threshold for IPSLA in milliseconds.
* `ipv6_traffic_class` - (Read-Only) The IPv6 Traffic Class.
