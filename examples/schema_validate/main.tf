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

data "mso_tenant" "example" {
	name = "example"
	display_name = "example"
}

resource "mso_schema" "example" {
	name = "example"
	template_name = "example"
	tenant_id = data.mso_tenant.test.id
}

data "mso_schema_validate" "example" {
  schema_id = mso_schema.example.id
}