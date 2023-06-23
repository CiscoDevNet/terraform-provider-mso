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

resource "mso_schema_site_bd" "foo_schema_site_bd" {
  schema_id     = mso_schema.demo_schema.id
  bd_name       = mso_schema_template_bd.demo_bd.name
  template_name = one(mso_schema.demo_schema.template).name
  site_id       = data.mso_site.demo_site.id
  host_route    = false
  svi_mac       = "00:22:BD:F8:19:FF"
}