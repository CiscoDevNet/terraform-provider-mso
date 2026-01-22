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

# Tenant template example

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id     = data.mso_tenant.example_tenant.id
}

# L3Out Node Routing Policy example with all configuration blocks

resource "mso_tenant_policies_l3out_node_routing_policy" "node_policy_full" {
  template_id             = mso_template.tenant_template.id
  name                    = "production_node_routing_policy"
  description             = "Production L3Out Node Routing Policy with BFD and BGP"
  as_path_multipath_relax = true
  
  bfd_multi_hop_settings {
    admin_state           = "enabled"
    detection_multiplier  = 3
    min_receive_interval  = 250
    min_transmit_interval = 250
  }
  
  bgp_node_settings {
    graceful_restart_helper = true
    keep_alive_interval     = 60
    hold_interval           = 180
    stale_interval          = 300
    max_as_limit            = 0
  }
}
