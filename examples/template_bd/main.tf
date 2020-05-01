provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template_bd" "bridge_domain" {
    schema_id = "5ea809672c00003bc40a2799"
    template_name = "Template1"
    name = "testBD3"
    display_name = "testwe"
    vrf_name = "demo"
    layer2_unknown_unicast = "proxy"
    dhcp_policy = {
      name = "etes"
      version = 1
      "dhcp_option_policy_name"    = "demo" 
      "dhcp_option_policy_version" = 2
    }

}


data "mso_schema_template_bd" "bd_data" {
    schema_id = "5ea809672c00003bc40a2799"
    template_name = "Template1"
    name = "testBD"
}

output "demo_bd" {
  value = data.mso_schema_template_bd.bd_data
}

