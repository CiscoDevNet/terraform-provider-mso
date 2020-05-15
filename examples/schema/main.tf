provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema" "schema1" {
  name          = "nkfi"
  template_name = "tempu"
  tenant_id     = "5eac0d982c00006dae0a28f6"

}