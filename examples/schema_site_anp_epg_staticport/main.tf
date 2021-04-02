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

resource "mso_schema_site_anp_epg_static_port" "foo_schema_site_anp_epg_static_port" {
  schema_id            = "5c4d5bb72700000401f80948"
  site_id              = "5c7c95b25100008f01c1ee3c"
  template_name        = "Template1"
  anp_name             = "ANP"
  epg_name             = "DB"
  path_type            = "port"
  deployment_immediacy = "lazy"
  pod                  = "pod-4"
  leaf                 = "106"
  path                 = "eth1/11"
  vlan                 = 200
  micro_seg_vlan       = 3
  mode                 = "untagged"
  fex                  = "10"
}

