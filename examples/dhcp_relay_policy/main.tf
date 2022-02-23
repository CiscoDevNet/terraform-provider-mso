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

resource "mso_tenant" "tenant1" {
  name         = var.tenant_name
  display_name = var.tenant_name
  description  = "DemoTenant"
}

resource "mso_schema" "schema1" {
  name          = var.schema_name
  template_name = var.template_name
  tenant_id     = data.mso_tenant.tenant1.id
}

resource "mso_schema_template_vrf" "vrf1" {
  schema_id    = mso_schema.schema1.id
  template     = mso_schema.schema1.template_name
  name         = var.vrf_name
  display_name = var.vrf_name
}

resource "mso_schema_template_external_epg" "epg1" {
  schema_id         = mso_schema.schema1.id
  template_name     = mso_schema.schema1.template_name
  external_epg_name = var.external_epg_name
  display_name      = var.external_epg_name
  vrf_name          = mso_schema_template_vrf.vrf1.name
}

resource "mso_dhcp_relay_policy" "dp1" {
  tenant_id   = mso_tenant.tenant1.id
  name        = var.relay_policy_name
  description = "desc"
  dhcp_relay_policy_provider {
    epg                 = var.epg
    dhcp_server_address = var.dhcp_server_address
  }
  dhcp_relay_policy_provider {
    external_epg        = mso_schema_template_external_epg.epg1.id
    dhcp_server_address = var.dhcp_server_address
  }
}
