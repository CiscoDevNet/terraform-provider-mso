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

# NetFlow Record example

resource "mso_tenant_policies_netflow_record" "netflow_record" {
  template_id      = mso_template.tenant_template.id
  name             = "netflow_record_1"
  description      = "My NetFlow Record"
  match_parameters = ["destination_ip", "source_ip", "ip_protocol"]
}

# NetFlow Record data source example

data "mso_tenant_policies_netflow_record" "netflow_record" {
  template_id = mso_template.tenant_template.id
  name        = mso_tenant_policies_netflow_record.netflow_record.name
}
