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

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id                  = "5eafca7d2c000052860a2902"
  template_name              = "stemplate1"
  anp_name                   = "sanp1"
  name                       = "nkusegte"
  bd_name                    = "testBD"
  vrf_name                   = "vrf1"
  display_name               = "nkuseg"
  useg_epg                   = true
  intra_epg                  = "enforced"
  intersite_multicast_source = true
  proxy_arp                  = true
  preferred_group            = true
  bd_template_name           = "stemplate1"
  vrf_schema_id              = "5eafeb792c0000a18e0a2900"
  bd_schema_id               = "5eafeb792c0000a18e0a2900"
  vrf_template_name          = "stemplate1"

}
