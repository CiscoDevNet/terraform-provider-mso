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

resource "mso_tenant" "tenant" {
  name         = "ServiceChainTenant"
  display_name = "ServiceChainTenant"
}

resource "mso_template" "device_template" {
  template_name = "DeviceTemplateForSC"
  template_type = "service_device"
  tenant_id     = mso_tenant.tenant.id
}

resource "mso_schema" "schema" {
  name = "SchemaForServiceChaining"
  template {
    name          = "Template1"
    display_name  = "Template1"
    tenant_id     = mso_tenant.tenant.id
    template_type = "aci_multi_site"
  }
}

resource "mso_schema_template_vrf" "vrf" {
  schema_id    = mso_schema.schema.id
  template     = "Template1"
  name         = "SC_VRF"
  display_name = "SC_VRF"
}

resource "mso_schema_template_bd" "bd1" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  name          = "SC_BD1"
  vrf_name      = mso_schema_template_vrf.vrf.name
}

resource "mso_schema_template_bd" "bd2" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  name          = "SC_BD2"
  vrf_name      = mso_schema_template_vrf.vrf.name
}

resource "mso_service_device_cluster" "fw_device" {
  template_id = mso_template.device_template.id
  name        = "FirewallCluster"
  device_mode = "layer3"
  device_type = "firewall"

  interface_properties {
    name    = "fw_interface"
    bd_uuid = mso_schema_template_bd.bd1.uuid
  }
}

resource "mso_service_device_cluster" "lb_device" {
  template_id = mso_template.device_template.id
  name        = "LoadBalancerCluster"
  device_mode = "layer3"
  device_type = "loadBalancer"

  interface_properties {
    name    = "lb_prov_if"
    bd_uuid = mso_schema_template_bd.bd1.uuid
  }
  interface_properties {
    name    = "lb_cons_if"
    bd_uuid = mso_schema_template_bd.bd2.uuid
  }
}

resource "mso_schema_template_contract" "contract" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  contract_name = "WebAppContract"
  display_name  = "WebAppContract"
  scope         = "context"
}

# Main resource for the Service Chain
resource "mso_schema_template_contract_service_chaining" "chain" {
  schema_id     = mso_schema.schema.id
  template_name = "Template1"
  contract_name = mso_schema_template_contract.contract.contract_name
  node_filter   = "allow-all"

  service_nodes {
    name        = "firewall-node"
    device_type = "firewall"
    device_ref  = mso_service_device_cluster.fw_device.uuid

    consumer_connector {
      interface_name = "fw_interface"
    }
    provider_connector {
      interface_name = "fw_interface"
    }
  }

  service_nodes {
    name        = "loadbalancer-node"
    device_type = "loadBalancer"
    device_ref  = mso_service_device_cluster.lb_device.uuid

    consumer_connector {
      interface_name = "lb_cons_if"
    }
    provider_connector {
      interface_name = "lb_prov_if"
    }
  }
}
