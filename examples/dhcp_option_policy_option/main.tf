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
  name         = var.tenant_name
  display_name = var.tenant_name
  description  = "DemoTenant"
}

resource "mso_dhcp_option_policy" "example" {
  tenant_id = mso_tenant.tenant1.id
  name = "example"
}

resource "mso_dhcp_option_policy_option" "example"{
    option_policy_name = mso_dhcp_option_policy.example.name
    option_name = "example"
    option_id = "1"
    option_data = "example_data"
}