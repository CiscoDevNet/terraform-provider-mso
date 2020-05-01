provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema" "s01" {
  name = "Shah"
  template_name = "template99"
  tenant_id = "5ea7e44b2c00007ebb0a2781"
  
}

# data "mso_schema_template" "st10" {
#   name = "Template101"
#   schema_id = "5c6c16d7270000c710f8094d"
# }

# output "demo1" {
#   value = "${data.mso_schema_template.st10}"
# }

resource "mso_schema_template" "st1" {
  schema_id = "${mso_schema.s01.id}"
  name = "Temp200"
  display_name = "Temp845"
  tenant_id = "5c4d9f3d2700007e01f80949" 

}
resource "mso_tenant" "tenant1" {
	name = "m22"
	display_name = "m22"
  description = "sfdgnhjm"
  site_associations{site_id = "5c7c95b25100008f01c1ee3c"}
  user_associations{user_id = "0000ffff0000000000000020"}
}
# resource "mso_schema_template" "st1" {
#   schema_id = "${mso_schema.s01.id}"
#   name = "Temp200"
#   display_name = "Temp845"
#   tenant_id = "5c4d9f3d2700007e01f80949"
  
# }
# resource "mso_tenant" "tenant1" {
# 		name = "mso1"
# 		display_name = "mso1"
# }

# data "mso_tenant" "schema1" {
#   name = "Campus-Integration"
#   display_name = "Campus-Integration"
# }

output "demo1" {
  value = "${data.mso_tenant.schema1}"
}
# output "demo1" {
#   value = "${data.mso_tenant.schema1.id}"
# }

# resource "mso_site" "site1" {
#   name = "mso2"
#   username = "admin"
#   password = "noir0!234"
#   apic_site_id = "18"
#   urls = [ "https://3.208.123.222" ]
# }

# data "mso_site" "schema10" {
#   name = "AWS-West"
# }

# output "demo1" {
#   value = "${data.mso_site.schema10.id}"
# }

resource "mso_schema_template_bd_subnet" "bdsub1" {
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  bd_name = "testBD"
  ip = "10.26.17.0/8"
  scope = "public"
  shared = false
  no_default_gateway = true
  querier = true
  
}

# data "mso_schema_template_bd_subnet" "sbd10" {
#   schema_id = "5ea809672c00003bc40a2799"
#   template_name = "Template1"
#   bd_name = "testBD"
#   ip = "10.22.13.0/8"

  
# }
# output "demo" {
#   value = "${data.mso_schema_template_bd_subnet.sbd10}"
# }

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  anp_name = "ANP"
  name = "mso_epg1"
  bd_name = "BD1"
  vrf_name = "DEVNET-VRF"
  display_name = "mso_epg1"
  useg_epg = true
  intra_epg = "unenforced"
  intersite_multicaste_source = false
  preferred_group = false
}

data "mso_schema_template_anp_epg" "sepg1" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  anp_name = "ANP"
  name = "mso_epg1"

}

output "demo" {
  value = "${data.mso_schema_template_anp_epg.sepg1}"
}