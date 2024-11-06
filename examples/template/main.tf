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

data "mso_site" "site_1" {
  name = "example_site_1"
}

data "mso_site" "site_2" {
  name = "example_site_2"
}

data "mso_tenant" "example_tenant" {
  name = "example_tenant"
}

# tenant template example

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id = data.mso_tenant.example_tenant.id
  sites = [data.mso_site.site_1.id, data.mso_site.site_2.id]
}

# l3out template example

resource "mso_template" "l3out_template" {
  template_name = "l3out_template"
  template_type = "l3out"
  tenant_id = data.mso_tenant.example_tenant.id
  sites = [data.mso_site.site_1.id]
}

# fabric policy template example

resource "mso_template" "fabric_policy_template" {
  template_name = "fabric_policy_template"
  template_type = "fabric_policy"
  sites = [data.mso_site.site_1.id, data.mso_site.site_2.id]
}

# fabric resource template example

resource "mso_template" "fabric_resource_template" {
  template_name = "fabric_resource_template"
  template_type = "fabric_resource"
  sites = [data.mso_site.site_1.id, data.mso_site.site_2.id]
}

# monitoring tenant template example

resource "mso_template" "monitoring_tenant_template" {
  template_name = "monitoring_tenant_template"
  template_type = "monitoring_tenant"
  tenant_id = data.mso_tenant.example_tenant.id
  sites = [data.mso_site.site_1.id]
}

# monitoring access template example

resource "mso_template" "monitoring_access_template" {
  template_name = "monitoring_access_template"
  template_type = "monitoring_access"
  sites = [data.mso_site.site_1.id]
}

# service device template example

resource "mso_template" "service_device_template" {
  template_name = "service_device_template"
  template_type = "service_device"
  tenant_id = data.mso_tenant.example_tenant.id
  sites = [data.mso_site.site_1.id, data.mso_site.site_2.id]
}
