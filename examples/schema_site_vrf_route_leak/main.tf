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

data "mso_site" "demo_site" {
  name = "demo_site"
}

resource "mso_tenant" "demo_tenant" {
  name         = "demo_tenant"
  display_name = "demo_tenant"
}

resource "mso_schema" "demo_schema" {
  name = "demo_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = mso_tenant.demo_tenant.id
  }
}

resource "mso_schema_site" "demo_schema_site" {
  schema_id     = mso_schema.demo_schema.id
  template_name = one(mso_schema.demo_schema.template).name
  site_id       = data.mso_site.demo_site.id
}

resource "mso_schema_template_vrf" "vrf1" {
  schema_id              = mso_schema.demo_schema.id
  template               = one(mso_schema.demo_schema.template).name
  name                   = "vrf1"
  display_name           = "vrf1"
  ip_data_plane_learning = "enabled"
}

resource "mso_schema_site_vrf_region" "vrf_region_azure1" {
  schema_id          = mso_schema.demo_schema.id
  template_name      = one(mso_schema.demo_schema.template).name
  site_id            = mso_schema_site.demo_schema_site.id
  vrf_name           = mso_schema_template_vrf.vrf1.name
  region_name        = "us-east-1"
  vpn_gateway        = true
  hub_network_enable = true
  hub_network = {
    name        = "default"
    tenant_name = "infra"
  }
  cidr {
    cidr_ip = "1.0.0.0/16"
    primary = true
    subnet {
      ip    = "1.0.0.0/24"
      name  = "subnet1"
      usage = "gateway"
    }
  }
}

resource "mso_schema_template_vrf" "vrf2" {
  schema_id              = mso_schema.demo_schema.id
  template               = one(mso_schema.demo_schema.template).name
  name                   = "vrf2"
  display_name           = "vrf2"
  ip_data_plane_learning = "enabled"
}

resource "mso_schema_site_vrf_region" "vrf_region_azure2" {
  schema_id          = mso_schema.demo_schema.id
  template_name      = one(mso_schema.demo_schema.template).name
  site_id            = mso_schema_site.demo_schema_site.id
  vrf_name           = mso_schema_template_vrf.vrf2.name
  region_name        = "us-east-1"
  vpn_gateway        = true
  hub_network_enable = true
  hub_network = {
    name        = "default"
    tenant_name = "infra"
  }
  cidr {
    cidr_ip = "2.0.0.0/16"
    primary = true
    subnet {
      ip    = "2.0.0.0/24"
      name  = "subnet1"
      usage = "gateway"
    }
  }
}

# LEAK_ALL specific example
resource "mso_schema_site_vrf_route_leak" "vrf1" {
  schema_id       = mso_schema.demo_schema.id
  template_name   = one(mso_schema.demo_schema.template).name
  site_id         = mso_schema_site.demo_schema_site.id
  vrf_name        = mso_schema_site_vrf_region.vrf_region_azure1.vrf_name
  target_vrf_name = mso_schema_site_vrf_region.vrf_region_azure2.vrf_name
  tenant_name     = mso_tenant.demo_tenant.name
}

# Subnet IP all example
resource "mso_schema_site_vrf_route_leak" "vrf2" {
  schema_id       = mso_schema.demo_schema.id
  template_name   = one(mso_schema.demo_schema.template).name
  site_id         = mso_schema_site.demo_schema_site.id
  vrf_name        = mso_schema_site_vrf_region.vrf_region_azure2.vrf_name
  target_vrf_name = mso_schema_site_vrf_region.vrf_region_azure1.vrf_name
  tenant_name     = mso_tenant.demo_tenant.name
  type            = "subnet_ip"
  subnet_ips      = ["2.0.0.1/32", "2.0.0.2/32"]
}
