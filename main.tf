provider "mso" {
  username = "admin"
  password = "Ins3965!Ins3965!"
  url      = "https://10.23.248.134"
  insecure = true
}

resource "mso_schema" "schema1" {
  name          = "nkp"
  template_name = "temp1"
  tenant_id     = "5e9fcbc30f00000d06fccecf"

}

# resource "mso_schema_site" "schemasite1" {
#     schema = "Test1"
#     template = "Template1"
#     site = "Site1"

# }


