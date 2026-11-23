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
  # platform = "nd"
  insecure = true
}

# DHCP policies live in a tenant policy template, separately from the schema
# template containing the Bridge Domain.
data "mso_tenant" "example" {
  name = "example_tenant"
}

resource "mso_schema" "example" {
  name = "example_schema"

  template {
    name         = "example_template"
    display_name = "example_template"
    tenant_id    = data.mso_tenant.example.id
  }
}

resource "mso_template" "dhcp_policies" {
  template_name = "example_dhcp_policies"
  template_type = "tenant"
  tenant_id     = data.mso_tenant.example.id
}

resource "mso_schema_template_vrf" "example" {
  schema_id    = mso_schema.example.id
  template     = one(mso_schema.example.template).name
  name         = "example_vrf"
  display_name = "example_vrf"
}

resource "mso_schema_template_bd" "example" {
  schema_id     = mso_schema.example.id
  template_name = one(mso_schema.example.template).name
  name          = "example_bd"
  display_name  = "example_bd"
  vrf_name      = mso_schema_template_vrf.example.name
}

resource "mso_schema_template_anp" "example" {
  schema_id    = mso_schema.example.id
  template     = one(mso_schema.example.template).name
  name         = "example_anp"
  display_name = "example_anp"
}

resource "mso_schema_template_anp_epg" "example" {
  schema_id     = mso_schema.example.id
  template_name = one(mso_schema.example.template).name
  anp_name      = mso_schema_template_anp.example.name
  name          = "example_epg"
  display_name  = "example_epg"
  bd_name       = mso_schema_template_bd.example.name
  vrf_name      = mso_schema_template_vrf.example.name
}

resource "mso_tenant_policies_dhcp_relay_policy" "example" {
  template_id = mso_template.dhcp_policies.id
  name        = "example_dhcp_relay_policy"

  dhcp_relay_providers {
    dhcp_server_address  = "192.0.2.10"
    application_epg_uuid = mso_schema_template_anp_epg.example.uuid
  }
}

resource "mso_tenant_policies_dhcp_option_policy" "example" {
  template_id = mso_template.dhcp_policies.id
  name        = "example_dhcp_option_policy"

  options {
    name = "routers"
    id   = 3
    data = "192.0.2.1"
  }
}

# Keep this BD's DHCP associations out of its inline dhcp_policies block when
# using the standalone resource.
resource "mso_schema_template_bd_dhcp_policy" "example" {
  schema_id               = mso_schema_template_bd.example.schema_id
  template_name           = mso_schema_template_bd.example.template_name
  bd_name                 = mso_schema_template_bd.example.name
  name                    = mso_tenant_policies_dhcp_relay_policy.example.name
  dhcp_option_policy_name = mso_tenant_policies_dhcp_option_policy.example.name
}

data "mso_schema_template_bd_dhcp_policy" "example" {
  schema_id     = mso_schema_template_bd_dhcp_policy.example.schema_id
  template_name = mso_schema_template_bd_dhcp_policy.example.template_name
  bd_name       = mso_schema_template_bd_dhcp_policy.example.bd_name
  name          = mso_schema_template_bd_dhcp_policy.example.name
}

output "dhcp_option_policy_name" {
  description = "DHCP Option Policy associated with the name-configured DHCP Relay Policy."
  value       = data.mso_schema_template_bd_dhcp_policy.example.dhcp_option_policy_name
}
