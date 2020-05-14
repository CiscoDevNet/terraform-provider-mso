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

# resource "mso_schema_template_externalepg_subnet" "subnet1" {
#   schema_id = "5ea809672c00003bc40a2799"
#   template_name = "Template1"
#   externalepg_name =  "UntitledExternalEPG1"
#   ip = "10.101.100.0/0"
#   name = "sddfgbany"
#   scope = ["shared-rtctrl", "export-rtctrl"]
#   aggregate = ["shared-rtctrl", "export-rtctrl"]
# }
# data "mso_schema_template_externalepg_subnet" "subnet1" {
#   schema_id = "5c6c16d7270000c710f8094d"
#   template_name = "Template1"
#   externalepg_name = "Internet"
#   ip = "30.1.1.0/24"
# }

# output "demo" {
#   value = "${data.mso_schema_template_externalepg_subnet.subnet1}"
# }

# resource "mso_schema_template_anp_epg_contract" "contract1" {
  # schema_id = "5e2dd7112c00005db60a268b"
  # template_name = "Template1"
  # anp_name = "ANP-Financial"
  # epg_name = "Web"
#   contract_name = "Web-to-Internet-Financial"
#   relationship_type = "provider"
  
# }

# resource "mso_schema_site_anp_epg" "site_anp_epg" {
#   schema_id = "5c4d9fca270000a101f8094a"
#   template_name = "Template1"
#   site_id = "5c7c95d9510000cf01c1ee3d"
#   anp_name = "ANP"
#   epg_name = "DB"
# }

# data "mso_schema_site_anp_epg" "anpEpg" {
#   schema_id = "5c4d5bb72700000401f80948"
#   template_name = "Template1"
#   site_id = "5c7c95b25100008f01c1ee3c"
#   anp_name = "ANP"
#   epg_name = "DB"
# }
# output "demo" {
#   value = "${data.mso_schema_site_anp_epg.anpEpg}"
# }


# resource "mso_schema_site_anp_epg_static_port" "port1" {
#   schema_id = "5c4d5bb72700000401f80948"
#   site_id = "5c7c95b25100008f01c1ee3c"
#   template_name = "Template1"
#   anp_name = "ANP"
#   epg_name = "DB"
#   path_type = "port"
#   deployment_immediacy = "lazy"
#   pod = "pod-4"
#   leaf = "106"
#   path = "eth1/10"
#   vlan = 200
#   mode = "untagged"

  
# }

# data "mso_schema_site_anp_epg_static_port" "port1" {
#  schema_id = "5c4d5bb72700000401f80948"
#   site_id = "5c7c95b25100008f01c1ee3c"
#   template_name = "Template1"
#   anp_name = "ANP"
#   epg_name = "DB"
#   path_type = "port"
#   pod = "pod-6"
#   leaf = "108"
#   path = "eth1/10"
  
# }

# output "demo" {
#   value = "${data.mso_schema_site_anp_epg_static_port.port1}"
# }

# resource "mso_schema_site_anp_epg_domain" "site_anp_epg_domain" {
#   schema_id = "5c4d9fca270000a101f8094a"
#   template_name = "Template1"
#   site_id = "5c7c95b25100008f01c1ee3c"
#   anp_name = "ANP"
#   epg_name = "Web"
#   domain_type = "vmmDomain"
#   dn = "VMware-abcd"
#   deploy_immediacy = "immediate"
#   resolution_immediacy = "immediate"
#   vlan_encap_mode = "static"
#   allow_micro_segmentation = false
#   switching_mode = "native"
#   switch_type = "default"
#   micro_seg_vlan_type = "vlan"
#   micro_seg_vlan = 46
#   port_encap_vlan_type = "vlan"
#   port_encap_vlan = 45
#   enhanced_lagpolicy_name = "name"
#   enhanced_lagpolicy_dn = "dn"

# }

# data "mso_schema_site_anp_epg_domain" "anpEpgDomain" {
#   schema_id = "5c4d9fca270000a101f8094a"
#   template_name = "Template1"
#   site_id = "5c7c95b25100008f01c1ee3c"
#   anp_name = "ANP"
#   epg_name = "Web"
#   dn = "VMware-ab"
#   domain_type = "vmmDomain"
# }
# output "demo" {
#   value = "${data.mso_schema_site_anp_epg_domain.anpEpgDomain}"
# }

# resource "mso_schema_site_bd_l3out" "bdL3out" {
#   schema_id = "5d5dbf3f2e0000580553ccce"
#   template_name = "Template1"
#   site_id = "5c7c95b25100008f01c1ee3c"
#   bd_name = "WebServer-Finance"
#   l3out_name = "zzz"
# }


# data "mso_schema_site_bd_l3out" "bdL3out" {
#   schema_id = "5d5dbf3f2e0000580553ccce"
#   template_name = "Template1"
#   site_id = "5c7c95b25100008f01c1ee3c"
#   bd_name = "WebServer-Finance"
#   l3out_name = "ccc"
# }
# output "demo" {
#   value = "${data.mso_schema_site_bd_l3out.bdL3out}"
# }

# resource "mso_schema_site_vrf_region" "vrfRegion" {
#   schema_id = "5d5dbf3f2e0000580553ccce"
#   template_name = "Template1"
#   site_id = "5ce2de773700006a008a2678"
#   vrf_name = "Campus"
#   region_name = "region3"

# }

# data "mso_schema_site_vrf_region" "vrfRegion" {
#   schema_id = "5d5dbf3f2e0000580553ccce"
#   site_id = "5ce2de773700006a008a2678"
#   vrf_name = "Campus"
#   region_name = "westus"
# }
# output "demo" {
#   value = "${data.mso_schema_site_vrf_region.vrfRegion}"
# }

# resource "mso_schema_site_anp_epg_subnet" "subnet1" {
#   schema_id = "5c4d5bb72700000401f80948"
#   site_id = "5c4d5bb72700000401f80948"
#   template_name = "Template1"
#   anp_name = "ANP"
#   epg_name = "DB"
#   ip = "10.0.7.0/8"
#   scope = "private"
#   shared = false
  
# }

resource "mso_schema_site_vrf_region_cidr" "vrfRegionCidr" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  template_name = "Template1"
  site_id = "5ce2de773700006a008a2678"
  vrf_name = "Campus"
  region_name = "region1"
  ip = "2.2.2.2/2"
  primary = false

}
data "mso_schema_site_vrf_region_cidr" "vrfRegionCidr" {
  schema_id = "5d5dbf3f2e0000580553ccce"
  site_id = "5ce2de773700006a008a2678"
  vrf_name = "Campus"
  region_name = "westus"
  ip = "192.168.241.0/24"
}
output "demo" {
  value = "${data.mso_schema_site_vrf_region_cidr.vrfRegionCidr}"
}