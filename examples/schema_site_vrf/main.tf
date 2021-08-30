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

resource "mso_schema_site_vrf" "foo_schema_site_vrf" {
  template_name = "Template1"
  site_id       = "5c7c95d9510000cf01c1ee3d"
  schema_id     = "5c6c16d7270000c710f8094d"
  vrf_name      = "vrf3"
}