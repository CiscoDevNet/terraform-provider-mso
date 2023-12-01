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

# ACI
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

data "aci_tenant" "ansible_test" {
  name = "ansible_test"
}

data "aci_l4_l7_device" "l4_l7_device_1" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "ansible_test_firewall1"
}

output "aci_l4_l7_device_1" {
  value = data.aci_l4_l7_device.l4_l7_device_1.id
}

data "aci_l4_l7_device" "l4_l7_device_2" {
  tenant_dn = data.aci_tenant.ansible_test.id
  name      = "ansible_test_other"
}

output "aci_l4_l7_device_2" {
  value = data.aci_l4_l7_device.l4_l7_device_2.id
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
  name = "ansible_test"
}

data "mso_site" "tf_site" {
  name = "ansible_test"
}

resource "mso_schema" "schema_test" {
  name = "terraform_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = data.mso_tenant.tf_tenant.id
  }
}

resource "mso_schema_template_service_graph" "test_sg" {
  schema_id          = mso_schema.schema_test.id
  template_name      = one(mso_schema.schema_test.template).name
  service_graph_name = "sgtf1"
  service_node {
    type = "firewall"
  }
  service_node {
    type = "other"
  }
  description = "Terraform Service Graph"
}

resource "mso_schema_site" "schema_site_1" {
  schema_id     = mso_schema.schema_test.id
  site_id       = data.mso_site.tf_site.id
  template_name = mso_schema_template_service_graph.test_sg.template_name
}

resource "mso_schema_site_service_graph" "test_sg_site" {
  schema_id          = mso_schema_site.schema_site_1.schema_id
  site_id            = mso_schema_site.schema_site_1.site_id
  template_name      = mso_schema_template_service_graph.test_sg.template_name
  service_graph_name = mso_schema_template_service_graph.test_sg.service_graph_name
  service_node {
    # for 1st item in the service graph list
    device_dn = data.aci_l4_l7_device.l4_l7_device_1.id
  }
  service_node {
    # for 2nd item in the service graph list
    device_dn = data.aci_l4_l7_device.l4_l7_device_2.id
  }
}

data "mso_schema_site_service_graph" "test_sg_site" {
  schema_id          = mso_schema_site_service_graph.test_sg_site.schema_id
  site_id            = mso_schema_site_service_graph.test_sg_site.site_id
  template_name      = mso_schema_site_service_graph.test_sg_site.template_name
  service_graph_name = mso_schema_site_service_graph.test_sg_site.service_graph_name
}

output "example" {
  value = data.mso_schema_site_service_graph.test_sg_site
}