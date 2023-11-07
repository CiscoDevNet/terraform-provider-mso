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

data "mso_tenant" "example_tenant" {
  name = "example_tenant"
}

data "mso_site" "example_site" {
  name = "example_site"
}

resource "mso_schema" "schema_1" {
  name = "schema_1"
  template {
    name         = "template_1"
    display_name = "template_1"
    tenant_id    = data.mso_tenant.example_tenant.id
  }
}

resource "mso_schema_template_vrf" "vrf_1" {
  schema_id    = mso_schema.schema_1.id
  template     = one(mso_schema.schema_1.template).name
  name         = "vrf_1"
  display_name = "vrf_1"
}

resource "mso_schema_template_bd" "bd_1" {
  schema_id                       = mso_schema.schema_1.id
  template_name                   = one(mso_schema.schema_1.template).name
  name                            = "bd_1"
  display_name                    = "bd_1"
  vrf_name                        = mso_schema_template_vrf.vrf_1.name
  layer2_unknown_unicast          = "proxy"
  intersite_bum_traffic           = false
  optimize_wan_bandwidth          = true
  layer2_stretch                  = false
  layer3_multicast                = false
  multi_destination_flooding      = "flood_in_encap"
  ipv6_unknown_multicast_flooding = "optimized_flooding"
  unknown_multicast_flooding      = "optimized_flooding"
  description                     = "bd_1_description"
}

resource "mso_schema_template_anp" "anp_1" {
  schema_id    = mso_schema.schema_1.id
  template     = one(mso_schema.schema_1.template).name
  name         = "anp_1"
  display_name = "anp_1"
}

resource "mso_schema_template_anp_epg" "anp_epg_1" {
  schema_id     = mso_schema.schema_1.id
  template_name = one(mso_schema.schema_1.template).name
  anp_name      = mso_schema_template_anp.anp_1.name
  name          = "anp_epg_1"
  display_name  = "anp_epg_1"
  bd_name       = mso_schema_template_bd.bd_1.name
  vrf_name      = mso_schema_template_vrf.vrf_1.name
  description   = "anp_epg1_description"
}

resource "mso_schema_site" "schema_site_1" {
  schema_id     = mso_schema.schema_1.id
  site_id       = data.mso_site.example_site.id
  template_name = one(mso_schema.schema_1.template).name
}

# Examples mso_schema_site_anp_epg_domain vmmDomain in version >= 4.2

resource "mso_schema_site_anp_epg_domain" "vmware_domain_id_4_2_up" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_type              = "vmmDomain"
  vmm_domain_type          = "Microsoft"
  domain_name              = "VM-Micro"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = true
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
  delimiter                = "|"
  binding_type             = "static"
  port_allocation          = "fixed"
  num_ports                = 3
  netflow                  = "disabled"
  allow_promiscuous        = "accept"
  mac_changes              = "reject"
  forged_transmits         = "reject"
  custom_epg_name          = "custom_epg_name_1"
}

resource "mso_schema_site_anp_epg_domain" "vmware_domain_name_4_2_up" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_dn                = "uni/vmmp-VMware/dom-TEST"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = true
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
  delimiter                = "|"
  binding_type             = "static"
  port_allocation          = "fixed"
  num_ports                = 3
  netflow                  = "disabled"
  allow_promiscuous        = "accept"
  mac_changes              = "reject"
  forged_transmits         = "reject"
  custom_epg_name          = "custom_epg_name_1"
}

# Examples mso_schema_site_anp_epg_domain vmmDomain in version < 4.2

resource "mso_schema_site_anp_epg_domain" "vmware_domain_with_name_pre_4_2" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_type              = "vmmDomain"
  vmm_domain_type          = "Microsoft"
  domain_name              = "VM-Micro"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = true
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
}

resource "mso_schema_site_anp_epg_domain" "vmware_domain_enhanced_lag_policy_with_domain_dn_pre_4_2" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_dn                = "uni/vmmp-VMware/dom-TEST"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = false
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
  enhanced_lag_policy_name = "Lacp"
  enhanced_lag_policy_dn   = "uni/vmmp-VMware/dom-TEST/vswitchpolcont/enlacplagp-Lacp"
}
