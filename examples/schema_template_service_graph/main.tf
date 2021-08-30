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

resource "mso_schema_template_service_graph" "test_sg" {
  schema_id          = "5f06a4c40f0000b63dbbd647"
  template_name      = "Template1"
  service_graph_name = "sgtf"
  service_node_type  = "firewall"
  description        = "hello"
  site_nodes {
    site_id     = "5f05c69f1900002234d0537e"
    tenant_name = "NkAutomation"
    node_name   = "nk-fw-2"
  }

}


data "mso_schema_template_service_graph" "test_sg" {
  schema_id          = "5f0830910e0000a348cff85e"
  template_name      = "Template1"
  service_graph_name = "sgt"
  node_index         = 2

}
