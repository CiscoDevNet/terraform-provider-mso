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

resource "mso_tenant" "tf_test_tenant" {
  name         = "tf_test_tenant"
  display_name = "tf_test_tenant"
}

resource "mso_schema" "schema1" {
  name = "tf_test_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = mso_tenant.tf_test_tenant.id
  }
}

resource "mso_schema" "schema2" {
  name = "tf_test_schema_2"
  template {
    name         = "Template2"
    display_name = "Template2"
    tenant_id    = mso_tenant.tf_test_tenant.id
  }
}

resource "mso_schema_template_vrf" "vrf1" {
  schema_id    = mso_schema.schema1.id
  template     = one(mso_schema.schema1.template).name
  name         = "vrf1"
  display_name = "vrf1"
}

resource "mso_schema_template_vrf" "vrf2" {
  schema_id    = mso_schema.schema2.id
  template     = one(mso_schema.schema2.template).name
  name         = "vrf2"
  display_name = "vrf2"
}

resource "mso_schema_template_bd" "bd1" {
  schema_id                       = mso_schema.schema1.id
  template_name                   = one(mso_schema.schema1.template).name
  name                            = "bd1"
  display_name                    = "bd1"
  vrf_name                        = mso_schema_template_vrf.vrf1.name
  layer2_unknown_unicast          = "proxy"
  intersite_bum_traffic           = false
  optimize_wan_bandwidth          = true
  layer2_stretch                  = false
  layer3_multicast                = false
  multi_destination_flooding      = "flood_in_encap"
  ipv6_unknown_multicast_flooding = "optimized_flooding"
  unknown_multicast_flooding      = "optimized_flooding"
  description                     = "bd1_description"
}

resource "mso_schema_template_bd" "bd2" {
  schema_id                       = mso_schema.schema2.id
  template_name                   = one(mso_schema.schema2.template).name
  name                            = "bd2"
  display_name                    = "bd2"
  vrf_name                        = mso_schema_template_vrf.vrf2.name
  vrf_schema_id                   = mso_schema_template_vrf.vrf2.schema_id
  vrf_template_name               = mso_schema_template_vrf.vrf2.template
  layer2_unknown_unicast          = "proxy"
  intersite_bum_traffic           = false
  optimize_wan_bandwidth          = true
  layer2_stretch                  = false
  layer3_multicast                = false
  multi_destination_flooding      = "flood_in_encap"
  ipv6_unknown_multicast_flooding = "optimized_flooding"
  unknown_multicast_flooding      = "optimized_flooding"
  description                     = "bd2_description"
}

resource "mso_schema_template_anp" "anp1" {
  schema_id    = mso_schema.schema1.id
  template     = one(mso_schema.schema1.template).name
  name         = "anp1"
  display_name = "anp1"
}

resource "mso_schema_template_anp_epg" "anp_epg1" {
  schema_id     = mso_schema.schema1.id
  template_name = one(mso_schema.schema1.template).name
  anp_name      = mso_schema_template_anp.anp1.name
  name          = "epg1"
  display_name  = "epg1"
  bd_name       = mso_schema_template_bd.bd1.name
  vrf_name      = mso_schema_template_vrf.vrf1.name
  description   = "anp_epg1_description"
}

resource "mso_schema_template_anp_epg" "anp_epg2" {
  schema_id                  = mso_schema.schema1.id
  template_name              = one(mso_schema.schema1.template).name
  anp_name                   = mso_schema_template_anp.anp1.name
  name                       = "epg2"
  display_name               = "epg2"
  bd_name                    = mso_schema_template_bd.bd2.name
  bd_template_name           = mso_schema_template_bd.bd2.template_name
  bd_schema_id               = mso_schema_template_bd.bd2.schema_id
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
  description                = "anp_epg2_description"
}
