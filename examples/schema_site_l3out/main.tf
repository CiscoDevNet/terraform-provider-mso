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

data "mso_site" "example" {
	name  = "example"
}
	  
data "mso_tenant" "example" {
	name = "example"
	display_name = "example"
}
	  
resource "mso_schema" "example" {
	name = "example"
	template_name = "example"
	tenant_id = data.mso_tenant.example.id
}
			
resource "mso_schema_site" "example" {
	schema_id =  mso_schema.example.id
	site_id = data.mso_site.example.id
	template_name = "example"
}

resource "mso_schema_site_vrf" "example" {
	template_name = mso_schema_site.example.template_name
	site_id = mso_schema_site.example.site_id
	schema_id = mso_schema_site.example.schema_id
	vrf_name = "example"
}

resource "mso_schema_site_l3out" "example" {
    schema_id = mso_schema_site.example.schema_id
	l3out_name = "example"
	template_name = mso_schema_site.example.template_name
	vrf_name = mso_schema_site_vrf.example.vrf_name
	site_id = mso_schema_site.example.site_id
}