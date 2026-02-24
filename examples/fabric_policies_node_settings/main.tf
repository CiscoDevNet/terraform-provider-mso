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

# fabric policies node settings example

resource "mso_fabric_policies_node_settings" "node_settings" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "tf_node_settings"
  description     = "Terraform Node Settings Policy"
  synce = {
    admin_state   = "enabled"
    quality_level = "option_2_generation_1"
  }
  ptp = {
    node_domain   = 25
    priority_2    = 99
  }
}
