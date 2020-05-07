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

# data "mso_schema_template" "st10" {
#   name = "Template101"
#   schema_id = "5c6c16d7270000c710f8094d"
# }

# output "demo1" {
#   value = "${data.mso_schema_template.st10}"
# }

# resource "mso_schema_template" "st1" {
#   schema_id = "${mso_schema.s01.id}"
#   name = "Temp200"
#   display_name = "Temp845"
#   tenant_id = "5c4d9f3d2700007e01f80949" 

# }
# resource "mso_tenant" "tenant1" {
# 	name = "m22"
# 	display_name = "m22"
#   description = "sfdgnhjm"
#   site_associations{site_id = "5c7c95b25100008f01c1ee3c"}
#   user_associations{user_id = "0000ffff0000000000000020"}
# }
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

# output "demo1" {
#   value = "${data.mso_tenant.schema1}"
# }
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

# resource "mso_schema_template_bd_subnet" "bdsub1" {
#   schema_id = "5ea809672c00003bc40a2799"
#   template_name = "Template1"
#   bd_name = "testBD"
#   ip = "10.26.17.0/8"
#   scope = "public"
#   shared = false
#   no_default_gateway = true
#   querier = true
  
# }

# data "mso_schema_template_bd_subnet" "sbd10" {
#   schema_id = "5ea809672c00003bc40a2799"
#   template_name = "Template1"
#   bd_name = "testBD"
#   ip = "10.22.13.0/8"

  
# }
# output "demo" {
#   value = "${data.mso_schema_template_bd_subnet.sbd10}"
# }

# 

# resource "mso_schema_template_l3out" "template_l3out" {
#   schema_id = "5c6c16d7270000c710f8094d"
#   template_name = "Template1"
#   l3out_name = "l3out1"
#   display_name = "l3out1"
#   vrf_name = "vrf2"
# }

# data "mso_schema_template_l3out" "sl3out1" {
#   schema_id = "5c6c16d7270000c710f8094d"
#   template_name = "Template1"
#   l3out_name = "Internet_L3Out"
  
# }
# output "demo" {
#   value = "${data.mso_schema_template_l3out.sl3out1}"
# }

# resource "mso_schema_template_anp_epg" "anp_epg" {
#  schema_id = "5eafca7d2c000052860a2902"
#  template_name = "stemplate1"
#  anp_name = "sanp1"
#  name = "anpepg111"
#  bd_name = "testBD"
#  vrf_name = "vrf1"
#  display_name = "anpepg111"
#  useg_epg = true
#  intra_epg = "enforced"
#  intersite_multicaste_source = true
#  preferred_group = true
#  bd_template_name = "stemplate1"
#  vrf_schema_id = "5eafeb792c0000a18e0a2900"
#  bd_schema_id = "5eafeb792c0000a18e0a2900"
#  vrf_template_name = "stemplate1"

# } 

# resource "mso_schema_template_contract" "template_contract" {
#   schema_id = "5c4d5bb72700000401f80948"
#   template_name = "Template1"
#   contract_name = "C2"
#   display_name = "C2"
#   filter_type = "bothWay"
#   scope = "context"
#   filter_relationships = {
#     filter_schema_id = "5c4d5bb72700000401f80948"
#     filter_template_name = "Template1"
#     filter_name = "filter1"
#   }
#   directives = ["none"]
# }

# data "mso_schema_template_contract" "contract1" {
#   schema_id = "5c6c16d7270000c710f8094d"
#   template_name = "template1"
#   contract_name = "web2-to-DB2"
# } 
# output "demo" {
#   value = "${data.mso_schema_template_contract.contract1}"
# }

# resource "mso_schema_template_externalepg" "template_externalepg" {
#   schema_id = "5ea809672c00003bc40a2799"
#   template_name = "Template1"
#   externalepg_name = "external_epg12"
#   display_name = "external_epg12"
#   vrf_name = "vrf1"
# }

# data "mso_schema_template_externalepg" "externalEpg" {
#   schema_id = "5ea809672c00003bc40a2799"
#   template_name = "Template1"
#   externalepg_name = "UntitledExternalEPG1"
# }

# output "demo" {
#   value = "${data.mso_schema_template_externalepg.externalEpg}"
# }

resource "mso_schema_template_externalepg_subnet" "subnet1" {
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  externalepg_name =  "UntitledExternalEPG1"
  ip = "10.101.100.0/0"
  name = "sddfgbany"
  scope = ["shared-rtctrl", "export-rtctrl"]
  aggregate = ["shared-rtctrl", "export-rtctrl"]
}
data "mso_schema_template_externalepg_subnet" "subnet1" {
  schema_id = "5c6c16d7270000c710f8094d"
  template_name = "Template1"
  externalepg_name = "Internet"
  ip = "30.1.1.0/24"
}

output "demo" {
  value = "${data.mso_schema_template_externalepg_subnet.subnet1}"
}