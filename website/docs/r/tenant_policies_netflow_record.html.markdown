---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_record"
sidebar_current: "docs-mso-resource-tenant_policies_netflow_record"
description: |-
  Manages NetFlow Record Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_netflow_record #

Manages NetFlow Record Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> NetFlow Record

## Example Usage ##

```hcl
resource "mso_tenant_policies_netflow_record" "netflow_record" {
  template_id      = mso_template.tenant_template.id
  name             = "netflow_record_1"
  description      = "Test NetFlow Record"
  match_parameters = ["destination_ip", "source_ip", "ip_protocol"]
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the NetFlow Record. Maximum 64 characters.
* `description` - (Optional) The description of the NetFlow Record.
* `match_parameters` - (Optional) The set of match parameters for the NetFlow Record. Valid values are: `destination_ip`, `destination_ipv4`, `destination_ipv6`, `destination_mac`, `destination_port`, `ethertype`, `ip_protocol`, `source_ip`, `source_ipv4`, `source_ipv6`, `source_mac`, `source_port`.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the NetFlow Record.
* `id` - (Read-Only) The unique terraform identifier of the NetFlow Record in the template.

## Importing ##

An existing MSO NetFlow Record can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html

```bash
terraform import mso_tenant_policies_netflow_record.netflow_record templateId/{template_id}/NetflowRecord/{name}
```
