---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_exporter"
sidebar_current: "docs-mso-data-source-tenant_policies_netflow_exporter"
description: |-
  Data source for NetFlow Exporter Policies.
---

# mso_tenant_policies_netflow_exporter #

Data source for NetFlow Exporter Policies. This data source is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> NetFlow Exporter

## Example Usage ##

```hcl
data "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
  template_id = mso_template.tenant_template.id
  name        = "netflow_exporter_1"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the NetFlow Exporter.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the NetFlow Exporter.
* `id` - (Read-Only) The unique terraform identifier of the NetFlow Exporter in the template.
