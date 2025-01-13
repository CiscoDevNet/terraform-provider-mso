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
  name = "test_bd_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = mso_tenant.tenant_1.id
  }
}

resource "mso_schema_template_vrf" "vrf" {
  schema_id    = mso_schema.schema_1.id
  template     = one(mso_schema.schema_1.template).name
  name         = "test_bd_vrf"
  display_name = "test_bd_vrf"
}

// Endpoint move detect mode

resource "mso_schema_template_bd" "bd_ep" {
  schema_id              = mso_schema.schema_1.id
  template_name          = one(mso_schema.schema_1.template).name
  vrf_name               = mso_schema_template_vrf.vrf.name
  name                   = "bd_ep_demo"
  display_name           = "bd_ep_demo"
  arp_flooding           = true
  ep_move_detection_mode = "garp"
}


// MSO versions 3.2 and higher
resource "mso_schema_template_bd" "bd" {
  schema_id                       = mso_schema.schema_1.id
  template_name                   = one(mso_schema.schema_1.template).name
  name                            = "bd_demo"
  display_name                    = "bd_demo"
  vrf_name                        = mso_schema_template_vrf.vrf.name
  vrf_schema_id                   = mso_schema_template_vrf.vrf.schema_id
  vrf_template_name               = mso_schema_template_vrf.vrf.template
  layer2_unknown_unicast          = "proxy"
  intersite_bum_traffic           = false
  optimize_wan_bandwidth          = true
  layer2_stretch                  = false
  layer3_multicast                = false
  multi_destination_flooding      = "flood_in_encap"
  ipv6_unknown_multicast_flooding = "optimized_flooding"
  unknown_multicast_flooding      = "optimized_flooding"
  dhcp_policies {
    name                       = "Policy1"
    version                    = 10
    dhcp_option_policy_name    = "Policy10"
    dhcp_option_policy_version = 12
  }
  dhcp_policies {
    name                       = "Policy1"
    version                    = 20
    dhcp_option_policy_name    = "Policy20"
    dhcp_option_policy_version = 22
  }
  description = "bd_description"
}

// MSO versions below 3.2

resource "mso_schema_template_bd" "bd1" {
  schema_id                       = mso_schema.schema_1.id
  template_name                   = one(mso_schema.schema_1.template).name
  name                            = "bd_demo1"
  display_name                    = "bd_demo1"
  vrf_name                        = mso_schema_template_vrf.vrf.name
  vrf_schema_id                   = mso_schema_template_vrf.vrf.schema_id
  vrf_template_name               = mso_schema_template_vrf.vrf.template
  layer2_unknown_unicast          = "proxy"
  intersite_bum_traffic           = false
  optimize_wan_bandwidth          = true
  layer2_stretch                  = false
  layer3_multicast                = false
  multi_destination_flooding      = "flood_in_encap"
  ipv6_unknown_multicast_flooding = "optimized_flooding"
  unknown_multicast_flooding      = "optimized_flooding"
  dhcp_policy = {
    name                       = "Policy1"
    version                    = 10
    dhcp_option_policy_name    = "Policy10"
    dhcp_option_policy_version = 12
  }
  description = "bd1_description"
}
