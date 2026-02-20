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

# Fabric policy template example
resource "mso_template" "fabric_policy_template" {
  template_name = "fabric_policy_template"
  template_type = "fabric_policy"
}

resource "mso_fabric_policies_mcp_global_policy" "test" {
    template_id                       = mso_template.fabric_policy_template.id
    name                              = "test"
    admin_state                       = "disabled"
    description                       = "Test"
    enable_mcp_pdu_per_vlan           = "disabled"
    initial_delay_time                = 180
    key                               = "test_key_123"
    loop_detect_multiplication_factor = 3
    port_disable_protection           = "enabled"
    transmission_frequency_msec       = 0
    transmission_frequency_sec        = 2
}
