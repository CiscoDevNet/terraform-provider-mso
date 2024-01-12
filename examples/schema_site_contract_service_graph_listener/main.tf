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

# Azure Cloud Network Controller Sample Configuration
resource "aci_tenant" "tf_tenant" {
  name = "tf_tenant"
}

# Cloud Subnet setup part
resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "vrf_1"
}

resource "aci_cloud_context_profile" "ctx1" {
  name                     = "tf_ctx1"
  tenant_dn                = aci_tenant.tf_tenant.id
  primary_cidr             = "10.1.0.0/16"
  region                   = "westus"
  cloud_vendor             = "azure"
  relation_cloud_rs_to_ctx = aci_vrf.vrf1.id
  hub_network              = "uni/tn-infra/gwrouterp-default"
}

resource "aci_cloud_cidr_pool" "cloud_cidr_pool" {
  cloud_context_profile_dn = aci_cloud_context_profile.ctx1.id
  addr                     = "10.1.0.0/16"
}

data "aci_cloud_provider_profile" "cloud_profile" {
  vendor = "azure"
}

data "aci_cloud_providers_region" "cloud_region" {
  cloud_provider_profile_dn = data.aci_cloud_provider_profile.cloud_profile.id
  name                      = "westus"
}

data "aci_cloud_availability_zone" "region_availability_zone" {
  cloud_providers_region_dn = data.aci_cloud_providers_region.cloud_region.id
  name                      = "default"
}

resource "aci_cloud_subnet" "cloud_subnet" {
  cloud_cidr_pool_dn = aci_cloud_cidr_pool.cloud_cidr_pool.id
  ip                 = "10.1.1.0/24"
  usage              = "gateway"
  zone               = data.aci_cloud_availability_zone.region_availability_zone.id
  scope              = ["shared", "private", "public"]
}

# AAA Domain setup part
resource "aci_aaa_domain" "aaa_domain_1" {
  name = "aaa_domain_1"
}

resource "aci_aaa_domain" "aaa_domain_2" {
  name = "aaa_domain_2"
}

# Application Load Balancer
resource "aci_cloud_l4_l7_native_load_balancer" "cloud_native_alb" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "cloud_native_alb"
  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id,
    aci_aaa_domain.aaa_domain_2.id
  ]
  relation_cloud_rs_ldev_to_cloud_subnet = [
    aci_cloud_subnet.cloud_subnet.id
  ]
  active_active                 = "no"
  allow_all                     = "no"
  auto_scaling                  = "no"
  context_aware                 = "multi-Context"
  device_type                   = "CLOUD"
  function_type                 = "GoTo"
  is_copy                       = "no"
  is_instantiation              = "no"
  is_static_ip                  = "no"
  managed                       = "no"
  mode                          = "legacy-Mode"
  promiscuous_mode              = "no"
  scheme                        = "internal"
  size                          = "medium"
  sku                           = "standard"
  service_type                  = "NATIVELB"
  target_mode                   = "primary"
  trunking                      = "no"
  cloud_l4l7_load_balancer_type = "application"
  instance_count                = "2"
  max_instance_count            = "10"
  min_instance_count            = "5"
}

# Network Load Balancer
resource "aci_cloud_l4_l7_native_load_balancer" "cloud_native_nlb" {
  tenant_dn                     = aci_tenant.tf_tenant.id
  name                          = "cloud_native_nlb"
  service_type                  = "NATIVELB"
  cloud_l4l7_load_balancer_type = "network"
  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id
  ]
  relation_cloud_rs_ldev_to_cloud_subnet = [
    aci_cloud_subnet.cloud_subnet.id
  ]
}


# ND Azure Cloud Network Controller - Schema Site Contract Service Graph Listener Sample Configuration
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

resource "mso_schema_template_anp" "anp1" {
  schema_id    = mso_schema.tf_schema_sg.id
  template     = one(mso_schema.tf_schema_sg.template).name
  name         = "anp1"
  display_name = "anp1"
}

resource "mso_schema_template_anp_epg" "epg1" {
  schema_id     = mso_schema.tf_schema_sg.id
  template_name = one(mso_schema.tf_schema_sg.template).name
  anp_name      = mso_schema_template_anp.anp1.name
  name          = "epg1"
  bd_name       = mso_schema_template_bd.bd1.name
  vrf_name      = mso_schema_template_vrf.vrf1.name
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

resource "mso_schema_template_anp_epg_contract" "contract_epg_provider" {
  schema_id         = mso_schema.tf_schema_sg.id
  template_name     = one(mso_schema.tf_schema_sg.template).name
  anp_name          = mso_schema_template_anp.anp1.name
  epg_name          = mso_schema_template_anp_epg.epg1.name
  contract_name     = mso_schema_template_contract.template_c1.contract_name
  relationship_type = "provider"
}

resource "mso_schema_template_vrf_contract" "template_vrf_contract" {
  schema_id         = mso_schema.tf_schema_sg.id
  template_name     = one(mso_schema.tf_schema_sg.template).name
  vrf_name          = mso_schema_template_vrf.vrf1.name
  relationship_type = "provider"
  contract_name     = mso_schema_template_contract.template_c1.contract_name
}

# I am currently here
resource "mso_schema_template_service_graph" "template_sg1" {
  schema_id          = mso_schema.tf_schema_sg.id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  service_graph_name = "service_graph1"
  service_node {
    type = "load-balancer"
  }
  service_node {
    type = "load-balancer"
  }
}

# If we want to use an external EPG
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
    device_dn = aci_cloud_l4_l7_native_load_balancer.cloud_native_alb.id
  }
  service_node {
    device_dn = aci_cloud_l4_l7_native_load_balancer.cloud_native_nlb.id
  }
}

resource "mso_schema_site_contract_service_graph" "site_contract_service_graph" {
  schema_id          = mso_schema.tf_schema_sg.id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  site_id            = mso_schema_site_service_graph.site_service_graph.site_id
  contract_name      = mso_schema_template_contract.template_c1.contract_name
  service_graph_name = mso_schema_template_contract_service_graph.template_contract_sg.service_graph_name
}

# Application Load Balancer - Config example
resource "mso_schema_site_contract_service_graph_listener" "application_load_balancer_node1" {
  contract_name      = mso_schema_site_contract_service_graph.site_contract_service_graph.contract_name
  listener_name      = "example"
  port               = 443
  protocol           = "https"
  schema_id          = mso_schema.tf_schema_sg.id
  security_policy    = "default"
  service_node_index = 1
  site_id            = mso_schema_site_service_graph.site_service_graph.site_id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  rules {
    action_type       = "redirect"
    content_type      = "text_plain"
    name              = "rule_1"
    port              = 80
    priority          = 1
    protocol          = "http"
    redirect_code     = "permanently_moved"
    redirect_port     = 80
    redirect_protocol = "http"
    response_code     = "204"
    target_ip_type    = "unspecified"
    url_type          = "original"
    health_check {
      host                = "3.3.3.3"
      interval            = 30
      path                = "/"
      port                = 443
      protocol            = "https"
      success_code        = "200"
      timeout             = 30
      unhealthy_threshold = 3
      use_host_from_rule  = false
    }
  }
  ssl_certificates {
    certificate_store = "default"
    name              = "ssl_certificate_key_ring"
    # Steps to create Key Ring
    # 1. Administrative -> Security -> Certificate Authorities
    # 2. Administrative -> Security -> Key Rings
    target_dn = "uni/tn-azure_tenant/certstore" # Certificate Authority -> Key Ring
  }
}

# frontend_ip attribute configuration only available for Network Load Balancer
resource "mso_schema_site_contract_service_graph_listener" "network_load_balancer_node2" {
  contract_name = mso_schema_template_contract.template_c1.contract_name
  # Steps to configure Frontend IP Name
  # 1. Application Management -> Services -> Create a "Network Load Balancer" - L4-L7 Device -> Advanced Settings -> Additional Frontend IPs -> Frontend IP Names
  frontend_ip_dn     = "uni/tn-azure_tenant/clb-nlb/vip-2.2.2.2"
  listener_name      = "example"
  port               = 80
  protocol           = "udp"
  schema_id          = mso_schema.tf_schema_sg.id
  security_policy    = "default"
  service_node_index = 0
  site_id            = mso_schema_site_service_graph.site_service_graph.site_id
  template_name      = one(mso_schema.tf_schema_sg.template).name
  rules {
    action_type       = "forward"
    content_type      = "text_plain"
    name              = "default"
    port              = 80
    priority          = 0
    protocol          = "udp"
    redirect_code     = "permanently_moved"
    redirect_port     = 0
    redirect_protocol = "inherit"
    target_ip_type    = "unspecified"
    url_type          = "original"
    health_check {
      interval            = 5
      port                = 80
      protocol            = "tcp"
      success_code        = "200"
      timeout             = 0
      unhealthy_threshold = 2
      use_host_from_rule  = false
    }
    provider_epg_ref {
      anp_name      = mso_schema_template_anp.anp1.name
      epg_name      = mso_schema_template_anp_epg_contract.contract_epg_provider.epg_name
      schema_id     = mso_schema.tf_schema_sg.id
      template_name = one(mso_schema.tf_schema_sg.template).name
    }
  }
}
