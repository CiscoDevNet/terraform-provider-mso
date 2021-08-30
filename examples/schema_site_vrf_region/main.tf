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

resource "mso_schema_site_vrf_region" "vrfRegion" {
  schema_id = "5efd6ea60f00005b0ebbd643"
  template_name = "Template1"
  site_id = "5efeb3c4190000cc12d05376"
  vrf_name = "Myvrf"
  region_name = "us-east-1"
  vpn_gateway = true
  hub_network_enable = true
  hub_network = {
    name = "hub-fualt"
    tenant_name = "infra"
  }
  cidr {
    cidr_ip = "2.2.2.2/10"
    primary = true
    subnet {
      ip = "1.20.30.4"
      zone = "us-east-1b"
      usage = "sdfkhsdkf"
    }
  }
}
