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

resource "mso_schema_site_anp_epg_subnet" "foo_schema_site_anp_epg_subnet" {
  schema_id          = "5c4d5bb72700000401f80948"
  site_id            = "5c7c95b25100008f01c1ee3c"
  template_name      = "Template1"
  anp_name           = "ANP"
  epg_name           = "DB"
  ip                 = "10.0.1.1/8"
  scope              = "private"
  shared             = false
  description        = "This is schema site anp epg subnet."
  no_default_gateway = false
  querier            = false
}

