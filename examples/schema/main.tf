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

// create resource with template_name and tenant_id
resource "mso_schema" "schema1" {
  name          = "demo_schema"
  template_name = "tempu"
  tenant_id     = "0000ffff0000000000000010"
}

// create resource with three template blocks
resource "mso_schema" "schema_blocks" {
  name            = "Schema3"
  template {
    name          = "Template1"
    display_name  = "TEMP1"
    tenant_id     = "623316531d0000abdd50343a"
  }
  template {
    name          = "Template2"
    display_name  = "TEMP2"
    tenant_id     = "623316531d0000abdd50343a"
  }
  template {
    name          = "Template3"
    display_name  = "TEMP3"
    tenant_id     = "0000ffff0000000000000010"
  }
}  

// Create ANPs associating them with all templates in mso_schema.schema_blocks
resource "mso_schema_template_anp" "anp1" {
  for_each  = { for template in tolist(mso_schema.schema_blocks.template) : template.name => template }
  schema_id     = mso_schema.schema_blocks.id
  template      = each.value.name
  name          = "anp1"
  display_name  = "anp1"
}

// Create ANP via index of template in mso_schema.schema_blocks
resource "mso_schema_template_anp" "anp2" {
  schema_id     = mso_schema.schema_blocks.id
  template = tolist(mso_schema.schema2.template)[1].name
  name          = "anp2"
  display_name  = "anp2"
}