---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_exporter"
sidebar_current: "docs-mso-resource-tenant_policies_netflow_exporter"
description: |-
  Manages NetFlow Exporter Policies on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_netflow_exporter #

Manages NetFlow Exporter Policies on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO v4.1 and higher.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Tenant Policies -> NetFlow Exporter

## Example Usage ##

```hcl
resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
  template_id = mso_template.tenant_template.id
  name        = "netflow_exporter_1"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the tenant policy template.
* `name` - (Required) The name of the NetFlow Exporter.

## Attribute Reference ##

* `uuid` - The NDO UUID of the NetFlow Exporter.
* `id` - The unique terraform identifier of the NetFlow Exporter in the template.

## Importing ##

An existing MSO NetFlow Exporter can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html

```bash
terraform import mso_tenant_policies_netflow_exporter.netflow_exporter templateId/{template_id}/NetflowExporter/{name}
```
