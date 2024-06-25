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

data "mso_tenant" "tenant_test" {
  name         = "ansible_test"
  display_name = "ansible_test"
}

resource "mso_schema" "schema_test" {
  name = "terraform_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = data.mso_tenant.tenant_test.id
  }
}

#  Associating VRF (template_vrf) to Template1
resource "mso_schema_template_vrf" "vrf" {
  schema_id        = mso_schema.schema_test.id
  template         = one(mso_schema.schema_test.template).name
  name             = "template_vrf"
  display_name     = "template_vrf"
  layer3_multicast = true
  vzany            = false
}

resource "mso_schema_template_vrf" "vrf_preferred_group" {
  schema_id              = mso_schema_template_vrf.vrf.schema_id
  template               = mso_schema_template_vrf.vrf.template
  name                   = "preferred_goup_vrf"
  display_name           = "preferred_goup_vrf"
  preferred_group        = true
  ip_data_plane_learning = "disabled"
}

resource "mso_schema_template_bd" "bd" {
  schema_id         = mso_schema_template_vrf.vrf_preferred_group.schema_id
  template_name     = mso_schema_template_vrf.vrf_preferred_group.template
  name              = "bd"
  display_name      = "bd"
  vrf_name          = mso_schema_template_vrf.vrf_preferred_group.name
  vrf_schema_id     = mso_schema_template_vrf.vrf_preferred_group.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf_preferred_group.template
  arp_flooding      = true
}

resource "mso_schema_template_anp" "anp" {
  schema_id    = mso_schema_template_bd.bd.schema_id
  template     = mso_schema_template_bd.bd.template_name
  name         = "anp"
  display_name = "anp"
}

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id                     = mso_schema.schema_test.id
  template_name                 = mso_schema_template_anp.anp.template
  anp_name                      = mso_schema_template_anp.anp.name
  name                          = "epg"
  display_name                  = "epg"
  bd_name                       = mso_schema_template_bd.bd.name
  vrf_name                      = mso_schema_template_vrf.vrf_preferred_group.name
  preferred_group               = true
  site_aware_policy_enforcement = true
}

data "mso_schema_template_vrf" "vrf_preferred_group_data" {
  schema_id = mso_schema_template_vrf.vrf_preferred_group.schema_id
  template  = mso_schema_template_vrf.vrf_preferred_group.template
  name      = mso_schema_template_vrf.vrf_preferred_group.name
}

output "schema_template_vrf_data" {
  value = data.mso_schema_template_vrf.vrf_preferred_group_data
}
