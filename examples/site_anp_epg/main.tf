provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

data "mso_schema_site_anp_epg" "anpEpg" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "abcd"
  site_id = "5c7c95b25100008f01c1ee3c"
  anp_name = "ANP"
  epg_name = "DB"
}