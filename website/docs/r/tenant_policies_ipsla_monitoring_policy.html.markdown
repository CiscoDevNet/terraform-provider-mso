---
layout: "mso"
page_title: "MSO: mso_tenant_policies_ipsla_monitoring_policy"
sidebar_current: "docs-mso-resource-tenant_policies_ipsla_monitoring_policy"
description: |-
Manages IPSLA Monitoring Policies on Cisco Nexus Dashboard Orchestrator (NDO).
---



# mso_tenant_policies_ipsla_monitoring_policy #

Manages Internet Protocol Service Level Agreement (IPSLA) Monitoring Policies on Cisco Nexus Dashboard Orchestrator (NDO).

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> IPSLA Monitoring Policy

## Example Usage ##

```hcl
resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
  template_id        = mso_template.template_tenant.id
  name               = "test_ipsla_policy"
  description        = "HTTP Type"
  sla_type           = "http"
  destination_port   = 80
  http_version       = "HTTP11"
  http_uri           = "/example"
  sla_frequency      = 120
  detect_multiplier  = 4
  request_data_size  = 64
  type_of_service    = 18
  operation_timeout  = 100
  threshold          = 100
  ipv6_traffic_class = 255
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the IPSLA monitoring policy.
* `description` - (Optional) The description of the IPSLA monitoring policy.
* `sla_type` - (Optional) The type of Service Level Agreement (SLA). Allowed values are `http`, `tcp`, `icmp`, `l2ping`.
* `destination_port` - (Optional) The destination port for the IPSLA. Valid range: 1-65535.
* `http_version` - (Optional) The HTTP version used for IPSLA. Allowed values are `HTTP10`, `HTTP11`.
* `http_uri` - (Optional) The URI used for HTTP IPSLA.
* `sla_frequency` - (Optional) The frequency of IPSLA monitoring in seconds. Valid range: 1-300.
* `detect_multiplier` - (Optional) The detection multiplier for IPSLA. Valid range: 1-100.
* `request_data_size` - (Optional) The size of the request data in bytes. Valid range: 1-17512.
* `type_of_service` - (Optional) The IPv4 Type of Service. Valid range: 0-255.
* `operation_timeout` - (Optional) The operation timeout for IPSLA in milliseconds. Valid range: 0-604800000.
* `threshold` - (Optional) The threshold for IPSLA in milliseconds. Valid range: 0-604800000.
* `ipv6_traffic_class` - (Optional) The IPv6 Traffic Class. Valid range: 0-255.

## Attribute Reference ##

* `uuid` - The UUID of the IPSLA monitoring policy.
* `id` - The unique identifier of the IPSLA monitoring policy in the template.

## Importing ##

An existing MSO IPSLA Monitoring Policy can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy templateId/{template_id}/IPSLAMonitoringPolicy/{name}
```
