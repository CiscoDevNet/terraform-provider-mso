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
  insecure = true
}

resource "mso_schema_template_anp_epg_subnet" "subnet1" {
  schema_id     = "5c4d5bb72700000401f80948"
  anp_name      = "ANP"
  epg_name      = "Web"
  template      = "Template1"
  ip            = "31.101.102.0/8"
  scope         = "public"
  shared        = true
}