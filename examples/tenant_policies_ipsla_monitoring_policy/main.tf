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

# tenant policies ipsla monitoring policy example

resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla_policy" {
  template_id                  = mso_template.tenant_template.id
  name                         = "ipsla_policy"
  description                  = "Example description"
  sla_type                     = "http"
  http_version                 = "HTTP11"
  http_uri                     = "/example1"
  sla_frequency                = 120
  detect_multiplier            = 4
  request_data_size            = 64
  type_of_service              = 18
  operation_timeout            = 100
  threshold                    = 100
  ipv6_traffic_class           = 255
}
