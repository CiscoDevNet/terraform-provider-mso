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
  platform = "nd"
}

data "mso_tenant" "tenant_test" {
  name         = "ansible_test"
  display_name = "ansible_test"
}

resource "mso_schema" "schema_test" {
  name = "terraform_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = data.mso_tenant.tenant_test.id
  }
}

resource "mso_schema_template_service_graph" "test_sg" {
  schema_id          = mso_schema.schema_test.id
  template_name      = one(mso_schema.schema_test.template).name
  service_graph_name = "sgtf1"
  service_node {
    type = "firewall"
  }
  service_node {
    type = "firewall"
  }
  description = "Terraform Service Graph"
}

data "mso_schema_template_service_graph" "test_sg" {
  schema_id          = mso_schema_template_service_graph.test_sg.schema_id
  template_name      = mso_schema_template_service_graph.test_sg.template_name
  service_graph_name = mso_schema_template_service_graph.test_sg.service_graph_name
}

output "example" {
  value = data.mso_schema_template_service_graph.test_sg
}
