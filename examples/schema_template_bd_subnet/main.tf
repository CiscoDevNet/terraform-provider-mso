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

resource "mso_schema_template_bd_subnet" "bdsubnet01" {
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  bd_name = "testBD"
  ip = "10.26.17.0/8"
  scope = "public"
  description = "SubnetDemo"
  shared = false
  no_default_gateway = true
  querier = true
  
  
}
