provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_site" "schemasite01" {
  schema_id      = "5c4d5bb72700000401f80948"
  site_id        = "5ec240152e00007d763c18b4"
  template_name  = "Template1"
  
}
