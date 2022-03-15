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
	vrf_name = "mso_schema_site_vrf_example"
}

resource "mso_schema_site_vrf_region" "example"{
    schema_id = mso_schema_site.example.schema_id
    template_name = mso_schema_site.example.template_name
    site_id = mso_schema_site.example.site_id
    vrf_name = mso_schema_site_vrf.example.vrf_name
    region_name = "example"
    cidr {
        cidr_ip = "2.2.2.2/10"
        primary = "true"
        subnet {
            ip = "1.20.30.4"
        }
    }
}

resource "mso_schema_site_vrf_region_hub_network" "example"{
    schema_id = mso_schema_site.example.schema_id
    template_name = mso_schema_site.example.template_name
    site_id = mso_schema_site.example.site_id
    vrf_name = mso_schema_site_vrf.example.vrf_name
    region_name = mso_schema_site_vrf_region.example.region_name
    name = "example"
    tenant_name = data.mso_tenant.example.id
}