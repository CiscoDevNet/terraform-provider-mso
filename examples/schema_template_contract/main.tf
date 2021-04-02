terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template_contract" "template_contract" {
  schema_id             = "5c4d5bb72700000401f80948"
  template_name         = "Template1"
  contract_name         = "C2"
  display_name          = "C2"
  filter_type           = "bothWay"
  scope                 = "context"
  filter_relationships  ={
    filter_schema_id    = "5c4d5bb72700000401f80948"
    filter_template_name = "Template1"
    filter_name = "filter1"
  }
  directives            = ["none"]
}