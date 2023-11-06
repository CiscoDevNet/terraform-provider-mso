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


data "mso_schema" "schema1" {
  name = "schema1"
}

data "mso_schema" "schema2" {
  name = "schema2"
}

data "mso_site" "site1" {
  name = "site1"
}

data "mso_schema_template" "t1" {
  name      = "t1"
  schema_id = data.mso_schema.schema1.id
}

data "mso_schema_template" "t2" {
  name      = "t2"
  schema_id = data.mso_schema.schema2.id
}

data "mso_schema_template_bd" "bd1" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "t1"
  name          = "bd1"
}

data "mso_schema_template_bd" "bd2" {
  schema_id     = data.mso_schema.schema2.id
  template_name = "t2"
  name          = "bd2"
}

data "mso_schema_template_contract" "c1" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "t1"
  contract_name = "c1"
}

data "mso_schema_template_service_graph" "sg1" {
  schema_id          = data.mso_schema.schema1.id
  template_name      = "t1"
  service_graph_name = "sg1"
}

resource "mso_schema_site" "schema_site_1" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = data.mso_schema_template_service_graph.sg1.template_name
}

resource "mso_schema_site_contract_service_graph" "example" {
  schema_id          = data.mso_schema.schema1.id
  template_name      = "t1"
  contract_name      = "c1"
  service_graph_name = "sg1"
  site_id            = data.mso_site.site1.id
  node_relationship {
    provider_connector_cluster_interface      = "example_provider_cluster_interface"
    provider_connector_redirect_policy_tenant = "example_tenant"
    provider_connector_redirect_policy        = "example_redirect_policy"
    consumer_connector_cluster_interface      = "example_consumer_cluster_interface"
    consumer_connector_redirect_policy_tenant = "example_tenant"
    consumer_connector_redirect_policy        = "example_redirect_policy"
    consumer_subnet_ips                       = ["1.1.1.1/24", "2.2.2.2/24"]
  }
  node_relationship {
    provider_connector_cluster_interface = "example_provider_cluster_interface"
    consumer_connector_cluster_interface = "example_consumer_cluster_interface"
  }
}

# Cloud Network Controller site configuration
resource "mso_schema_site_contract_service_graph" "example" {
  schema_id          = data.mso_schema.schema1.id
  template_name      = "t1"
  contract_name      = "c1"
  service_graph_name = "sg1"
  site_id            = data.mso_site.site1.id
}
