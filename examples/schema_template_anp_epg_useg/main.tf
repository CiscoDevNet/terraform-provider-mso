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

resource "mso_schema_template_anp_epg_useg_attr" "useg_attrs" {
  schema_id     = "5eafca7d2c000052860a2902"
  anp_name      = "sanp1"
  epg_name      = "nkuseg"
  template_name = "stemplate1"
  name          = "usg_test"
  useg_type     = "tag"
  operator      = "startsWith"
  category      = "tagger"
  value         = "10.2.3.4"
  useg_subnet   = true

}

# data "mso_schema_template_anp_epg_useg_attr" "useg_attrs" {
#   schema_id     = "5eafca7d2c000052860a2902"
#   anp_name      = "sanp1"
#   epg_name      = "nkuseg"
#   template_name = "stemplate1"
#   name          = "usg_test"
# }
