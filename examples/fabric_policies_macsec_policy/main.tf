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

# fabric policy template example

resource "mso_template" "fabric_policy_template" {
  template_name = "fabric_policy_template"
  template_type = "fabric_policy"
}

# fabric policies macsec policy example

resource "mso_fabric_policies_macsec_policy" "macsec_policy" {
  template_id            = mso_template.fabric_policy_template.id
  name                   = "macsec_policy"
  description            = "Example description"
  admin_state            = "enabled"
  interface_type         = "access"
  cipher_suite           = "256GcmAes"
  window_size            = 128
  security_policy        = "shouldSecure"
  sak_expire_time        = 60
  confidentiality_offset = "offset30"
  key_server_priority    = 8
  macsec_key {
    key_name             = "abc123"
    psk                  = "AA111111111111111111111111111111111111111111111111111111111111aa"
    start_time           = "now"
    end_time             = "2027-09-23 00:00:00"
	}
}
