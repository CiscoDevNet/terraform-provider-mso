provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_site_anp_epg_staticleaf" "foo_schema_site_anp_epg_staticleaf" {
  schema_id       = "5c4d9fca270000a101f8094a"
  template_name   = "Template1"
  site_id         = "5c7c95b25100008f01c1ee3c"
  anp_name        = "ANP"
  epg_name        = "Web"
  path            = "topology/pod-1/paths-103/pathep-[eth1/111]"
  port_encap_vlan = 100
}
