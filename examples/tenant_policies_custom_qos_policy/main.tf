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

data "mso_tenant" "example_tenant" {
  name = "example_tenant"
}

# tenant template example

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id = data.mso_tenant.example_tenant.id
}

# tenant policies custom qos policy example

resource "mso_tenant_policies_custom_qos_policy" "qos_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_custom_qos_policy"
  description = "Custom QoS Policy"

  dscp_mappings {
      dscp_from    = "af11"
      dscp_to      = "af12"
      dscp_target  = "af11"
      target_cos   = "background"
      qos_priority = "level1"
  }

  cos_mappings {
      dot1p_from   = "background"
      dot1p_to     = "best_effort"
      dscp_target  = "af11"
      target_cos   = "background"
      qos_priority = "level1"
  }
}
