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
  match_type  = "prefix_and_routing_policy"
}

resource "mso_tenant_policies_set_rule_policy" "set_rule_policy" {
  template_id = mso_template.tenant_template.id
  name        = "example_set_rule_policy"
  set_metric  = 100

  set_route_tag {
    value = 0
  }
  set_weight {
    value = 0
  }
  set_preference {
    value = 0
  }
}

resource "mso_tenant_policies_match_rule_policy" "match_rule_policy" {
  template_id = mso_template.tenant_template.id
  name        = "example_match_rule_policy"

  match_prefix {
    prefix = "10.0.0.0/8"
  }
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.route_map_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.route_map_policy.uuid
  name                  = "example_context"
  description           = "Example Route Control Context"
  order                 = 1
  action                = "permit"
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.set_rule_policy.uuid
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.match_rule_policy.uuid]
}

data "mso_tenant_policies_route_map_policy_route_control_context" "context" {
  template_id           = mso_tenant_policies_route_map_policy_route_control_context.context.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control_context.context.route_map_policy_uuid
  name                  = mso_tenant_policies_route_map_policy_route_control_context.context.name
}
