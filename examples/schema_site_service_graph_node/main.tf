provider "mso" {
  username = "terraform_github_ci"
  password = "Crest@123456"
  url      = "https://173.36.219.66/"
  insecure = true
}

resource "mso_schema_site_service_graph_node" "test_sg" {
  schema_id          = "5f06a4c40f0000b63dbbd647"
  template_name      = "Template1"
  service_graph_name = "sgtf"
  service_node_type  = "firewall"
  site_nodes {
    site_id     = "5f05c69f1900002234d0537e"
    tenant_name = "NkAutomation"
    node_name   = "nk-fw-2"
  }

}

