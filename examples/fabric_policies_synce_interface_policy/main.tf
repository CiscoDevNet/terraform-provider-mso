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

# fabric policies synce interface policy example

resource "mso_fabric_policies_synce_interface_policy" "synce_interface_policy" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "synce_interface_policy"
  description     = "Example description"
  admin_state     = "enabled"
  sync_state_msg  = "enabled"
  selection_input = "enabled"
  src_priority    = 120
  wait_to_restore = 6
}
