provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
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