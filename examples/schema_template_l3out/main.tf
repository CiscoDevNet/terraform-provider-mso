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

resource "mso_schema_template_l3out" "template_l3out" {
  schema_id             = "5c6c16d7270000c710f8094d"
  template_name         = "Template1"
  l3out_name            = "l3out100"
  display_name          = "l3out100"
  vrf_name              = "vrf2"
  vrf_schema_id         = "5c6c16d7270000c710f8094d"
  vrf_template_name     = "Template1"
}