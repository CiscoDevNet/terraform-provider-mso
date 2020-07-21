provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
#   insecure = true
#   proxy_url = "https://proxy_server:proxy_port"
}

resource "mso_schema_template_contract_service_graph" "one" {
  schema_id = "5f11b0e22c00001c4a812a2a"
  site_id = "5c7c95b25100008f01c1ee3c"
  template_name = "Template1"
  contract_name = "UntitledContract1"
  service_graph_name = "sg1"
  node_relationship {
    provider_connector_bd_name = "BD1"
    consumer_connector_bd_name = "BD2"
    provider_connector_cluster_interface = "test"
    consumer_connector_cluster_interface = "test"
    provider_connector_redirect_policy_tenant = "NkAutomation"
    provider_connector_redirect_policy = "test2"
    consumer_connector_redirect_policy_tenant = "NkAutomation"
    consumer_connector_redirect_policy = "test2"
    provider_subnet_ips = ["1.2.3.4/20"]
    consumer_subnet_ips = ["1.2.3.4/20"]
  }
}