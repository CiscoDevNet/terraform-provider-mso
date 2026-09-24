---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_exporter_site"
sidebar_current: "docs-mso-resource-tenant_policies_netflow_exporter_site"
description: |-
  Manages a site-specific NetFlow Exporter configuration on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_netflow_exporter_site #

Manages a site-specific configuration for a template-level NetFlow Exporter on Cisco Nexus Dashboard Orchestrator (NDO).

The referenced `mso_tenant_policies_netflow_exporter` must exist before this resource is created, and the selected site must be attached to the tenant policy template.

## Example Usage ##

```hcl
resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
  template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
  site_id               = data.mso_site.example_site.id
  netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

  source_ip_address_type = "ptep"
  destination_port = "10"
  destination_ip   = "10.0.0.15"
  dscp             = "cs1"

  external_epg {
    tenant_name       = "example_tenant"
    l3out_name        = "example_l3out"
    external_epg_name = "example_external_epg"
    vrf_name          = "default"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template. Changing this forces a new resource.
* `site_id` - (Required) The ID of the site on which to configure the exporter. Changing this forces a new resource.
* `netflow_exporter_uuid` - (Required) The UUID of the template-level NetFlow Exporter. Changing this forces a new resource.
* `description` - (Optional) The site-specific description of the NetFlow Exporter.
* `source_ip_address_type` - (Required) Source IP address type. Allowed values are `inbandMgmtIP`, `oobMgmtIP`, `ptep`, and `customSrcIP`.
* `source_ip` - (Optional) The custom source IP address or prefix. Required when `source_ip_address_type` is `customSrcIP`.
* `destination_port` - (Required) The destination port of the NetFlow Exporter.
* `destination_ip` - (Required) The IPv4 or IPv6 destination address.
* `dscp` - (Required) The QoS DSCP value.
* `application_epg` - (Optional) The application EPG associated with the NetFlow Exporter. Cannot be combined with `external_epg`.
    * `tenant_name` - (Required) The name of the tenant.
    * `anp_name` - (Required) The name of the application profile.
    * `epg_name` - (Required) The name of the application EPG.
    * `vrf_name` - (Required) The name of the VRF.
    * `vrf_tenant_name` - (Optional) The name of the tenant that contains the VRF. If omitted, the tenant name is used.
* `external_epg` - (Optional) The external EPG associated with the NetFlow Exporter. Cannot be combined with `application_epg`.
    * `tenant_name` - (Required) The name of the tenant.
    * `l3out_name` - (Required) The name of the L3Out.
    * `external_epg_name` - (Required) The name of the external EPG.
    * `vrf_name` - (Required) The name of the VRF.
    * `vrf_tenant_name` - (Optional) The name of the tenant that contains the VRF. If omitted, the tenant name is used.

## Attribute Reference ##

* `id` - The unique Terraform identifier in the form `templateId/{template_id}/site/{site_id}/NetflowExporter/{netflow_exporter_uuid}`.

## Importing ##

```bash
terraform import mso_tenant_policies_netflow_exporter_site.netflow_exporter_site templateId/{template_id}/site/{site_id}/NetflowExporter/{netflow_exporter_uuid}
```
