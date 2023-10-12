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

resource "mso_schema_template_contract_service_graph" "example" {
  schema_id          = data.mso_schema.schema1.id
  template_name      = data.mso_schema_template.t1.name
  contract_name      = data.mso_schema_template_contract.c1.name
  service_graph_name = data.mso_schema_template_service_graph.sg1.service_graph_name
  node_relationship {
    consumer_connector_bd_name          = data.mso_schema_template_bd.bd1.name
    provider_connector_bd_template_name = data.mso_schema_template.t2.name
    provider_connector_bd_schema_id     = data.mso_schema.schema2.id
    provider_connector_bd_name          = data.mso_schema_template_bd.bd2.name
  }
  node_relationship {
    consumer_connector_bd_name = data.mso_schema_template_bd.bd1.name
    provider_connector_bd_name = data.mso_schema_template_bd.bd1.name
  }
  node_relationship {
    consumer_connector_bd_name = data.mso_schema_template_bd.bd1.name
    provider_connector_bd_name = data.mso_schema_template_bd.bd1.name
  }
}
