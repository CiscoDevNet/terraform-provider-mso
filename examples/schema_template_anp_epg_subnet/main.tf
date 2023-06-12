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

data "mso_site" "demo_site" {
  name = "demo_site"
}

resource "mso_tenant" "demo_tenant" {
  name         = "demo_tenant"
  display_name = "demo_tenant"
}

resource "mso_schema" "demo_schema" {
  name = "demo_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = mso_tenant.demo_tenant.id
  }
}

resource "mso_schema_site" "demo_schema_site" {
  schema_id     = mso_schema.demo_schema.id
  template_name = one(mso_schema.demo_schema.template).name
  site_id       = data.mso_site.demo_site.id
}

resource "mso_schema_template_vrf" "demo_vrf" {
  schema_id              = mso_schema.demo_schema.id
  template               = one(mso_schema.demo_schema.template).name
  name                   = "demo_vrf"
  display_name           = "demo_vrf"
  ip_data_plane_learning = "enabled"
}

resource "mso_schema_template_bd" "demo_bd" {
  schema_id     = mso_schema.demo_schema.id
  template_name = one(mso_schema.demo_schema.template).name
  name          = "demo_bd"
  display_name  = "demo_bd"
  vrf_name      = mso_schema_template_vrf.demo_vrf.name
  arp_flooding  = true
}

resource "mso_schema_template_anp" "demo_ap" {
  schema_id    = mso_schema.demo_schema.id
  template     = one(mso_schema.demo_schema.template).name
  name         = "demo_ap"
  display_name = "demo_ap"
}

resource "mso_schema_template_anp_epg" "demo_epg" {
  schema_id     = mso_schema.demo_schema.id
  template_name = one(mso_schema.demo_schema.template).name
  anp_name      = mso_schema_template_anp.demo_ap.name
  name          = "demo_epg"
  display_name  = "demo_epg"
  bd_name       = mso_schema_template_bd.demo_bd.name
  vrf_name      = mso_schema_template_vrf.demo_vrf.name
}

resource "mso_schema_template_anp_epg_subnet" "subnet1" {
  schema_id          = mso_schema.demo_schema.id
  template           = one(mso_schema.demo_schema.template).name
  anp_name           = mso_schema_template_anp.demo_ap.name
  epg_name           = mso_schema_template_anp_epg.demo_epg.name
  ip                 = "1.1.1.1/32"
  scope              = "public"
  primary            = true
  no_default_gateway = true
}