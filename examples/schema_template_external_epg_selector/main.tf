provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
}

resource "mso_schema_template_external_epg" "template_externalepg" {
	schema_id = "5ea809672c00003bc40a2799"
	template_name = "Template1"
	external_epg_name = "check_anp01"
	display_name = "check_anp01"
	vrf_name = "demo"
}

resource "mso_schema_template_external_epg_selector" "selector1" {
	schema_id = "${mso_schema_template_external_epg.template_externalepg.schema_id}"
	template = "${mso_schema_template_external_epg.template_externalepg.template_name}"
	external_epg_name = "${mso_schema_template_external_epg.template_externalepg.external_epg_name}"
	name = "check01"
    expressions {
      value = "1.20.56.44"
    }
    expressions{
      value = "5.6.7.8"
    }
}