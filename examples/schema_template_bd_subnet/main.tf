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

resource "mso_schema_template_vrf" "demo_vrf" {
  schema_id    = mso_schema.demo_schema.id
  template     = one(mso_schema.demo_schema.template).name
  name         = "demo_vrf"
  display_name = "demo_vrf"
}

resource "mso_schema_template_bd" "demo_bd" {
  schema_id     = mso_schema.demo_schema.id
  template_name = one(mso_schema.demo_schema.template).name
  name          = "demo_bd"
  display_name  = "demo_bd"
  vrf_name      = mso_schema_template_vrf.demo_vrf.name
}

resource "mso_schema_template_bd_subnet" "bdsubnet01" {
  schema_id              = mso_schema.demo_schema.id
  template_name          = one(mso_schema.demo_schema.template).name
  bd_name                = mso_schema_template_bd.demo_bd.name
  ip                     = "10.26.17.1/24"
  scope                  = "public"
  ip_data_plane_learning = "enabled"
  description            = "SubnetDemo"
  shared                 = false
  no_default_gateway     = false
  querier                = false
}
