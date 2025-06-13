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

# fabric policies vlan pool example

resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "vlan_pool"
  description     = "Example description"
  vlan_range {
    from            = 200
    to              = 202
  }
  vlan_range {
    from            = 204
    to              = 209
  }
}

# fabric policies physical domain example

resource "mso_fabric_policies_physical_domain" "physical_domain" {
  template_id     = mso_template.fabric_policy_template.id
  name            = "physical_domain"
  description     = "Example description"
  vlan_pool_uuid  = mso_fabric_policies_vlan_pool.vlan_pool.uuid
}