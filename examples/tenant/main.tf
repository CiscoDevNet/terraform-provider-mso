provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_tenant" "tenant01" {
  name              = "m22"
  display_name      = "m22"
  description       = "DemoTenant"
  site_associations{
      site_id = "5c7c95b25100008f01c1ee3c"
      }
  user_associations{
      user_id = "0000ffff0000000000000020"
    }
}