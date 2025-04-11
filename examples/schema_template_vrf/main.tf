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
  name = "example_tenant"
}

resource "mso_schema" "schema_test" {
  name = "terraform_schema"
  template {
    name         = "terraform_schema_template"
    display_name = "terraform_schema_template"
    tenant_id    = data.mso_tenant.tenant_test.id
  }
}

# Example 1: Creating a VRF with minimum configuration

resource "mso_schema_template_vrf" "vrf" {
  schema_id    = mso_schema.schema_test.id
  template     = one(mso_schema.schema_test.template).name
  name         = "template_vrf"
  display_name = "template_vrf"
}

#  Example 2: Creating a VRF with layer3 Multicast enabled and with a defined Rendezvous Point.

data "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
  template_id = mso_template.tenant_template.id
  name        = "route_map_policy_multicast"
}

resource "mso_schema_template_vrf" "vrf_layer3_multicast" {
  schema_id        = mso_schema_template_vrf.vrf.schema_id
  template         = mso_schema_template_vrf.vrf.template
  name             = "layer3_multicast_vrf"
  display_name     = "layer3_multicast_vrf"
  layer3_multicast = true
  rendezvous_points {
    ip_address                      = "1.1.1.2"
    type                            = "static"
    route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
  }
  rendezvous_points {
    ip_address = "1.1.1.3"
    type       = "fabric"
  }
}

#  Example 3: Creating a VRF with preferred group enabled and associating it to an Application EPG.

resource "mso_schema_template_vrf" "vrf_preferred_group" {
  schema_id              = mso_schema_template_vrf.vrf.schema_id
  template               = mso_schema_template_vrf.vrf.template
  name                   = "preferred_goup_vrf"
  display_name           = "preferred_goup_vrf"
  preferred_group        = true
  vzany                  = false
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
  schema_id       = mso_schema.schema_test.id
  template_name   = mso_schema_template_anp.anp.template
  anp_name        = mso_schema_template_anp.anp.name
  name            = "epg"
  display_name    = "epg"
  bd_name         = mso_schema_template_bd.bd.name
  vrf_name        = mso_schema_template_vrf.vrf_preferred_group.name
  preferred_group = true
}

data "mso_schema_template_vrf" "vrf_preferred_group_data" {
  schema_id = mso_schema_template_vrf.vrf_preferred_group.schema_id
  template  = mso_schema_template_vrf.vrf_preferred_group.template
  name      = mso_schema_template_vrf.vrf_preferred_group.name
}

output "schema_template_vrf_data" {
  value = data.mso_schema_template_vrf.vrf_preferred_group_data
}
