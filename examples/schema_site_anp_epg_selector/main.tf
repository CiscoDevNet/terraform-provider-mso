provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
}

resource "mso_schema_site_anp_epg" "site_anp_epg" {
  schema_id = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id = "5c7c95d9510000cf01c1ee3d"
  anp_name = "ANP"
  epg_name = "DB"
}

resource "mso_schema_site_anp_epg_selector" "check" {
  schema_id     = "${mso_schema_site_anp_epg.site_anp_epg.schema_id}"
  site_id       = "${mso_schema_site_anp_epg.site_anp_epg.site_id}"
  template      = "${mso_schema_site_anp_epg.site_anp_epg.template_name}"
  anp_name      = "${mso_schema_site_anp_epg.site_anp_epg.anp_name}"
  epg_name      = "${mso_schema_site_anp_epg.site_anp_epg.epg_name}"
  name          = "check01"
  expressions {
    key         = "one"
    operator    = "equals"
    value       = "1"
  }
  expressions {
    key         = "two"
    operator    = "notEquals"
    value       = "22"
  }
}