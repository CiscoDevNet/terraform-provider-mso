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

data "mso_schema" "demo_schema" {
  name = "demo_schema"
}

data "mso_schema_template_contract" "demo_contract" {
  schema_id     = data.mso_schema.demo_schema.id
  template_name = one(data.mso_schema.demo_schema.template).name
  contract_name = "demo_contract"
}

resource "mso_schema_template_filter_entry" "example" {
  schema_id          = data.mso_schema.demo_schema.id
  template_name      = one(data.mso_schema.demo_schema.template).name
  name               = "filter"
  display_name       = "filter"
  entry_name         = "entry"
  entry_display_name = "entry"
}

resource "mso_schema_template_contract_filter" "example" {
  schema_id     = data.mso_schema.demo_schema.id
  template_name = one(data.mso_schema.demo_schema.template).name
  contract_name = data.mso_schema_template_contract.demo_contract.contract_name
  filter_name   = mso_schema_template_filter_entry.example.name
  filter_type   = "bothWay"
  action        = "deny"
  directives    = ["log", "no_stats"]
  priority      = "level1"
}
