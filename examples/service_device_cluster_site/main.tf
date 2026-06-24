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
  platform = "nd"
}

data "mso_site" "site1" {
  name = "site1"
}

resource "mso_tenant" "tenant1" {
  name = "tenant1"

  site_associations {
    site_id = data.mso_site.site1.id
  }
}

resource "mso_schema" "schema1" {
  name = "schema1"
  template {
    name         = "template1"
    display_name = "template1"
    tenant_id    = mso_tenant.tenant1.id
  }
}

resource "mso_schema_site" "schema_site1" {
  schema_id     = mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = one(mso_schema.schema1.template).name
}

resource "mso_schema_template_vrf" "vrf1" {
  schema_id    = mso_schema.schema1.id
  template     = one(mso_schema.schema1.template).name
  name         = "vrf1"
  display_name = "vrf1"
  depends_on   = [mso_schema_site.schema_site1]
}

resource "mso_schema_template_bd" "bd1" {
  schema_id      = mso_schema.schema1.id
  template_name  = one(mso_schema.schema1.template).name
  name           = "bd1"
  display_name   = "bd1"
  vrf_name       = mso_schema_template_vrf.vrf1.name
  arp_flooding   = true
  layer2_stretch = false
}

resource "mso_schema_template_deploy_ndo" "schema1" {
  schema_id           = mso_schema.schema1.id
  template_name       = one(mso_schema.schema1.template).name
  site_ids            = [data.mso_site.site1.id]
  undeploy_on_destroy = true
  force_apply         = ""
  depends_on          = [mso_schema_template_deploy_ndo.tenant_template]
}

resource "mso_template" "fabric_policy" {
  template_name = "fabric_policy_template"
  template_type = "fabric_policy"
  sites         = [data.mso_site.site1.id]
}

resource "mso_fabric_policies_vlan_pool" "vlan_pool" {
  template_id = mso_template.fabric_policy.id
  name        = "vlan_pool"
  vlan_range {
    from = 200
    to   = 250
  }
}

resource "mso_fabric_policies_physical_domain" "physical_domain" {
  template_id    = mso_template.fabric_policy.id
  name           = "physical_domain"
  vlan_pool_uuid = mso_fabric_policies_vlan_pool.vlan_pool.uuid
}

resource "mso_schema_template_deploy_ndo" "fabric_policy" {
  template_id         = mso_fabric_policies_physical_domain.physical_domain.template_id
  template_type       = "fabric_policy"
  undeploy_on_destroy = true
  site_ids            = [data.mso_site.site1.id]
  force_apply         = ""
}

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id     = mso_tenant.tenant1.id
  sites         = [data.mso_site.site1.id]
}

resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
  template_id = mso_template.tenant_template.id
  name        = "ipsla_policy"
  sla_type    = "icmp"
}

resource "mso_schema_template_deploy_ndo" "tenant_template" {
  template_id         = mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy.template_id
  template_type       = "tenant"
  undeploy_on_destroy = true
  site_ids            = [data.mso_site.site1.id]
  force_apply         = ""
  depends_on          = [mso_schema_template_deploy_ndo.fabric_policy]
}

resource "mso_template" "device_template" {
  template_name = "device_template"
  template_type = "service_device"
  tenant_id     = mso_tenant.tenant1.id
  sites         = [data.mso_site.site1.id]
  depends_on    = [mso_schema_template_deploy_ndo.fabric_policy, mso_schema_template_deploy_ndo.tenant_template, mso_schema_template_deploy_ndo.schema1]
}

resource "mso_service_device_cluster" "cluster" {
  template_id = mso_template.device_template.id
  name        = "device_cluster"
  device_mode = "layer3"
  device_type = "firewall"

  interface_properties {
    name                         = "Internal"
    bd_uuid                      = mso_schema_template_bd.bd1.uuid
    ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla_policy.uuid
    config_static_mac            = true
    is_backup_redirect_ip        = true
    resilient_hashing            = true
    tag_based_sorting            = true
    min_threshold                = 3
    max_threshold                = 100
    threshold_down_action        = "permit"
  }

  interface_properties {
    name    = "External1"
    bd_uuid = mso_schema_template_bd.bd1.uuid
  }

  interface_properties {
    name    = "External2"
    bd_uuid = mso_schema_template_bd.bd1.uuid
  }
}

resource "mso_service_device_cluster_site" "cluster_site" {
  template_id = mso_template.device_template.id
  site_id     = data.mso_site.site1.id
  name        = mso_service_device_cluster.cluster.name

  domain_type = "physicalDomain"
  domain_name = mso_fabric_policies_physical_domain.physical_domain.name

  interfaces {
    name = "Internal"
    vlan = 201

    fabric_to_device_connectivity {
      pod_id    = "1"
      node_id   = ["101"]
      path      = "eth1/20"
      port_type = "port"
    }

    pbr_destinations {
      ip                     = "10.0.0.10"
      mac                    = "00:22:BD:F8:19:FE"
      pod_id                 = "1"
      additional_tracking_ip = "0.0.0.0"
      tag                    = "internal"
    }
  }

  interfaces {
    name = "External1"
    vlan = 202

    fabric_to_device_connectivity {
      pod_id    = "1"
      node_id   = ["101"]
      path      = "eth1/25"
      port_type = "port"
    }
  }

  interfaces {
    name = "External2"
    vlan = 203

    fabric_to_device_connectivity {
      pod_id    = "1"
      node_id   = ["102"]
      path      = "eth1/12"
      port_type = "port"
    }
  }
}

data "mso_service_device_cluster_site" "cluster_site" {
  template_id = mso_service_device_cluster_site.cluster_site.template_id
  site_id     = mso_service_device_cluster_site.cluster_site.site_id
  name        = mso_service_device_cluster_site.cluster_site.name
}
