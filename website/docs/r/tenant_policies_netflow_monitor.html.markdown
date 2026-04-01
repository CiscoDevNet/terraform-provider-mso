---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_monitor"
sidebar_current: "docs-mso-resource-tenant_policies_netflow_monitor"
description: |-
  Manages NetFlow Monitor Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_netflow_monitor #

Manages NetFlow Monitor Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> NetFlow Monitor

## Example Usage ##

```hcl
resource "mso_tenant_policies_netflow_monitor" "netflow_monitor" {
  template_id            = mso_template.tenant_template.id
  name                   = "netflow_monitor_1"
  description            = "Test NetFlow Monitor"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.netflow_record.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.netflow_exporter.uuid]
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the NetFlow Monitor. Maximum 64 characters.
* `description` - (Optional) The description of the NetFlow Monitor.
* `netflow_record_uuid` - (Optional) The UUID of the NetFlow Record to associate with this monitor.
* `netflow_exporter_uuids` - (Required) The list of UUIDs of the NetFlow Exporters to associate with this monitor. At least one exporter must be provided.

## Attribute Reference ##

* `uuid` - The NDO UUID of the NetFlow Monitor.
* `id` - The unique terraform identifier of the NetFlow Monitor in the template.

## Importing ##

An existing MSO NetFlow Monitor can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html

```bash
terraform import mso_tenant_policies_netflow_monitor.netflow_monitor templateId/{template_id}/NetflowMonitor/{name}
```
