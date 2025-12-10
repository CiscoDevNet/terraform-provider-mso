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

# VLAN pool example

resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
  template_id = mso_template.fabric_policy_template.id
  name        = "production_vlan_pool"
  description = "Production VLAN Pool for L3 Domains"
  
  vlan_range {
    from        = 100
    to          = 199
    description = "VLAN range for site A"
  }
}

# L3 domain example with VLAN pool

resource "mso_fabric_policies_l3_domain" "l3_domain" {
  template_id    = mso_template.fabric_policy_template.id
  name           = "production_l3_domain"
  description    = "Production L3 Domain for external connectivity"
  vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
}
