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
	template_name = mso_schema.example.template_name
}

resource "mso_schema_site_anp" "example" {
  schema_id     = mso_schema.example.id
  anp_name      = "ANP_EXAMPLE"
  template_name = mso_schema_site.example.template_name
  site_id       = data.mso_site.example.id
}

resource "mso_schema_site_anp_epg" "example" {
  schema_id     = mso_schema.example.id
  template_name = mso_schema_site.example.template_name
  site_id       = data.mso_site.example.id
  anp_name      = mso_schema_site_anp.example.anp_name
  epg_name      = "EPG_EXAMPLE"
}

resource "mso_schema_site_anp_epg_useg_attr" "useg_attrs" {
  schema_id     = mso_schema.example.id
  site_id       = data.mso_site.example.id
  anp_name      = mso_schema_site_anp.example.anp_name
  epg_name      = mso_schema_site_anp_epg.example.epg_name
  template_name = mso_schema_site.example.template_name
  useg_name     = "useg_site_test"
  useg_type     = "tag"
  operator      = "startsWith"
  category      = "tagger"
  value         = "10.2.3.4"
}
