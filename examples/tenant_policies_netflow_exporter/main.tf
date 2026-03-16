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

data "mso_tenant" "example_tenant" {
  name = "example_tenant"
}

# Tenant template example

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id     = data.mso_tenant.example_tenant.id
}

# NetFlow Exporter example

resource "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
  template_id = mso_template.tenant_template.id
  name        = "netflow_exporter_1"
}

# NetFlow Exporter data source example

data "mso_tenant_policies_netflow_exporter" "netflow_exporter" {
  template_id = mso_template.tenant_template.id
  name        = mso_tenant_policies_netflow_exporter.netflow_exporter.name
}