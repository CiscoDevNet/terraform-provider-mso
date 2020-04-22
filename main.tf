provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema" "schema1" {
  name          = "nkp"
  template_name = "temp1"
  tenant_id     = "5e9d09482c000068500a269a"

}

resource "mso_schema_site" "schemasite1" {
    schema_id = mso_schema.schema1.id
    template_name = "temp1"
    site_id = "5c7c95b25100008f01c1ee3c"
}


