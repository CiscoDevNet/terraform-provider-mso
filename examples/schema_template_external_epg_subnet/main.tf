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
  # platform = "nd" # Use it when logging in ND
}

resource "mso_tenant" "test_tenant" {
  name         = "eepg_subnet_tenant"
  display_name = "eepg_subnet_tenant"
  description  = "DemoTenant"
}

resource "mso_schema" "test_schema" {
  name = "eepg_subnet_schema"
  template {
    name         = "eepg_subnet_template"
    display_name = "eepg_subnet_template"
    tenant_id    = mso_tenant.test_tenant.id
  }
}

resource "mso_schema_template_vrf" "test_vrf" {
  schema_id              = mso_schema.test_schema.id
  template               = one(mso_schema.test_schema.template).name
  name                   = "eepg_subnet_vrf"
  display_name           = "eepg_subnet_vrf"
  ip_data_plane_learning = "disabled"
  layer3_multicast       = false
}

resource "mso_schema_template_external_epg" "template_externalepg" {
  schema_id         = mso_schema.test_schema.id
  template_name     = one(mso_schema.test_schema.template).name
  external_epg_name = "eepg"
  display_name      = "eepg"
  vrf_name          = mso_schema_template_vrf.test_vrf.name
  external_epg_type = "on-premise"
}

resource "mso_schema_template_external_epg_subnet" "subnet1" {
  schema_id         = mso_schema.test_schema.id
  template_name     = one(mso_schema.test_schema.template).name
  external_epg_name = mso_schema_template_external_epg.template_externalepg.external_epg_name
  ip                = "10.102.100.0/0"
  scope             = ["shared-rtctrl", "export-rtctrl"]
  aggregate         = ["shared-rtctrl", "export-rtctrl"]
}
