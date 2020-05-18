provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template_externalepg" "template_externalepg" {
  schema_id             = "5ea809672c00003bc40a2799"
  template_name         = "Template1"
  externalepg_name      = "external_epg12"
  display_name          = "external_epg12"
  vrf_name              = "vrf1"
  vrf_schema_id         = "5c6c16d7270000c710f8094d"
  vrf_template_name     = "Template1"
}