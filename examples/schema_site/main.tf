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

resource "mso_site" "site_1" {
  name         = var.site_name_1
  username     = "" # <site username>
  password     = "" # <site pwd>
  apic_site_id = "105"
  urls         = ["https://10.10.100.111"] # <site url>
  location = {
    lat  = 78.946
    long = 95.623
  }
}

resource "mso_site" "site_2" {
  name         = var.site_name_2
  username     = "" # <site username>
  password     = "" # <site pwd>
  apic_site_id = "106"
  urls         = ["https://10.10.10.125"] # <site url>
  location = {
    lat  = 79.946
    long = 96.623
  }
}

resource "mso_tenant" "tenant_1" {
  name         = var.tenant_name
  display_name = var.tenant_name
  description  = "DemoTenant"
  site_associations {
    site_id = mso_site.site_1.id
  }
  site_associations {
    site_id = mso_site.site_2.id
  }
}

resource "mso_schema" "schema_1" {
  name          = var.schema_name
  template_name = var.template_name
  tenant_id     = mso_tenant.tenant_1.id
}

resource "mso_schema_site" "schema_site_1" {
  schema_id     = mso_schema.schema_1.id
  site_id       = mso_site.ansible_test.id
  template_name = var.template_name
}

// when a template should be undeployed from site before disassociation the 'undeploy_on_delete' argument should be set to true prior to terraform destroy  
resource "mso_schema_site" "schema_site_2" {
  schema_id          = mso_schema.schema_1.id
  site_id            = mso_site.site_test_1.id
  template_name      = var.template_name
  undeploy_on_delete = true
}