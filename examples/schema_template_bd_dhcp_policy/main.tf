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

resource "mso_tenant" "tenant1" {
  name         = "test_tenant"
  display_name = "test_tenant"
}

resource "mso_schema" "schema1" {
  name          = "test_schema"
  template_name = "Template1"
  tenant_id     = mso_tenant.tenant1.id
}

resource "mso_schema_template_vrf" "demovrf01" {
  schema_id        = mso_schema.schema1.id
  template         = mso_schema.schema1.template_name
  name             = "test_vrf"
  display_name     = "test_vrf"
  layer3_multicast = false
}

resource "mso_schema_template_bd" "bridge_domain" {
  schema_id              = mso_schema.schema1.id
  template_name          = mso_schema.schema1.template_name
  name                   = "test_bd"
  display_name           = "test_bd"
  vrf_name               = mso_schema_template_vrf.demovrf01.name
  layer2_unknown_unicast = "proxy"
}

resource "mso_dhcp_relay_policy" "example" {
  tenant_id   = mso_tenant.tenant1.id
  name        = "dhcpRelayPol"
  description = "from Terraform"
}

resource "mso_dhcp_option_policy" "example" {
  tenant_id   = mso_tenant.tenant1.id
  name        = "dhcpOptionPol"
  description = "from Terraform"
}

resource "mso_schema_template_bd_dhcp_policy" "exp" {
  schema_id           = mso_schema.schema1.id
  template_name       = mso_schema.schema1.template_name
  bd_name             = mso_schema_template_bd.bridge_domain.name
  name                = mso_dhcp_relay_policy.example.name
  version             = 1
  dhcp_option_name    = mso_dhcp_option_policy.example.name
  dhcp_option_version = 1
}

data "mso_schema_template_bd_dhcp_policy" "exp2" {
  schema_id           = mso_schema.schema1.id
  template_name       = mso_schema.schema1.template_name
  bd_name             = mso_schema_template_bd.bridge_domain.name
  name                = mso_schema_template_bd_dhcp_policy.exp.name
}