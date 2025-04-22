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

data "mso_site" "ansible_test" {
  name = "ansible_test"
}

resource "mso_tenant" "tf_test_tenant" {
  name         = "tf_test_tenant"
  display_name = "tf_test_tenant"
  site_associations {
    site_id = data.mso_site.ansible_test.id
  }
}

resource "mso_schema" "tf_test_mso_schema" {
  name = "tf_test_mso_schema"
  template {
    name         = "tf_test_mso_schema_template"
    display_name = "tf_test_mso_schema_template"
    tenant_id    = mso_tenant.tf_test_tenant.id
  }
}

resource "mso_schema_template_anp" "tf_test_anp" {
  name         = "tf_test_anp"
  display_name = "tf_test_anp"
  schema_id    = mso_schema.tf_test_mso_schema.id
  template     = "tf_test_mso_schema_template"
}

resource "mso_schema_template_anp_epg" "tf_test_anp_epg" {
  name          = "tf_test_anp_epg"
  display_name  = "tf_test_anp_epg"
  anp_name      = mso_schema_template_anp.tf_test_anp.name
  schema_id     = mso_schema.tf_test_mso_schema.id
  template_name = mso_schema_template_anp.tf_test_anp.template
}

resource "mso_schema_template_vrf" "tf_test_vrf" {
  name         = "tf_test_vrf"
  display_name = "tf_test_vrf"
  schema_id    = mso_schema.tf_test_mso_schema.id
  template     = "tf_test_mso_schema_template"
}

resource "mso_schema_template_external_epg" "tf_test_ext_epg" {
  external_epg_name = "tf_test_ext_epg"
  display_name      = "tf_test_ext_epg"
  vrf_name          = mso_schema_template_vrf.tf_test_vrf.name
  schema_id         = mso_schema.tf_test_mso_schema.id
  template_name     = "tf_test_mso_schema_template"
}
resource "mso_template" "tf_test_tenant_policy_template" {
  template_name = "tf_test_tenant_policy_template"
  template_type = "tenant"
  tenant_id     = mso_tenant.tf_test_tenant.id
}

resource "mso_tenant_policies_dhcp_relay_policy" "tf_test_dhcp_relay_policy" {
  name        = "tf_test_dhcp_relay_policy"
  template_id = mso_template.tf_test_tenant_policy_template.id
  providers {
    dhcp_server_address  = "1.1.1.1"
    application_epg_uuid = mso_schema_template_anp_epg.tf_test_anp_epg.uuid
  }
  providers {
    dhcp_server_address        = "2.2.2.2"
    external_epg_uuid          = mso_schema_template_external_epg.tf_test_ext_epg.uuid
    dhcp_server_vrf_preference = true
  }
}
