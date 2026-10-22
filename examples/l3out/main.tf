terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

# Set MSO_URL, MSO_USERNAME, and MSO_PASSWORD in the environment.
provider "mso" {
  platform = "nd"
}

data "mso_site" "example" {
  name = var.site_name
}

resource "mso_tenant" "example" {
  name = "example_l3out_tenant"

  site_associations {
    site_id = data.mso_site.example.id
  }
}

resource "mso_schema" "example" {
  name = "example_l3out_schema"

  template {
    name         = "example_l3out_schema_template"
    display_name = "example_l3out_schema_template"
    tenant_id    = mso_tenant.example.id
  }
}

resource "mso_schema_site" "example" {
  schema_id     = mso_schema.example.id
  template_name = "example_l3out_schema_template"
  site_id       = data.mso_site.example.id
}

resource "mso_schema_template_vrf" "example" {
  schema_id        = mso_schema.example.id
  template         = "example_l3out_schema_template"
  name             = "example_l3out_vrf"
  display_name     = "example_l3out_vrf"
  layer3_multicast = true

  depends_on = [mso_schema_site.example]
}

resource "mso_schema_template_deploy_ndo" "vrf" {
  schema_id           = mso_schema.example.id
  template_name       = "example_l3out_schema_template"
  force_apply         = ""
  undeploy_on_destroy = true

  depends_on = [mso_schema_template_vrf.example]
}

resource "mso_template" "fabric_policy" {
  template_name = "example_l3out_fabric_policy"
  template_type = "fabric_policy"
  sites         = [data.mso_site.example.id]
}

resource "mso_fabric_policies_l3_domain" "example" {
  template_id = mso_template.fabric_policy.id
  name        = "example_l3_domain"
}

resource "mso_template" "tenant_policy" {
  template_name = "example_l3out_tenant_policy"
  template_type = "tenant"
  tenant_id     = mso_tenant.example.id
  sites         = [data.mso_site.example.id]
}

resource "mso_tenant_policies_custom_qos_policy" "example" {
  template_id = mso_template.tenant_policy.id
  name        = "example_custom_qos"
}

resource "mso_tenant_policies_l3out_interface_routing_policy" "example" {
  template_id = mso_template.tenant_policy.id
  name        = "example_interface_routing"

  bfd_settings {
    admin_state = "enabled"
  }

  depends_on = [mso_tenant_policies_custom_qos_policy.example]
}

resource "mso_tenant_policies_netflow_exporter" "example" {
  template_id = mso_template.tenant_policy.id
  name        = "example_netflow_exporter"

  depends_on = [mso_tenant_policies_l3out_interface_routing_policy.example]
}

resource "mso_tenant_policies_netflow_record" "example" {
  template_id = mso_template.tenant_policy.id
  name        = "example_netflow_record"

  depends_on = [mso_tenant_policies_netflow_exporter.example]
}

resource "mso_tenant_policies_netflow_monitor" "ipv4" {
  template_id            = mso_template.tenant_policy.id
  name                   = "example_ipv4_monitor"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.example.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.example.uuid]
}

resource "mso_tenant_policies_netflow_monitor" "ipv6" {
  template_id            = mso_template.tenant_policy.id
  name                   = "example_ipv6_monitor"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.example.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.example.uuid]

  depends_on = [mso_tenant_policies_netflow_monitor.ipv4]
}

resource "mso_tenant_policies_netflow_monitor" "ce" {
  template_id            = mso_template.tenant_policy.id
  name                   = "example_ce_monitor"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.example.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.example.uuid]

  depends_on = [mso_tenant_policies_netflow_monitor.ipv6]
}

resource "mso_tenant_policies_netflow_monitor" "unspecified" {
  template_id            = mso_template.tenant_policy.id
  name                   = "example_unspecified_monitor"
  netflow_record_uuid    = mso_tenant_policies_netflow_record.example.uuid
  netflow_exporter_uuids = [mso_tenant_policies_netflow_exporter.example.uuid]

  depends_on = [mso_tenant_policies_netflow_monitor.ce]
}

resource "mso_template" "l3out" {
  template_name = "example_l3out_template"
  template_type = "l3out"
  tenant_id     = mso_tenant.example.id
  sites         = [data.mso_site.example.id]

  depends_on = [mso_schema_template_deploy_ndo.vrf]
}

resource "mso_l3out" "example" {
  template_id                  = mso_template.l3out.id
  name                         = "example_l3out"
  description                  = "L3Out managed with annotations and interface groups"
  vrf_uuid                     = mso_schema_template_vrf.example.uuid
  l3_domain                    = mso_fabric_policies_l3_domain.example.name
  target_dscp                  = "unspecified"
  pim_enabled                  = true
  import_route_control_enabled = true
  originate_default_route      = "in_addition"

  bgp = {
    enabled = true
  }

  ospf = {
    enabled                                       = true
    area_id                                       = "0.0.0.10"
    area_type                                     = "regular"
    cost                                          = 1
    originate_default_route_always                = false
    send_redistributed_lsa                        = true
    originate_summary_lsa                         = true
    suppress_forwarding_address_in_translated_lsa = false
  }

  annotations = {
    owner   = "network"
    purpose = "external-routing"
  }

  interface_groups = {
    edge = {
      description                   = "Edge interfaces"
      interface_routing_policy_uuid = mso_tenant_policies_l3out_interface_routing_policy.example.uuid
      custom_qos_policy_uuid        = mso_tenant_policies_custom_qos_policy.example.uuid
      qos_priority                  = "level6"

      netflow_monitor_uuids = {
        ipv4        = mso_tenant_policies_netflow_monitor.ipv4.uuid
        ipv6        = mso_tenant_policies_netflow_monitor.ipv6.uuid
        ce          = mso_tenant_policies_netflow_monitor.ce.uuid
        unspecified = mso_tenant_policies_netflow_monitor.unspecified.uuid
      }

      bfd = {
        enabled                = true
        authentication_enabled = true
        key_id                 = 20
        key                    = var.bfd_key
      }

      bfd_multi_hop = {
        enabled                = true
        authentication_enabled = true
        key_id                 = 30
        key                    = var.bfd_multi_hop_key
      }

      ospf = {
        enabled             = true
        authentication_type = "md5"
        key_id              = 10
        key                 = var.ospf_key
      }
    }
  }
}

data "mso_l3out" "example" {
  template_id = mso_template.l3out.id
  name        = mso_l3out.example.name
}
