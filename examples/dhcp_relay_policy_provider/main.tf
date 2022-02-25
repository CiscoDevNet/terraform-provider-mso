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

resource "mso_dhcp_relay_policy" "example" {
	tenant_id = data.mso_tenant.example.id
	name = "example"		
}

resource mso_schema "example"{
	name = "example"
	template_name = "example"
	tenant_id = data.mso_tenant.example.id
}

resource mso_schema_template_vrf "example" {
	schema_id = mso_schema.example.id
	template = mso_schema.example.template_name
	name = "example"
	display_name= "example"
}

resource "mso_schema_template_external_epg" "example" {
	schema_id = mso_schema.example.id
	template_name = mso_schema.example.template_name
	external_epg_name = "example"
	display_name = "example"
	vrf_name = mso_schema_template_vrf.example.name
}

resource "mso_dhcp_relay_policy_provider" "example" {
	dhcp_relay_policy_name = mso_dhcp_relay_policy.example.name
	dhcp_server_address = "1.2.3.4"
	external_epg_ref = mso_schema_template_external_epg.example.id
}