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

data "mso_site" "ansible_test" {
  name = "ansible_test"
}

resource "mso_tenant" "mso_tenant" {
  name         = "mso_tenant"
  display_name = "mso_tenant"
  site_associations {
    site_id = data.mso_site.ansible_test.id
  }
}

resource "mso_schema" "mso_schema" {
  name = "mso_schema"
  template {
    name         = "mso_schema_template"
    display_name = "mso_schema_template"
    tenant_id    = mso_tenant.mso_tenant.id
  }
}

resource "mso_schema_template_vrf" "example_vrf" {
  name         = "example_vrf"
  display_name = "example_vrf"
  schema_id    = mso_schema.mso_schema.id
  template     = "mso_schema_template"
}

resource "mso_schema_template_bd" "example_bd" {
  schema_id              = mso_schema.mso_schema.id
  template_name          = "mso_schema_template"
  name                   = "example_bd"
  display_name           = "example_bd"
  layer2_unknown_unicast = "proxy"
  vrf_name               = mso_schema_template_vrf.example_vrf.name
}

resource "mso_template" "tenant_policy_template" {
  template_name = "tenant_policy_template"
  template_type = "tenant"
  tenant_id     = mso_tenant.mso_tenant.id
}

resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_monitoring_policy" {
  template_id        = mso_template.tenant_policy_template.id
  name               = "ipsla_monitoring_policy"
  sla_type           = "http"
  destination_port   = 80
  http_version       = "HTTP11"
  http_uri           = "/example"
  sla_frequency      = 120
  detect_multiplier  = 4
  request_data_size  = 64
  type_of_service    = 18
  operation_timeout  = 100
  threshold          = 100
  ipv6_traffic_class = 255
}

resource "mso_tenant_policies_ipsla_track_list" "ipsla_track_list" {
  template_id    = mso_template.tenant_policy_template.id
  name           = "ipsla_track_list"
  description    = "Terraform test IPSLA Track List"
  threshold_down = 11
  threshold_up   = 12
  type           = "weight"
  members {
    destination_ip               = "1.1.1.1"
    ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla_monitoring_policy.uuid
    scope_type                   = "bd"
    scope_uuid                   = mso_schema_template_bd.example_bd.uuid
  }
}
