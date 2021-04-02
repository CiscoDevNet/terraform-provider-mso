terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template_bd" "bridgedomain" {
    schema_id              = "5c4d5bb72700000401f80948"
    template_name          = "Temp200"
    name                   = "testBD"
    display_name           = "test"
    vrf_name               = "vrf982"
    vrf_schema_id          = "5c4d5bb72700000401f80948"
    vrf_template_name      = "Temp200"
    layer2_unknown_unicast = "proxy" 
    intersite_bum_traffic  = false
    optimize_wan_bandwidth = true
    layer2_stretch         = true
    layer3_multicast       = false
    dhcp_policy = {
        name = "Policy1"
        version = 10
        dhcp_option_policy_name = "Policy10"
        dhcp_option_policy_version = 12
    }

  
}
