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

resource "mso_schema_site_anp_epg" "foo_schema_site_anp_epg" {
  schema_id     = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id       = "5c7c95d9510000cf01c1ee3d"
  anp_name      = "ANP"
  epg_name      = "DB"
  private_link_label {
    name = "Cloud"
    }
}
