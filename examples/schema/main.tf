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
  // platform = "nd"  // use when logging in ND.
}

resource "mso_site" "test_site" {
  name         = "test_site"
  username     = "" # <APIC username>
  password     = "" # <APIC pwd>
  apic_site_id = "105"
  urls         = "" # <APIC site url>
  location = {
    lat  = 78.946
    long = 95.623
  }
}

resource "mso_tenant" "tenant1" {
  name         = "test_tenant"
  display_name = "test_tenant"
  description  = "DemoTenant"
  site_associations {
    site_id = mso_site.test_site.id
  }
}

// create resource with template_name and tenant_id
resource "mso_schema" "schema1" {
  name          = "demo_schema"
  template_name = "tempu"
  tenant_id     = mso_tenant.tenant1.id
}

// create resource with three template blocks
resource "mso_schema" "schema_blocks" {
  name = "Schema3"
  template {
    name         = "Template1"
    display_name = "TEMP1"
    tenant_id    = mso_tenant.tenant1.id
  }
  template {
    name         = "Template2"
    display_name = "TEMP2"
    tenant_id    = mso_tenant.tenant1.id
  }
  template {
    name         = "Template3"
    display_name = "TEMP3"
    tenant_id    = mso_tenant.tenant1.id
  }
}

// Create ANPs associating them with all templates in mso_schema.schema_blocks
resource "mso_schema_template_anp" "anp1" {
  for_each     = { for template in tolist(mso_schema.schema_blocks.template) : template.name => template }
  schema_id    = mso_schema.schema_blocks.id
  template     = each.value.name
  name         = "anp1"
  display_name = "anp1"
}

// Create ANP via index of template in mso_schema.schema_blocks
resource "mso_schema_template_anp" "anp2" {
  schema_id    = mso_schema.schema_blocks.id
  template     = tolist(mso_schema.schema2.template)[1].name
  name         = "anp2"
  display_name = "anp2"
}