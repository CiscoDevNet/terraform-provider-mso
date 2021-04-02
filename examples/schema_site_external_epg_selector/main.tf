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
#   insecure = true
#   proxy_url = "https://proxy_server:proxy_port"
}

resource "mso_schema_template_external_epg" "template_externalepg" {
	schema_id           = "5f043b3b2c0000f47e812a0b"
	template_name       = "Template1"
	external_epg_name   = "temp_epg"
    external_epg_type   = "cloud"
	display_name        = "temp_epg"
	vrf_name            = "Myvrf"
    anp_name            = "ap1"
    l3out_name          = "temp"
    site_id             = ["5c7c95d9510000cf01c1ee3d"]
    selector_name       = "check02"
    selector_ip         = "12.23.34.45"
}

resource "mso_schema_site_external_epg_selector" "sel1" {
  schema_id = "${mso_schema_template_external_epg.template_externalepg.schema_id}"
  template_name = "${mso_schema_template_external_epg.template_externalepg.template_name}"
  site_id = "${mso_schema_template_external_epg.template_externalepg.site_id}"
  external_epg_name = "${mso_schema_template_external_epg.template_externalepg.external_epg_name}"
  name = "second"
  ip = "12.25.70.50"
}