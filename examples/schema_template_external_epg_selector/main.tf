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

resource "mso_schema_template_external_epg" "template_externalepg" {
  schema_id         = mso_schema.test_schema.id
  template_name     = one(mso_schema.test_schema.template).name
  external_epg_name = "eepg"
  display_name      = "eepg"
  vrf_name          = mso_schema_template_vrf.test_vrf.name
}

resource "mso_schema_template_external_epg_selector" "selector1" {
  schema_id         = mso_schema_template_external_epg.template_externalepg.schema_id
  template_name     = mso_schema_template_external_epg.template_externalepg.template_name
  external_epg_name = mso_schema_template_external_epg.template_externalepg.external_epg_name
  name              = "eepg_selector"
  expressions {
    value = "1.20.56.44"
  }
  expressions {
    value = "5.6.7.8"
  }
}
