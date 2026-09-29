---
layout: "mso"
page_title: "MSO: mso_tenant_policies_netflow_exporter_site"
sidebar_current: "docs-mso-data-source-tenant_policies_netflow_exporter_site"
description: |-
  Retrieves a site-specific NetFlow Exporter configuration on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_tenant_policies_netflow_exporter_site #

Retrieves the site-specific configuration of a template-level NetFlow Exporter on Cisco Nexus Dashboard Orchestrator (NDO).

## Example Usage ##

```hcl
data "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
  template_id           = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.template_id
  site_id               = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.site_id
  netflow_exporter_uuid = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.netflow_exporter_uuid
}
```

## Argument Reference ##

* `template_id` - (Required) The ID of the tenant policy template.
* `site_id` - (Required) The ID of the site on which the NetFlow Exporter is configured.
* `netflow_exporter_uuid` - (Required) The UUID of the template-level NetFlow Exporter.

## Attribute Reference ##

* `id` - The unique Terraform identifier in the form `templateId/{template_id}/site/{site_id}/NetflowExporter/{netflow_exporter_uuid}`.
* `description` - The site-specific description of the NetFlow Exporter.
* `source_ip_address_type` - The source IP address type.
* `source_ip` - The custom source IP address or prefix.
* `destination_port` - The destination port of the NetFlow Exporter.
* `destination_ip` - The IPv4 or IPv6 destination address.
* `dscp` - The QoS DSCP value.
* `application_epg` - The application EPG associated with the NetFlow Exporter, if configured.
    * `tenant_name` - The name of the tenant.
    * `anp_name` - The name of the application profile.
    * `epg_name` - The name of the EPG.
    * `vrf_name` - The name of the VRF.
    * `vrf_tenant_name` - The name of the tenant that contains the VRF.
* `external_epg` - The external EPG associated with the NetFlow Exporter, if configured.
    * `tenant_name` - The name of the tenant.
    * `l3out_name` - The name of the L3Out.
    * `external_epg_name` - The name of the external EPG.
    * `vrf_name` - The name of the VRF.
    * `vrf_tenant_name` - The name of the tenant that contains the VRF.
