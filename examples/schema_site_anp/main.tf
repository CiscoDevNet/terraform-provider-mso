terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO pwd>
  url      = "" # <MSO URL>
  insecure = true
}

resource "mso_schema_site_anp" "anp1" {
  schema_id     = "5c6c16d7270000c710f8094d"
  anp_name      = "AP1234"
  template_name = "Template1"
  site_id       = "5c7c95d9510000cf01c1ee3d"
}
