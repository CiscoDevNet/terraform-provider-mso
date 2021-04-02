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

resource "mso_schema_site_bd_subnet" "foo_schema_site_bd_subnet" {
  schema_id          = "5d5dbf3f2e0000580553ccce"
  template_name      = "Template1"
  site_id            = "5c7c95b25100008f01c1ee3c"
  bd_name            = "WebServer-Finance"
  ip                 = "200.168.240.1/24"
  description        = "This is schema site bd subnet."
  shared             = false
  scope              = "private"
  querier            = false
  no_default_gateway = false
}