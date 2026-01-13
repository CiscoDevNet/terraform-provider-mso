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

# tenant policies dhcp option policy example

resource "mso_tenant_policies_dhcp_option_policy" "dhcp_policy" {
  template_id = mso_template.tenant_template.id
  name        = "test_dhcp_option_policy"
  description = "Test DHCP Option Policy"
    
  options {
    name = "option_1"
    id   = 1
    data = "data_1"
  }
    
  options {
    name = "option_2"
    id   = 2
    data = "data_2"
  }
}
