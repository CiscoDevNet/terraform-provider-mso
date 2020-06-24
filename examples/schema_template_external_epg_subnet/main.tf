provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
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