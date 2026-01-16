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

# L3Out Interface Routing Policy example with all configuration blocks

resource "mso_tenant_policies_l3out_interface_routing_policy" "routing_policy_full" {
  template_id = mso_template.tenant_template.id
  name        = "production_routing_policy"
  description = "Production L3Out Interface Routing Policy with BFD and OSPF"
  
  bfd_multi_hop_settings {
    admin_state           = "enabled"
    detection_multiplier  = 3
    min_receive_interval  = 250
    min_transmit_interval = 250
  }
  
  bfd_settings {
    admin_state           = "enabled"
    detection_multiplier  = 3
    min_receive_interval  = 50
    min_transmit_interval = 50
    echo_receive_interval = 50
    echo_admin_state      = "enabled"
    interface_control     = false
  }
  
  ospf_interface_settings {
    network_type          = "point_to_point"
    priority              = 100
    cost_of_interface     = 10
    hello_interval        = 10
    dead_interval         = 40
    retransmit_interval   = 5
    transmit_delay        = 1
    advertise_subnet      = true
    bfd                   = true
    mtu_ignore            = false
    passive_participation = false
  }
}
