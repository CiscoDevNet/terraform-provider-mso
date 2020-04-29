provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

# resource "mso_schema" "s01" {
#   name = "Shah"
#   template_name = "template99"
#   tenant_id = "5ea7e44b2c00007ebb0a2781"
  
# }

data "mso_schema_template" "st10" {
  name = "Template101"
  schema_id = "5c6c16d7270000c710f8094d"
}

output "demo1" {
  value = "${data.mso_schema_template.st10}"
}

resource "mso_schema_template" "st1" {
  schema_id = "${mso_schema.s01.id}"
  name = "Temp200"
  display_name = "Temp845"
  tenant_id = "5c4d9f3d2700007e01f80949"
  
}
resource "mso_tenant" "tenant1" {
		name = "mso1"
		display_name = "mso1"
}

data "mso_tenant" "schema1" {
  name = "Campus-Integration"
  display_name = "Campus-Integration"
}

output "demo1" {
  value = "${data.mso_tenant.schema1.id}"
}

resource "mso_site" "site1" {
  name = "mso2"
  username = "admin"
  password = "noir0!234"
  apic_site_id = "18"
  urls = [ "https://3.208.123.222" ]
}

data "mso_site" "schema10" {
  name = "AWS-West"
}

output "demo1" {
  value = "${data.mso_site.schema10.id}"
}
