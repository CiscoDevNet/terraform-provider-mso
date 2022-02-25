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
	tenant_id = data.mso_tenant.test.id
	name = "example"		
}

resource mso_schema "example"{
	name = "example"
	template_name = "example"
	tenant_id = data.mso_tenant.test.id
}

resource mso_schema_template_vrf "example" {
	schema_id = mso_schema.test.id
	template = mso_schema.test.template_name
	name = "example"
	display_name= "example"
}

resource "mso_schema_template_external_epg" "example" {
	schema_id = mso_schema.test.id
	template_name = mso_schema.test.template_name
	external_epg_name = "example"
	display_name = "example"
	vrf_name = mso_schema_template_vrf.test.name
}

resource "mso_dhcp_relay_policy_provider" "example" {
	dhcp_relay_policy_name = mso_dhcp_relay_policy.test.name
	dhcp_server_address = "example"
	external_epg_ref = mso_schema_template_external_epg.test.id
}