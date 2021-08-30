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

resource "mso_schema_template_anp" "anp1" {
  schema_id     = "5c4d5bb72700000401f80948"
  template      = "Template1"
  name          = "Demo_ANP"
  display_name  = "anp1234"
}
