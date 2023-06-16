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

data "mso_tenant" "demo_tenant" {
  name = "demo_tenant"
}

data "mso_schema" "demo_schema" {
  name = "demo_schema"
}

resource "mso_schema_template" "template1" {
  schema_id     = data.mso_schema.demo_schema.id
  name          = "Template1"
  display_name  = "Template1"
  tenant_id     = data.mso_tenant.demo_tenant.id
  template_type = "aci_multi_site"
}
