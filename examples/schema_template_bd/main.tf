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
  # platform = "nd"
  insecure = true
}

resource "mso_tenant" "tenant_1" {
  name         = "test_bd_tenant"
  display_name = "test_bd_tenant"
  description  = "DemoTenant"
}

resource "mso_schema" "schema_1" {
  name          = "test_bd_schema"
  template_name = "test_bd"
  tenant_id     = mso_tenant.tenant_1.id
}

resource "mso_schema_template_vrf" "vrf" {
  schema_id     = mso_schema.schema_1.id
  template      = mso_schema.schema_1.template_name
  name          = "test_bd_vrf"
  display_name = "test_bd_vrf"
}

// MSO versions 3.2 and higher

resource "mso_schema_template_bd" "bd" {
  schema_id              = mso_schema.schema_1.id
  template_name          = mso_schema.schema_1.template_name
  name                   = "bd_demo"
  display_name           = "bd_demo"
  vrf_name               = mso_schema_template_vrf.vrf.name
  vrf_schema_id          = mso_schema_template_vrf.vrf.schema_id
  vrf_template_name      = mso_schema_template_vrf.vrf.template
  layer2_unknown_unicast = "proxy"
  intersite_bum_traffic  = false
  optimize_wan_bandwidth = true
  layer2_stretch         = true
  layer3_multicast       = false
  multi_destination_flooding = "flood_in_encap"
  ipv6_unknown_multicast_flooding = "optimized_flooding"
  unknown_multicast_flooding = "optimized_flooding"
  dhcp_policies {
      name = "Policy1"
      version = 10
      dhcp_option_policy_name = "Policy10"
      dhcp_option_policy_version = 12
  }
}

// MSO versions below 3.2

resource "mso_schema_template_bd" "bd" {
  schema_id              = mso_schema.schema_1.id
  template_name          = mso_schema.schema_1.template_name
  name                   = "bd_demo"
  display_name           = "bd_demo"
  vrf_name               = mso_schema_template_vrf.vrf.name
  vrf_schema_id          = mso_schema_template_vrf.vrf.schema_id
  vrf_template_name      = mso_schema_template_vrf.vrf.template
  layer2_unknown_unicast = "proxy"
  intersite_bum_traffic  = false
  optimize_wan_bandwidth = true
  layer2_stretch         = true
  layer3_multicast       = false
  multi_destination_flooding = "flood_in_encap"
  ipv6_unknown_multicast_flooding = "optimized_flooding"
  unknown_multicast_flooding = "optimized_flooding"
  dhcp_policy = {
      name = "Policy1"
      version = 10
      dhcp_option_policy_name = "Policy10"
      dhcp_option_policy_version = 12
  }
}
