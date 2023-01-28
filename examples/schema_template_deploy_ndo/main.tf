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

resource "mso_site" "site_test_1" {
  name         = "site1"
  username     = "admin"
  password     = "test"
  apic_site_id = "100"
  urls         = ["https://3.208.123.222"]
}

resource "mso_tenant" "tenant1" {
  name         = "test_tenant"
  display_name = "test_tenant"
  site_associations {
    site_id = mso_site.site_test_1.id
  }
}

resource "mso_schema" "schema1" {
  name = "test_schema"
  template {
    name         = "test_template"
    display_name = "test_template"
    tenant_id    = mso_tenant.tenant1.id
  }
}

resource "mso_schema_template_anp" "anp1" {
  schema_id    = mso_schema.schema1.id
  template     = tolist(mso_schema.schema1.template)[0].name
  name         = "anp1"
  display_name = "anp1"
}

// when a template should be undeployed from a site before disassociation the 'undeploy_on_destroy' argument should be set to true prior to terraform destroy  
resource "mso_schema_site" "schema_site1" {
  schema_id           = mso_schema.schema1.id
  template_name       = tolist(mso_schema.schema1.template)[0].name
  site_id             = mso_site.site_test_1.id
  undeploy_on_destroy = true
}

resource "mso_schema_template_deploy_ndo" "deploy_ndo" {
  schema_id     = mso_schema.schema1.id
  template_name = tolist(mso_schema.schema1.template)[0].name
}

// when a redeploy is preferred the boolean flag of re_deploy should be set to true as shown below
resource "mso_schema_template_deploy_ndo" "redeploy_ndo" {
  schema_id     = mso_schema.schema1.id
  template_name = tolist(mso_schema.schema1.template)[0].name
  re_deploy     = true
}
