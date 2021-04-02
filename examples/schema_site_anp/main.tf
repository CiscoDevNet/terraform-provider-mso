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

resource "mso_schema_site_anp" "anp1" {
  schema_id     = "5c6c16d7270000c710f8094d"
  anp_name      = "AP1234"
  template_name = "Template1"
  site_id       = "5c7c95d9510000cf01c1ee3d"
}
