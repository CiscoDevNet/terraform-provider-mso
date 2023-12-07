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

# AZURE CLOUD / AWS CLOUD
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

# ACI
data "aci_tenant" "ansible_test" {
  name = "ansible_test_cloud"
}

data "aci_cloud_l4_l7_native_load_balancer" "application_load_balancer" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "tf_application_load_balancer"
}

data "aci_cloud_l4_l7_native_load_balancer" "network_load_balancer" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "tf_network_load_balancer"
}

data "aci_cloud_l4_l7_third_party_device" "third_party_load_balancer" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "tf_third_party_load_balancer"
}

data "aci_cloud_l4_l7_third_party_device" "third_party_firewall" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "tf_third_party_firewall"
}


# MSO
provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO pwd>
  url      = "" # <MSO URL>
  insecure = true
  platform = "nd"
}

data "mso_tenant" "tf_tenant" {
  name = "ansible_test_cloud"
}

data "mso_site" "tf_site" {
  name = "azure_ansible_test_2"
}

resource "mso_schema" "schema_test" {
  name = "terraform_schema_cloud"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = data.mso_tenant.tf_tenant.id
  }
}

resource "mso_schema_template_service_graph" "test_sg_cloud" {
  schema_id          = mso_schema.schema_test.id
  template_name      = one(mso_schema.schema_test.template).name
  service_graph_name = "sgtf1"
  service_node {
    type = "other"
  }
  service_node {
    type = "load-balancer"
  }
  service_node {
    type = "firewall"
  }
  description = "Terraform Service Graph"
}

resource "mso_schema_site" "schema_site_1" {
  schema_id     = mso_schema.schema_test.id
  site_id       = data.mso_site.tf_site.id
  template_name = mso_schema_template_service_graph.test_sg_cloud.template_name
}

resource "mso_schema_site_service_graph" "test_sg_cloud_site" {
  schema_id          = mso_schema_site.schema_site_1.schema_id
  site_id            = mso_schema_site.schema_site_1.site_id
  template_name      = mso_schema_template_service_graph.test_sg_cloud.template_name
  service_graph_name = mso_schema_template_service_graph.test_sg_cloud.service_graph_name
  service_node {
    # for 1st item in the service graph list - network load balancer
    device_dn               = data.aci_cloud_l4_l7_native_load_balancer.network_load_balancer.id
    provider_connector_type = "redir"
    consumer_connector_type = "redir"
  }
  service_node {
    # for 2nd item in the service graph list - 3rd party load balancer
    device_dn          = data.aci_cloud_l4_l7_third_party_device.third_party_load_balancer.id
    consumer_interface = tolist(data.aci_cloud_l4_l7_third_party_device.third_party_load_balancer.interface_selectors)[0].name
    provider_interface = tolist(data.aci_cloud_l4_l7_third_party_device.third_party_load_balancer.interface_selectors)[0].name
  }
  service_node {
    # for 3rd item in the service graph list - 3rd party firewall
    device_dn                        = data.aci_cloud_l4_l7_third_party_device.third_party_firewall.id
    firewall_provider_connector_type = "snat_dnat"
    consumer_connector_type          = "redir"
    consumer_interface               = tolist(data.aci_cloud_l4_l7_third_party_device.third_party_load_balancer.interface_selectors)[0].name
    provider_interface               = tolist(data.aci_cloud_l4_l7_third_party_device.third_party_load_balancer.interface_selectors)[0].name
  }
}

data "mso_schema_site_service_graph" "test_sg_cloud_site" {
  schema_id          = mso_schema_site_service_graph.test_sg_cloud_site.schema_id
  site_id            = mso_schema_site_service_graph.test_sg_cloud_site.site_id
  template_name      = mso_schema_site_service_graph.test_sg_cloud_site.template_name
  service_graph_name = mso_schema_site_service_graph.test_sg_cloud_site.service_graph_name
}

output "example" {
  value = data.mso_schema_site_service_graph.test_sg_cloud_site
}

resource "mso_schema_template_service_graph" "test_sg_cloud_2" {
  schema_id          = mso_schema.schema_test.id
  template_name      = one(mso_schema.schema_test.template).name
  service_graph_name = "sgtf2"
  service_node {
    type = "other"
  }
}

resource "mso_schema_site_service_graph" "test_sg_cloud_site_2" {
  schema_id          = mso_schema_site.schema_site_1.schema_id
  site_id            = mso_schema_site.schema_site_1.site_id
  template_name      = mso_schema_template_service_graph.test_sg_cloud_2.template_name
  service_graph_name = mso_schema_template_service_graph.test_sg_cloud_2.service_graph_name
  service_node {
    # for 1st item in the service graph list - application load balancer
    device_dn = data.aci_cloud_l4_l7_native_load_balancer.application_load_balancer.id
  }
}