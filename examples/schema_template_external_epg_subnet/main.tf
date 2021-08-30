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

resource "mso_schema_template_external_epg_subnet" "subnet1" {
  schema_id             = "5ea809672c00003bc40a2799"
  template_name         = "Template1"
  external_epg_name      =  "UntitledExternalEPG1"
  ip                    = "10.102.100.0/0"
  name                  = "sddfgbany"
  scope                 = ["shared-rtctrl", "export-rtctrl"]
  aggregate             = ["shared-rtctrl", "export-rtctrl"]
}