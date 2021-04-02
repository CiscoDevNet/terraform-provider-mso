terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
}

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id = "5ea809672c00003bc40a2799"    #id for existing schema
  template_name = "Template1"               #template name of existing template into above schema
  anp_name = "ap1"                          #anp name of existing anp in above template
  name = "mso_epg1"
  bd_name = "BD1"
  vrf_name = "DEVNET-VRF"
  display_name = "mso_epg1"
  useg_epg = true
  intra_epg = "unenforced"
  intersite_multicast_source = false
  preferred_group = false
}

resource "mso_schema_template_anp_epg_selector" "check" {
  schema_id = "${mso_schema_template_anp_epg.anp_epg.schema_id}"
  template_name = "${mso_schema_template_anp_epg.anp_epg.template_name}"
  anp_name = "${mso_schema_template_anp_epg.anp_epg.anp_name}"
  epg_name = "${mso_schema_template_anp_epg.anp_epg.name}"
  name = "check01"
  expressions {
    key = "one"
    operator = "equals"
    value = "1"
  }
  expressions {
    key = "two"
    operator = "equals"
    value = "2"
  }
}