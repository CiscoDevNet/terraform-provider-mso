
provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}
# data "mso_schema_template_anp_epg_contract" "contract" {
#  schema_id = "5eb009212c00003d960a2906"
#  template_name = "stemplate"
#  anp_name = "sanp1"
#  epg_name = "anpepg1"
#  contract_name = "sC14"
# }

resource "mso_schema_template_anp_epg_contract" "contract2" {
   count = 5
 schema_id = "5eb009212c00003d960a2906"
 template_name = "stemplate"
 anp_name = "sanp1"
 epg_name = "anpepg1"
 contract_name = "sC1${count.index}"
 relationship_type = "provider"
}