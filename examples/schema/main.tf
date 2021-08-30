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
  // platform = "nd"  // use when logging in ND.
}

// create data source
resource "mso_schema" "schema1" {
  name          = "demo_schema"
  template_name = "tempu"
  tenant_id     = "0000ffff0000000000000010"
}