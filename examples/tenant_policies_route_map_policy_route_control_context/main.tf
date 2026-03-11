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

# tenant template example

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id     = data.mso_tenant.example_tenant.id
}

resource "mso_tenant_policies_route_map_policy_route_control" "route_map_policy" {
  template_id = mso_template.tenant_template.id
  name        = "example_route_map_policy"
  description = "Example Route Map Policy"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "context" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.route_map_policy.id
  name        = "example_context"
  description = "Example Route Control Context"
  order       = 1
  action      = "permit"
}
