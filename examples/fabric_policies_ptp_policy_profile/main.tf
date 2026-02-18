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

# fabric policies ptp policy example

resource "mso_fabric_policies_ptp_policy" "ptp_policy" {
  template_id               = mso_template.fabric_policy_template.id
  name                      = "ptp_policy"
  description               = "Example description"
  admin_state               = "enabled"
  global_priority1          = 250
  global_priority2          = 100
  global domain             = 99
  fabric_profile_template   = "aes67"
  fabric_announce_interval  = 1
  fabric_sync_interval      = -3
  fabric_delay_interval     = -2
  fabric_announce_timeout   = 3
}

# fabric policies ptp policy profile example

resource "mso_fabric_policies_ptp_policy_profile" "ptp_policy_profile" {
  template_id         = mso_template.fabric_policy_template.id
  name                  = "ptp_policy_profile"
  description           = "Example description"
  delay_interval        = -2
  sync_interval         = -3
  announce_timeout      = 3
  announce_interval     = 1
  profile_template      = "aes67"
  override_node_profile = false
}
