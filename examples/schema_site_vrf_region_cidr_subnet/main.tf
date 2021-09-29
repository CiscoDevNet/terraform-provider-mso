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

resource "mso_schema_site_vrf_region_cidr_subnet" "azure_schema_site_vrf_region_cidr_subnet" {
  schema_id     = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id       = "5ce2de773700006a008a2678"
  vrf_name      = "Campus"
  region_name   = "westus"
  cidr_ip       = "203.168.0.0/16"
  ip            = "203.168.240.0/24"
  usage         = "gateway"
}

resource "mso_schema_site_vrf_region_cidr_subnet" "aws_schema_site_vrf_region_cidr_subnet" {
  schema_id     = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id       = "5ce2de773700006a008a2679"
  vrf_name      = "Campus"
  region_name   = "us-east"
  cidr_ip       = "203.167.0.0/16"
  ip            = "203.167.240.0/24"
  zone          = "us-east-1b"
}