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


resource "mso_schema_template_anp_epg_contract" "contract01" {
  schema_id              = "5e2dd7112c00005db60a268b"
  template_name          = "Template1"
  anp_name               = "ANP-Financial"
  epg_name               = "Web"
  contract_name          = "Web-to-Internet-Financial"
  relationship_type      = "provider"
  contract_schema_id     = "5e2dd7112c00005db60a268b"
  contract_template_name = "Template1"
}
