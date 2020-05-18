provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template_contract_filter" "contractfilter01" {
  schema_id             = "5c4d5bb72700000401f80948"
  template_name         = "Template1"
  contract_name         = "Web-to-DB"
  filter_type           = "provider_to_consumer"
  filter_name           = "Any100"
  filter_schema_id      = "5c4d5bb72700000401f80948"
  filter_template_name  = "Template1"
  directives            = ["none","log"]
}
