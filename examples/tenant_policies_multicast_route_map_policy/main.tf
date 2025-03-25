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

# tenant policies multicast route map policy example

resource "mso_tenant_policies_multicast_route_map_policy" "multicast_route_map_policy" {
  template_id = mso_template.tenant_template.id
  name        = "multicast_route_map_policy"
  description = "Example description"
  multicast_route_map_entries {
    order     = 1
    group_ip  = "226.2.2.2/8"
    source_ip = "1.1.1.1/1"
    rp_ip     = "1.1.1.2"
    action    = "permit"
  }
  multicast_route_map_entries {
    order     = 2
    group_ip  = "230.3.3.3/32"
    action    = "deny"
  }
}
