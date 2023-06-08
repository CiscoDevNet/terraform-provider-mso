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

resource "mso_tenant" "tf_tenant" {
  name         = "tf_tenant"
  display_name = "tf_tenant"
}

resource "mso_schema" "tf_schema" {
  name = "tf_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = mso_tenant.tf_tenant.id
  }
}

resource "mso_schema_template_vrf" "vrf" {
  schema_id              = mso_schema.tf_schema.id
  template               = "Template1"
  name                   = "template_vrf"
  display_name           = "template_vrf"
  layer3_multicast       = false
  vzany                  = false
  ip_data_plane_learning = "disabled"
}

resource "mso_schema_template_external_epg" "template_externalepg" {
  schema_id                  = mso_schema.tf_schema.id
  template_name              = "Template1"
  external_epg_name          = "external_epg12"
  display_name               = "external_epg12"
  vrf_name                   = mso_schema_template_vrf.vrf.name
  vrf_schema_id              = mso_schema.tf_schema.id
  vrf_template_name          = "Template1"
  external_epg_type          = "on-premise"
  include_in_preferred_group = false
}
