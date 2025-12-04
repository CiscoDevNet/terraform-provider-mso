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

resource "mso_tenant_policies_route_map_policy_multicast" "state_limit" {
  template_id = mso_template.tenant_template.id
  name        = "tf_test_state_limit"
  description = "Terraform test Route Map Policy for Multicast"
  route_map_multicast_entries {
    order               = 1
    group_ip            = "226.2.2.2/8"
    source_ip           = "1.1.1.1/1"
    rendezvous_point_ip = "1.1.1.2"
    action              = "permit"
  }
}

resource "mso_tenant_policies_route_map_policy_multicast" "report_policy" {
  template_id = mso_tenant_policies_route_map_policy_multicast.state_limit.template_id
  name        = "tf_test_report_policy"
  description = "Terraform test Route Map Policy for Multicast"
  route_map_multicast_entries {
    order               = 1
    group_ip            = "226.2.2.2/8"
    source_ip           = "1.1.1.1/1"
    rendezvous_point_ip = "1.1.1.2"
    action              = "permit"
  }
}

resource "mso_tenant_policies_route_map_policy_multicast" "static_report" {
  template_id = mso_tenant_policies_route_map_policy_multicast.report_policy.template_id
  name        = "tf_test_static_report"
  description = "Terraform test Route Map Policy for Multicast"
  route_map_multicast_entries {
    order               = 1
    group_ip            = "226.2.2.2/8"
    source_ip           = "1.1.1.1/1"
    rendezvous_point_ip = "1.1.1.2"
    action              = "permit"
  }
}

# tenant policies igmp interface policy example

resource "mso_tenant_policies_igmp_interface_policy" "igmp_policy" {
  template_id                  = mso_template.tenant_template.id
  name                         = "test_igmp_interface_policy"
  description                  = "With Route Maps"
  igmp_version                 = "v3"
  state_limit_route_map_uuid   = mso_tenant_policies_route_map_policy_multicast.state_limit.uuid
  report_policy_route_map_uuid = mso_tenant_policies_route_map_policy_multicast.report_policy.uuid
  static_report_route_map_uuid = mso_tenant_policies_route_map_policy_multicast.static_report.uuid
  maximum_multicast_entries    = 5000000
}
