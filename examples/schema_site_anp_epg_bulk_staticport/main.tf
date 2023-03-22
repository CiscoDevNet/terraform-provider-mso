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

// Create multiple Static Ports.
resource "mso_site" "test_site" {
  name             = var.site_name
  username         = var.site_username
  password         = var.site_password
  apic_site_id     = 105
  urls             = ["https://10.23.248.102"]
  # login_domain     = "radius_test"
  # maintenance_mode = true
  location = {
    lat  = 78.946
    long = 95.623
  }
}

resource "mso_tenant" "tenant_test" {
  name = var.tenant_name
  display_name = var.tenant_name
  site_associations {
    site_id     = data.mso_site.test_site.id
  }
}

resource "mso_schema" "schema_test" {
  name = var.schema_name
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = data.mso_tenant.test_tenant.id
  }
}

resource "mso_schema_template_vrf" "vrf" {
  schema_id    = mso_schema.schema_test.id
  template     = tolist(mso_schema.schema_test.template)[0].name
  name         = var.vrf_name
  display_name = var.vrf_name
}

resource "mso_schema_template_bd" "bd" {
  schema_id         = mso_schema.schema_test.id
  template_name     = mso_schema_template_vrf.vrf.template
  name              = var.bd_name
  display_name      = var.bd_name
  vrf_name          = mso_schema_template_vrf.vrf.name
  vrf_schema_id     = mso_schema_template_vrf.vrf.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf.template
}

resource "mso_schema_template_anp" "anp" {
  schema_id    = mso_schema.schema_test.id
  template     = tolist(mso_schema.schema_test.template)[0].name
  name         = var.anp_name
  display_name = var.anp_name
}

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id     = mso_schema.schema_test.id
  template_name = mso_schema_template_anp.anp.template
  anp_name      = mso_schema_template_anp.anp.name
  name          = var.epg_name
  display_name  = var.epg_name
  bd_name       = mso_schema_template_bd.bd.name
  vrf_name      = mso_schema_template_vrf.vrf.name
}

resource "mso_schema_site" "schema_site" {
  schema_id     = mso_schema.schema_test.id
  site_id       = data.mso_site.test_site.id
  template_name = mso_schema_template_anp_epg.anp_epg.template_name
}

resource "mso_schema_site_anp" "site_anp" {
  schema_id     = mso_schema.schema_test.id
  template_name = mso_schema_site.schema_site.template_name
  site_id       = data.mso_site.test_site.id
  anp_name      = mso_schema_template_anp.anp.name
}

resource "mso_schema_site_anp_epg" "site_anp_epg" {
  schema_id     = mso_schema.schema_test.id
  template_name = mso_schema_site_anp.site_anp.template_name
  site_id       = data.mso_site.test_site.id
  anp_name      = mso_schema_site_anp.site_anp.anp_name
  epg_name      = mso_schema_template_anp_epg.anp_epg.name
}

resource "mso_schema_site_anp_epg_bulk_staticport" "bulk_static_port" {
  schema_id            = mso_schema.schema_test.id
  site_id              = data.mso_site.test_site.id
  template_name        = tolist(mso_schema.schema_test.template)[0].name
  anp_name             = var.anp_name
  epg_name             = mso_schema_site_anp_epg.site_anp_epg.epg_name
  static_ports {
    path_type            = "vpc"
    deployment_immediacy = "lazy"
    pod                  = "pod-4"
    leaf                 = "105"
    path                 = "eth1/4"
    vlan                 = 207
    mode                 = "regular"
  }
  static_ports {
    path_type            = "port"
    deployment_immediacy = "immediate"
    pod                  = "pod-1"
    leaf                 = "102"
    path                 = "eth1/11"
    vlan                 = 200
    micro_seg_vlan       = 3
    mode                 = "untagged"
  }
}

data "mso_schema_site_anp_epg_bulk_staticport" "example" {
  schema_id            = mso_schema_site_anp_epg_bulk_staticport.bulk_static_port.schema_id
  site_id              = mso_schema_site_anp_epg_bulk_staticport.bulk_static_port.site_id
  template_name        = mso_schema_site_anp_epg_bulk_staticport.bulk_static_port.template_name
  anp_name             = mso_schema_site_anp_epg_bulk_staticport.bulk_static_port.anp_name
  epg_name             = mso_schema_site_anp_epg_bulk_staticport.bulk_static_port.epg_name
}

output "bulk_static_ports" {
  value = data.mso_schema_site_anp_epg_bulk_staticport.example
}
