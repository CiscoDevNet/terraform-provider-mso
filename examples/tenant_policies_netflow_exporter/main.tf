terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO pwd>
  url      = "" # <MSO URL>
  insecure = true
}

data "mso_tenant" "example_tenant" {
  name = "example_tenant"
}

# Tenant template example

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id     = data.mso_tenant.example_tenant.id
  sites         = [data.mso_site.example_site.id]
}

# NetFlow Exporter example

resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
  template_id = mso_template.tenant_template.id
  name        = "netflow_exporter_1"
}

data "mso_site" "example_site" {
  name = "example_site"
}

# Site-specific NetFlow Exporter example. The site must already be associated
# with the tenant template and the referenced EPG objects must exist on the site.

resource "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
  template_id           = mso_tenant_policies_netflow_exporter.netflow_exporter.template_id
  site_id               = data.mso_site.example_site.id
  netflow_exporter_uuid = mso_tenant_policies_netflow_exporter.netflow_exporter.uuid

  source_ip_address_type = "ptep"
  destination_port       = "10"
  destination_ip         = "10.0.0.15"
  dscp                   = "cs1"

  # Omit vrf_tenant_name when the VRF is in the same tenant as the EPG.
  # Set it explicitly, for example to "common", when the VRF is in another tenant.
  application_epg {
    tenant_name = "example_tenant"
    anp_name    = "example_anp"
    epg_name    = "example_epg"
    vrf_name    = "example_vrf"
    # vrf_tenant_name = "common"
  }
}

# Site-specific NetFlow Exporter data source example

data "mso_tenant_policies_netflow_exporter_site" "netflow_exporter_site" {
  template_id           = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.template_id
  site_id               = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.site_id
  netflow_exporter_uuid = mso_tenant_policies_netflow_exporter_site.netflow_exporter_site.netflow_exporter_uuid
}

# NetFlow Exporter data source example

data "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
  template_id = mso_template.tenant_template.id
  name        = mso_tenant_policies_netflow_exporter.netflow_exporter.name
}
