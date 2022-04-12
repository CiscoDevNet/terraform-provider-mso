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

resource "mso_schema_site_anp_epg_domain" "foo_schema_site_anp_epg_domain" {
  schema_id                 = "5c4d9fca270000a101f8094a"
  template_name             = "Template1"
  site_id                   = "5c7c95b25100008f01c1ee3c"
  anp_name                  = "ANP"
  epg_name                  = "Web"
  domain_type               = "vmmDomain"
  vmm_domain_type           = "Microsoft"
  domain_name               = "VM-Micro"
  deploy_immediacy          = "immediate"
  resolution_immediacy      = "immediate"
  vlan_encap_mode           = "static"
  allow_micro_segmentation  = true
  switching_mode            = "native"
  switch_type               = "default"
  micro_seg_vlan_type       = "vlan"
  micro_seg_vlan            = 46
  port_encap_vlan_type      = "vlan"
  port_encap_vlan           = 45
  enhanced_lag_policy_name  = "name"
  enhanced_lag_policy_dn    = "dn"
}

resource "mso_schema_site_anp_epg_domain" "foo_schema_site_anp_epg_domain" {
  schema_id                 = "5c4d9fca270000a101f8094a"
  template_name             = "Template1"
  site_id                   = "5c7c95b25100008f01c1ee3c"
  anp_name                  = "Anp"
  epg_name                  = "Web"
  domain_dn                 = "uni/vmmp-VMware/dom-TEST"
  deploy_immediacy          = "immediate"
  resolution_immediacy      = "immediate"
  vlan_encap_mode           = "static"
  allow_micro_segmentation  = false
  switching_mode            = "native"
  switch_type               = "default"
  micro_seg_vlan_type       = "vlan"
  micro_seg_vlan            = 46
  port_encap_vlan_type      = "vlan"
  port_encap_vlan           = 45
  enhanced_lag_policy_name  = "name"
  enhanced_lag_policy_dn    = "uni/vmmp-VMware/dom-TEST/vswitchpolcont/enlacplagp-Lacp"
}
