---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_record"
sidebar_current: "docs-mso-data-source-tenant_policies_netflow_record"
description: |-
  Data source for NetFlow Record Policies.
---

# mso_tenant_policies_netflow_record #

Data source for NetFlow Record Policies. This data source is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> NetFlow Record

## Example Usage ##

```hcl
data "mso_tenant_policies_netflow_record" "netflow_record" {
  template_id = mso_template.tenant_template.id
  name        = "netflow_record_1"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the template.
* `name` - (Required) The name of the NetFlow Record.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the NetFlow Record.
* `description` - (Read-Only) The description of the NetFlow Record.
* `match_parameters` - (Read-Only) The match parameters of the NetFlow Record.
* `id` - (Read-Only) The unique terraform identifier of the NetFlow Record in the template.
