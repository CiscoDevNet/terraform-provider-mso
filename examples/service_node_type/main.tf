provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_service_node_type" "node_type" {
  name         = "tftst"
  display_name = "terrform type"
}

data "mso_service_node_type" "node_data" {
  name = "tftst"
}
