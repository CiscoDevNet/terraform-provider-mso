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

resource "mso_service_node_type" "node_type" {
  name         = "tftst"
  display_name = "terrform type"
}

data "mso_service_node_type" "node_data" {
  name = "tftst"
}
