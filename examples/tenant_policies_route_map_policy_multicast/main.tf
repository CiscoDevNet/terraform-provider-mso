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
  tenant_id     = data.mso_tenant.example_tenant.id
}

# tenant policies route map policy for multicast example

resource "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
  template_id = mso_template.tenant_template.id
  name        = "route_map_policy_multicast"
  description = "Example description"
  route_map_multicast_entries {
    order               = 1
    group_ip            = "226.2.2.2/8"
    source_ip           = "1.1.1.1/1"
    rendezvous_point_ip = "1.1.1.2"
    action              = "permit"
  }
  route_map_multicast_entries {
    order    = 2
    group_ip = "230.3.3.3/32"
    action   = "deny"
  }
}
