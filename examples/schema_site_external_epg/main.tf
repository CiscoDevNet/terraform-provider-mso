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
  platform = "nd"
}

data "mso_tenant" "demo_tenant" {
  name         = "demo_tenant"
  display_name = "demo_tenant"
}

data "mso_site" "demo_site" {
  name = "demo_site"
}

resource "mso_schema" "schema_blocks" {
  name = "demo_schema_blocks"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = data.mso_tenant.demo_tenant.id
  }
}

resource "mso_schema_site" "demo_schema_site" {
  schema_id     = mso_schema.schema_blocks.id
  template_name = one(mso_schema.schema_blocks.template).name
  site_id       = data.mso_site.demo_site.id
}

resource "mso_schema_template_vrf" "vrf1" {
  schema_id              = mso_schema.schema_blocks.id
  template               = one(mso_schema.schema_blocks.template).name
  name                   = "vrf1"
  display_name           = "vrf1"
  ip_data_plane_learning = "enabled"
}

resource "mso_schema_template_l3out" "template_l3out" {
  schema_id         = mso_schema.schema_blocks.id
  template_name     = one(mso_schema.schema_blocks.template).name
  l3out_name        = "l3out1"
  display_name      = "l3out1"
  vrf_name          = mso_schema_template_vrf.vrf1.name
  vrf_schema_id     = mso_schema_template_vrf.vrf1.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf1.template
}

resource "mso_schema_template_external_epg" "extepg1" {
  schema_id         = mso_schema.schema_blocks.id
  template_name     = one(mso_schema.schema_blocks.template).name
  external_epg_name = "extepg1"
  display_name      = "extepg1"
  vrf_name          = mso_schema_template_vrf.vrf1.name
  vrf_schema_id     = mso_schema_template_vrf.vrf1.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf1.template
  l3out_name        = mso_schema_template_l3out.template_l3out.l3out_name
}

resource "mso_schema_site_external_epg" "site_extepg1" {
  site_id           = mso_schema_site.demo_schema_site.id
  schema_id         = mso_schema.schema_blocks.id
  template_name     = one(mso_schema.schema_blocks.template).name
  external_epg_name = mso_schema_template_external_epg.extepg1.external_epg_name
  l3out_name        = mso_schema_template_l3out.template_l3out.l3out_name
}

resource "mso_schema_template_external_epg" "extepg2" {
  schema_id         = mso_schema.schema_blocks.id
  template_name     = one(mso_schema.schema_blocks.template).name
  external_epg_name = "extepg2"
  display_name      = "extepg2"
  vrf_name          = mso_schema_template_vrf.vrf1.name
  vrf_schema_id     = mso_schema_template_vrf.vrf1.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf1.template
}

resource "mso_schema_site_external_epg" "site_extepg2" {
  site_id           = mso_schema_site.demo_schema_site.id
  schema_id         = mso_schema.schema_blocks.id
  template_name     = one(mso_schema.schema_blocks.template).name
  external_epg_name = mso_schema_template_external_epg.extepg2.external_epg_name
  l3out_name        = "L3out2"
  l3out_on_apic     = true
}