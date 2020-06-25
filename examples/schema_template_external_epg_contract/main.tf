provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}


resource "mso_schema_template_external_epg_contract" "c1" {
  schema_id                 = "5ea809672c00003bc40a2799"
  template_name             = "Template1"
  contract_name             = "contractdemo"
  external_epg_name         = "UntitledExternalEPG1"
  relationship_type         = "consumer"
  contract_schema_id        = "5ea809672c00003bc40a2799"
  contract_template_name    = "Template1"
}