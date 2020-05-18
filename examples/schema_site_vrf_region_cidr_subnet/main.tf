provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_site_vrf_region_cidr_subnet" "foo_schema_site_vrf_region_cidr_subnet" {
  schema_id     = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id       = "5ce2de773700006a008a2678"
  vrf_name      = "Campus"
  region_name   = "westus"
  cidr_ip       = "1.1.1.1/24"
  ip            = "203.168.240.1/24"
  zone          = "West"
  usage         = "gateway"
}