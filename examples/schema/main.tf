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
  name         = "demo_tenant"
  display_name = "demo_tenant"
}

// create resource with four template blocks
resource "mso_schema" "schema_blocks" {
  name = "demo_schema_blocks"
  template {
    name         = "Template1"
    display_name = "TEMP1"
    tenant_id    = data.mso_tenant.demo_tenant.id
    template_type = "aci_multi_site"
  }
  template {
    name         = "Template2"
    display_name = "TEMP2"
    tenant_id    = data.mso_tenant.demo_tenant.id
    template_type = "aci_autonomous"
  }
  template {
    name         = "Template3"
    display_name = "TEMP3"
    tenant_id    = data.mso_tenant.demo_tenant.id
    template_type = "cloud_local"
  }
  template {
    name         = "Template4"
    display_name = "TEMP4"
    tenant_id    = data.mso_tenant.demo_tenant.id
    template_type = "ndfc"
  }
}

// Create ANPs associating them with all templates in mso_schema.schema_blocks
resource "mso_schema_template_anp" "anp_loop" {
  for_each     = { for template in tolist(mso_schema.schema_blocks.template) : template.name => template }
  schema_id    = mso_schema.schema_blocks.id
  template     = each.value.name
  name         = "anp1"
  display_name = "anp1"
}

// Create ANP via index of template in mso_schema.schema_blocks
resource "mso_schema_template_anp" "anp_single" {
  schema_id    = mso_schema.schema_blocks.id
  template     = tolist(mso_schema.schema_blocks.template)[1].name
  name         = "anp2"
  display_name = "anp2"
}