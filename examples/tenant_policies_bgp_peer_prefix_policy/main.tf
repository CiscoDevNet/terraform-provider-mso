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

data "mso_tenant" "example_tenant" {
  name = "example_tenant"
}

# tenant template example

resource "mso_template" "tenant_template" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id = data.mso_tenant.example_tenant.id
}

# tenant policies bgp peer prefix policy example

resource "mso_tenant_policies_bgp_peer_prefix_policy" "bgp_policy" {
  template_id             = mso_template.tenant_template.id
  name                    = "test_bgp_peer_prefix_policy"
  description             = "Test BGP Peer Prefix Policy"
  action                  = "restart"
  max_number_of_prefixes  = 1000
  threshold_percentage    = 50
  restart_time            = 60
}
