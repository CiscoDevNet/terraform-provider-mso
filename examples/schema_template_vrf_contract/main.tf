terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "terraform_github_ci"
  password = "Crest@123456"
  url      = "https://173.36.219.66/"
  insecure = true
}

resource "mso_schema_template_vrf_contract" "demovrf01" {
  schema_id              = "5eff091b0e00008318cff859"
  template_name          = "Template1"
  vrf_name               = "myVrf"
  relationship_type      = "provider"
  contract_name          = "hubcon"
  contract_schema_id     = "5efd6ea60f00005b0ebbd643"
  contract_template_name = "Template1"
}

data "mso_schema_template_vrf_contract" "demovrf01" {
  schema_id         = "5eff091b0e00008318cff859"
  template_name     = "Template1"
  vrf_name          = "myVrf"
  relationship_type = "provider"
  contract_name     = "hubcon"
}

