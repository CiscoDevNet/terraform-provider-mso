provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template" "schematemplate01" {
  schema_id     = "5c4d5bb72700000401f80948"
  name          = "Temp200"
  display_name  = "Temp845"
  tenant_id     = "5c4d9f3d2700007e01f80949" 
}
