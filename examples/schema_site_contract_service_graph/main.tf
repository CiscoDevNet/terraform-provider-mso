terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO pwd>
  url      = "" # <MSO URL>
  insecure = true
  platform = "nd"
}

# ACI config begins

data "aci_tenant" "ansible_test" {
  name = "ansible_test"
}

data "aci_l4_l7_device" "l4_l7_fw" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "l4_l7_fw"
}

data "aci_l4_l7_device" "l4_l7_adc_lb" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "l4_l7_adc_lb"
}

data "aci_l4_l7_device" "l4_l7_others" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "l4_l7_others"
}

data "aci_l4_l7_logical_interface" "l4_l7_fw_prov_prov_inf" {
  l4_l7_device_dn = data.aci_l4_l7_device.l4_l7_fw.id
  name            = "prov_inf"
}
data "aci_l4_l7_logical_interface" "l4_l7_fw_prov_cons_inf" {
  l4_l7_device_dn = data.aci_l4_l7_device.l4_l7_fw.id
  name            = "cons_inf"
}

data "aci_l4_l7_logical_interface" "l4_l7_adc_lb_prov_inf" {
  l4_l7_device_dn = data.aci_l4_l7_device.l4_l7_adc_lb.id
  name            = "prov_inf"
}

data "aci_l4_l7_logical_interface" "l4_l7_adc_lb_cons_inf" {
  l4_l7_device_dn = data.aci_l4_l7_device.l4_l7_adc_lb.id
  name            = "cons_inf"
}

data "aci_l4_l7_logical_interface" "l4_l7_others_prov_inf" {
  l4_l7_device_dn = data.aci_l4_l7_device.l4_l7_others.id
  name            = "prov_inf"
}

data "aci_l4_l7_logical_interface" "l4_l7_others_cons_inf" {
  l4_l7_device_dn = data.aci_l4_l7_device.l4_l7_others.id
  name            = "cons_inf"
}

# ND Config begins

data "mso_site" "ansible_test" {
  name = "ansible_test"
}

data "mso_tenant" "ansible_test" {
  name = "ansible_test"
}

resource "mso_schema" "tf_schema_sg" {
  name = "tf_schema_sg"
  template {
    name         = "template1"
    display_name = "template1"
    tenant_id    = data.mso_tenant.ansible_test.id
  }
}

resource "mso_schema_template_vrf" "vrf1" {
  schema_id    = mso_schema.tf_schema_sg.id
  template     = one(mso_schema.tf_schema_sg.template).name
  name         = "vrf1"
  display_name = "vrf1"
}

resource "mso_schema_template_vrf" "vrf2" {
  schema_id    = mso_schema.tf_schema_sg.id
  template     = one(mso_schema.tf_schema_sg.template).name
  name         = "vrf2"
  display_name = "vrf2"
}

resource "mso_schema_template_bd" "bd1" {
  schema_id     = mso_schema.tf_schema_sg.id
  template_name = one(mso_schema.tf_schema_sg.template).name
  name          = "bd1"
  display_name  = "bd1"
  vrf_name      = mso_schema_template_vrf.vrf1.name
  arp_flooding  = true
}

resource "mso_schema_template_bd" "bd2" {
  schema_id     = mso_schema.tf_schema_sg.id
  template_name = one(mso_schema.tf_schema_sg.template).name
  name          = "bd2"
  display_name  = "bd2"
  vrf_name      = mso_schema_template_vrf.vrf2.name
  arp_flooding  = true
}

resource "mso_schema_template_filter_entry" "filter_entry1" {
  schema_id          = mso_schema.tf_schema_sg.id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  name               = "filter1"
  display_name       = "filter1"
  entry_name         = "filter_entry1"
  entry_display_name = "filter_entry1"
}

resource "mso_schema_template_contract" "template_c1" {
  schema_id     = mso_schema.tf_schema_sg.id
  template_name = one(mso_schema.tf_schema_sg.template).name
  contract_name = "contract1"
  filter_relationship {
    filter_name = mso_schema_template_filter_entry.filter_entry1.name
    filter_type = "bothWay"
  }
}

resource "mso_schema_template_vrf_contract" "template_vrf_contract" {
  schema_id         = mso_schema.tf_schema_sg.id
  template_name     = one(mso_schema.tf_schema_sg.template).name
  vrf_name          = mso_schema_template_vrf.vrf1.name
  relationship_type = "provider"
  contract_name     = mso_schema_template_contract.template_c1.contract_name
}

resource "mso_schema_template_service_graph" "template_sg1" {
  schema_id          = mso_schema.tf_schema_sg.id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  service_graph_name = "service_graph1"
  service_node {
    type = "firewall"
  }
  service_node {
    type = "load-balancer"
  }
  service_node {
    type = "other"
  }
}

resource "mso_schema_template_l3out" "template_l3out" {
  schema_id     = mso_schema.tf_schema_sg.id
  template_name = one(mso_schema.tf_schema_sg.template).name
  l3out_name    = "l3out"
  display_name  = "l3out"
  vrf_name      = mso_schema_template_vrf.vrf1.name
}

resource "mso_schema_template_external_epg" "ext_epg" {
  schema_id         = mso_schema.tf_schema_sg.id
  template_name     = one(mso_schema.tf_schema_sg.template).name
  external_epg_name = "ext_epg"
  display_name      = "ext_epg"
  vrf_name          = mso_schema_template_vrf.vrf1.name
  l3out_name        = mso_schema_template_l3out.template_l3out.l3out_name
}

resource "mso_schema_template_external_epg_contract" "ext_epg_c1_cons" {
  schema_id         = mso_schema.tf_schema_sg.id
  template_name     = one(mso_schema.tf_schema_sg.template).name
  external_epg_name = mso_schema_template_external_epg.ext_epg.external_epg_name
  relationship_type = "consumer"
  contract_name     = mso_schema_template_contract.template_c1.contract_name
}

resource "mso_schema_template_external_epg_contract" "ext_epg_c1_prov" {
  schema_id         = mso_schema.tf_schema_sg.id
  template_name     = one(mso_schema.tf_schema_sg.template).name
  external_epg_name = mso_schema_template_external_epg.ext_epg.external_epg_name
  relationship_type = "provider"
  contract_name     = mso_schema_template_contract.template_c1.contract_name
  depends_on = [mso_schema_template_external_epg_contract.ext_epg_c1_cons
  ]
}

resource "mso_schema_template_contract_service_graph" "template_contract_sg" {
  schema_id          = mso_schema.tf_schema_sg.id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  contract_name      = mso_schema_template_contract.template_c1.contract_name
  service_graph_name = mso_schema_template_service_graph.template_sg1.service_graph_name
  node_relationship {
    consumer_connector_bd_name = mso_schema_template_bd.bd1.name
    provider_connector_bd_name = mso_schema_template_bd.bd2.name
  }
  node_relationship {
    consumer_connector_bd_name = mso_schema_template_bd.bd1.name
    provider_connector_bd_name = mso_schema_template_bd.bd2.name
  }
  node_relationship {
    consumer_connector_bd_name = mso_schema_template_bd.bd1.name
    provider_connector_bd_name = mso_schema_template_bd.bd2.name
  }
  depends_on = [mso_schema_template_external_epg_contract.ext_epg_c1_prov, mso_schema_template_external_epg_contract.ext_epg_c1_cons
  ]
}

resource "mso_schema_site" "site1" {
  schema_id     = mso_schema.tf_schema_sg.id
  template_name = one(mso_schema.tf_schema_sg.template).name
  site_id       = data.mso_site.ansible_test.id
}

# Note: This resource(mso_schema_site_service_graph) is supported only for NDO 4.1.1i and above.
# Deletion of site Service Graph is not supported by the API. 
# Site Service Graph will be removed when site is disassociated from the template or when Service Graph is removed at the template level.
resource "mso_schema_site_service_graph" "site_service_graph" {
  schema_id          = mso_schema.tf_schema_sg.id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  site_id            = mso_schema_site.site1.site_id
  service_graph_name = mso_schema_template_contract_service_graph.template_contract_sg.service_graph_name
  service_node {
    device_dn = data.aci_l4_l7_device.l4_l7_fw.id
  }
  service_node {
    device_dn = data.aci_l4_l7_device.l4_l7_adc_lb.id
  }
  service_node {
    device_dn = data.aci_l4_l7_device.l4_l7_others.id
  }
}

resource "mso_schema_site_contract_service_graph" "site_contract_service_graph" {
  schema_id          = mso_schema.tf_schema_sg.id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  site_id            = mso_schema_site_service_graph.site_service_graph.site_id
  contract_name      = mso_schema_template_contract.template_c1.contract_name
  service_graph_name = mso_schema_template_contract_service_graph.template_contract_sg.service_graph_name
  node_relationship {
    provider_connector_cluster_interface = data.aci_l4_l7_logical_interface.l4_l7_fw_prov_prov_inf.name
    consumer_connector_cluster_interface = data.aci_l4_l7_logical_interface.l4_l7_fw_prov_cons_inf.name
  }
  node_relationship {
    provider_connector_cluster_interface = data.aci_l4_l7_logical_interface.l4_l7_adc_lb_prov_inf.name
    consumer_connector_cluster_interface = data.aci_l4_l7_logical_interface.l4_l7_adc_lb_cons_inf.name
    consumer_subnet_ips = [
      "1.1.1.1/24",
      "2.2.2.2/24"
    ]
  }
  node_relationship {
    provider_connector_cluster_interface = data.aci_l4_l7_logical_interface.l4_l7_others_prov_inf.name
    consumer_connector_cluster_interface = data.aci_l4_l7_logical_interface.l4_l7_others_cons_inf.name
  }
}
