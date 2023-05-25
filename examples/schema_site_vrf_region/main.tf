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
  name = "demo_tenant"
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

// AZURE specific example
resource "mso_schema_site_vrf_region" "vrf_region_azure" {
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
    cidr_ip = "2.2.2.2/10"
    primary = true
    subnet {
      ip    = "1.20.30.4"
      name  = "subnet1"
      usage = "gateway"
    }
  }
}

// AWS specific example
resource "mso_schema_site_vrf_region" "vrf_region_aws" {
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
    cidr_ip = "2.2.2.2/10"
    primary = true
    subnet {
      ip   = "1.20.30.4"
      name = "subnet1"
      zone = "us-east-1b"
    }
  }
}

// GCP specific example
resource "mso_schema_site_vrf_region" "vrf_region_gcp" {
  schema_id          = mso_schema.abr_schema.id
  template_name      = one(mso_schema.abr_schema.template).name
  site_id            = mso_schema_site.demo_schema_site.id
  region_name        = "southamerica-east1"
  vrf_name           = mso_schema_template_vrf.vrf1.name
  hub_network_enable = true
  hub_network = {
    name        = "default"
    tenant_name = "infra"
  }
  cidr {
    cidr_ip = "1.1.1.0/24"
    primary = true
    subnet {
      ip           = "1.1.1.0/28"
      name         = "subnet1"
      subnet_group = "test_group"
    }
  }
}