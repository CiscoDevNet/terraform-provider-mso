provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_site_bd" "foo_schema_site_bd" {
  schema_id     = "5d5dbf3f2e0000580553ccce"
  bd_name       = "bd4"
  template_name = "Template1"
  site_id       = "5c7c95b25100008f01c1ee3c"
  host_route          = false
}