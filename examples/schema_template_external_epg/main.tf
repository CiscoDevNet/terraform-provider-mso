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

resource "mso_schema_template_external_epg" "template_externalepg" {
  schema_id                  = "5eba96b22c0000ed0981291e"
  template_name              = "Template1"
  external_epg_name          = "external_epg12"
  display_name               = "external_epg12"
  vrf_name                   = "demo_vrf"
  vrf_schema_id              = "5eba96b22c0000ed0981291e"
  vrf_template_name          = "Template1"
  external_epg_type          = "on-premise"
  l3out_name                 = "nk_l3out"
  include_in_preferred_group = false
  anp_name                   = "demo"
}
