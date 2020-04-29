provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_template_bd" "bridge_domain" {
    schema_id = "5ea809672c00003bc40a2799"
    template_name = "Template1"
    name = "testBD"
    display_name = "test"
    vrf_name = "demo"
    layer2_unknown_unicast = "proxy"
  
}
