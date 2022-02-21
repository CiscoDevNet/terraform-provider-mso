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

resource "mso_dhcp_option_policy" "dp1" {
  tenant_id = mso_tenant.tenant1.id
  name = "dhcp_opx"
  description = "desc"
  option {
    name = "op2"
    data = "d1"
    id = "2"
  }
  option {
    name = "op3"
    data = "d3"
    id = "4"
  }
  option {
    name = "op1"
    data = "d1"
    id = "1"
  }
}
