provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema" "schema1" {
  name          = "nkp1002"
  template_name = "temp1"
  tenant_id     = "5e9d09482c000068500a269a"

}

resource "mso_schema_site" "schemasite1" {
    schema_id = "${mso_schema.schema1.id}"
    template_name = "temp1"
    site_id = "5c7c95b25100008f01c1ee3c"
}

data "mso_schema" "schema10" {
  name = "nkp1002"
}

output "demo1" {
  value = "${data.mso_schema.schema10}"
}



data "mso_schema_site" "schemasite10" {
  name = "On-premises"
  schema_id = "${mso_schema.schema1.id}"
  
}

output "demo" {
  value = "${data.mso_schema_site.schemasite10}"
}



