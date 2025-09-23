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

resource "mso_template" "device_template" {
  name          = "test_device_template"
  display_name  = "test_device_template"
  template_type = "service_device"
}

resource "mso_template" "tenant_template" {
  name          = "test_tenant_template_for_device"
  display_name  = "test_tenant_template_for_device"
  template_type = "tenant_policy"
}

resource "mso_schema" "schema1" {
  name          = "SchemaForDeviceClusterTest"
  display_name  = "SchemaForDeviceClusterTest"
  template_name = "Template1"
}

resource "mso_schema_template_bd" "bd1" {
  schema_id    = mso_schema.schema1.id
  template_name = "Template1"
  name         = "test_bd_1"
  display_name = "test_bd_1"
}

resource "mso_schema_template_external_epg" "epg1" {
  schema_id    = mso_schema.schema1.id
  template_name = "Template1"
  name         = "test_epg_1"
  display_name = "test_epg_1"
}

resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla1" {
  template_id = mso_template.tenant_template.id
  name        = "test_ipsla_for_device"
  sla_type    = "icmp"
}

# Main resource for the Service Device Cluster
resource "mso_service_device_cluster" "cluster" {
  template_id = mso_template.device_template.id
  name        = "test_device_cluster"
  device_mode = "layer3"
  device_type = "firewall"

  interface_properties {
    name                         = "interface1"
    external_epg_uuid            = mso_schema_template_external_epg.epg1.uuid
    ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
    load_balance_hashing         = "source_ip"
    min_threshold                = 10
    max_threshold                = 90
    threshold_down_action        = "permit"
    preferred_group              = true
  }

  interface_properties {
    name                         = "interface2"
    bd_uuid                      = mso_schema_template_bd.bd1.uuid
    ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
    load_balance_hashing         = "destination_ip"
  }
}
