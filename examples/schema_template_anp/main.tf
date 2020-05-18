provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template_anp" "anp1" {
  schema_id     = "5c4d5bb72700000401f80948"
  template      = "Template1"
  name          = "Demo_ANP"
  display_name  = "anp1234"
}
