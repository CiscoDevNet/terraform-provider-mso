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

# Tenant policies MLD snooping policy example

resource "mso_tenant_policies_mld_snooping_policy" "mld_policy" {
  template_id                = mso_template.tenant_template.id
  name                       = "mld_policy"
  description                = "Example description"
  admin_state                = "enabled"
  fast_leave_control         = true
  querier_control            = true
  querier_version            = "v2"
  query_interval             = 125
  query_response_interval    = 10
  last_member_query_interval = 1
  start_query_interval       = 31
  start_query_count          = 2
}
