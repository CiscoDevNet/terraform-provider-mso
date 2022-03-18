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

data "mso_tenant" "common" {
  name         = "common"
  display_name = "common"
}

resource "mso_schema" "schema1" {
  name          = "tf_test_schema"
  template_name = "Template1"
  tenant_id     = data.mso_tenant.common.id
}

resource "mso_schema" "schema2" {
  name          = "tf_test_schema_2"
  template_name = "Template2"
  tenant_id     = data.mso_tenant.common.id
}

resource "mso_schema_template_vrf" "vrf1" {
  schema_id    = mso_schema.schema1.id
  template     = mso_schema.schema1.template_name
  name         = "vrf1"
  display_name = "vrf1"
}

resource "mso_schema_template_vrf" "vrf2" {
  schema_id    = mso_schema.schema2.id
  template     = mso_schema.schema2.template_name
  name         = "vrf2"
  display_name = "vrf2"
}

resource "mso_schema_template_bd" "bd1" {
  schema_id     = mso_schema.schema1.id
  template_name = mso_schema.schema1.template_name
  name          = "bd1"
  display_name  = "bd1"
  vrf_name      = mso_schema_template_vrf.vrf1.name
}

resource "mso_schema_template_bd" "bd2" {
  schema_id         = mso_schema.schema2.id
  template_name     = mso_schema.schema2.template_name
  name              = "bd2"
  display_name      = "bd2"
  vrf_name          = mso_schema_template_vrf.vrf1.name
  vrf_schema_id     = mso_schema_template_vrf.vrf2.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf2.template
}

resource "mso_schema_template_anp" "anp1" {
  schema_id    = mso_schema.schema1.id
  template     = mso_schema.schema1.template_name
  name         = "anp1"
  display_name = "anp1"
}

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id     = mso_schema.schema1.id
  template_name = mso_schema.schema1.template_name
  anp_name      = mso_schema_template_anp.anp1.name
  name          = "epg1"
  display_name  = "epg1"
  bd_name       = mso_schema_template_bd.bd1.name
  vrf_name      = mso_schema_template_vrf.vrf1.name
}

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id                  = mso_schema.schema1.id
  template_name              = mso_schema.schema1.template_name
  anp_name                   = mso_schema_template_anp.anp1.name
  name                       = "epg1"
  display_name               = "epg1"
  bd_name                    = mso_schema_template_bd.bd2.name
  bd_template_name           = mso_schema_template_bd.bd2.schema_id
  bd_schema_id               = mso_schema_template_bd.bd2.template_name
  vrf_name                   = mso_schema_template_vrf.vrf2.name
  vrf_schema_id              = mso_schema_template_vrf.vrf2.schema_id
  vrf_template_name          = mso_schema_template_vrf.vrf2.template
  useg_epg                   = false
  intra_epg                  = "unenforced"
  intersite_multicast_source = false
  proxy_arp                  = false
  preferred_group            = false
  access_type                = "private"
  deployment_type            = "cloud_native"
  service_type               = "custom"
  custom_service_type        = "My_Custom_Type"
}
