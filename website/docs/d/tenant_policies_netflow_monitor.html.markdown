---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_monitor"
sidebar_current: "docs-mso-data-source-tenant_policies_netflow_monitor"
description: |-
  Data source for NetFlow Monitor Policies.
---

# mso_tenant_policies_netflow_monitor #

Data source for NetFlow Monitor Policies. This data source is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> NetFlow Monitor

## Example Usage ##

```hcl
data "mso_tenant_policies_netflow_monitor" "netflow_monitor" {
  template_id = mso_template.tenant_template.id
  name        = "netflow_monitor_1"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the NetFlow Monitor.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the NetFlow Monitor.
* `description` - (Read-Only) The description of the NetFlow Monitor.
* `netflow_record_uuid` - (Read-Only) The UUID of the associated NetFlow Record.
* `netflow_exporter_uuids` - (Read-Only) The list of UUIDs of the associated NetFlow Exporters.
* `id` - (Read-Only) The unique terraform identifier of the NetFlow Monitor in the template.
