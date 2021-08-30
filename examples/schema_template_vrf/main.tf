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

resource "mso_schema_template_vrf" "demovrf01" {
    schema_id       = "5c4d5bb72700000401f80948"
    template        ="Temp200"
	name            = "vrf982"
	display_name    ="vrf982"
	layer3_multicast=false
  
}
