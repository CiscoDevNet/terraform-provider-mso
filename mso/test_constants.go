package mso

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const msoTemplateSiteName1 = "ansible_test"
const msoTemplateSiteName2 = "ansible_test_2"

var msoTenantName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpEpgName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateVrfName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateExtEpgName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPolicyTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateMCPGlobalPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateBdName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPolicyTemplateIPSLAMonitoringPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPolicyTemplateIPSLATrackListName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateInterfaceSettingName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateL3DomainName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateSyncEInterfacePolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateMacsecPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateFilterName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateContractName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateBdL3MulticastName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateVrfL3MulticastName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

const msoSchemaTemplateAnpEpgSubnetIp = "10.0.0.1/24"
const msoSchemaTemplateAnpEpgSubnetIp2 = "10.0.0.2/24"

func testSiteConfigAnsibleTest() string {
	return fmt.Sprintf(`
data "mso_site" "%[1]s" {
	name = "%[1]s"
}
`, msoTemplateSiteName1)
}

func testSiteConfigAnsibleTest2() string {
	return fmt.Sprintf(`
data "mso_site" "%[1]s" {
	name = "%[1]s"
}
`, msoTemplateSiteName2)
}

func testTenantConfig() string {
	return fmt.Sprintf(`
resource "mso_tenant" "%[1]s" {
	name         = "%[1]s"
	display_name = "%[1]s"
	site_associations { 
		site_id = data.mso_site.%[2]s.id
	}
}
`, msoTenantName, msoTemplateSiteName1)
}

func testTenantPolicyTemplateConfig() string {
	return fmt.Sprintf(`
resource "mso_template" "%[1]s" {
	template_name = "%[1]s"
	template_type = "tenant"
	tenant_id     = mso_tenant.%[2]s.id
}
`, msoTenantPolicyTemplateName, msoTenantName)
}

func testSchemaConfig() string {
	return fmt.Sprintf(`
resource "mso_schema" "%[1]s" {
	name = "%[1]s"
	template {
		name         = "%[2]s"
		display_name = "%[2]s"
		tenant_id    = mso_tenant.%[3]s.id
	}
}
`, msoSchemaName, msoSchemaTemplateName, msoTenantName)
}

func testSchemaTemplateAnpConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp" "%[1]s" {
	name         = "%[1]s"
	display_name = "%[1]s"
	schema_id    = mso_schema.%[2]s.id
	template     = "%[3]s"
}
`, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateAnpEpgConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp_epg" "%[1]s" {
	name          = "%[1]s"
	display_name  = "%[1]s"
	anp_name      = mso_schema_template_anp.%[2]s.name
	schema_id     = mso_schema.%[3]s.id
	template_name = "%[4]s"
}
`, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateVrfConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_vrf" "%[1]s" {
	name         = "%[1]s"
	display_name = "%[1]s"
	schema_id    = mso_schema.%[2]s.id
	template     = "%[3]s"
}
`, msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateExtEpgConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_external_epg" "%[1]s" {
	external_epg_name = "%[1]s"
	display_name      = "%[1]s"
	vrf_name          = mso_schema_template_vrf.%[2]s.name
	schema_id         = mso_schema.%[3]s.id
	template_name     = "%[4]s"
}
`, msoSchemaTemplateExtEpgName, msoSchemaTemplateVrfName, msoSchemaName, msoSchemaTemplateName)
}

func testFabricPolicyTemplateConfig() string {
	return fmt.Sprintf(`
resource "mso_template" "%[1]s" {
	template_name = "%[1]s"
	template_type = "fabric_policy"
}
`, msoFabricPolicyTemplateName)
}

func testSchemaTemplateBdConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_bd" "%[1]s" {
	schema_id				= mso_schema.%[2]s.id
	template_name			= "%[3]s"
	name					= "%[1]s"
	display_name			= "%[1]s"
	layer2_unknown_unicast 	= "proxy"
	vrf_name				= mso_schema_template_vrf.%[4]s.name
}
`, msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testTenantPolicyTemplateIPSLAMonitoringPolicyConfig() string {
	return fmt.Sprintf(`
resource "mso_tenant_policies_ipsla_monitoring_policy" "%[1]s" {
	template_id        = mso_template.%[2]s.id
	name               = "%[1]s"
	sla_type           = "http"
	destination_port   = 80
	http_version       = "HTTP11"
	http_uri           = "/example"
	sla_frequency      = 120
	detect_multiplier  = 4
	request_data_size  = 64
	type_of_service    = 18
	operation_timeout  = 100
	threshold          = 100
	ipv6_traffic_class = 255
}`, msoTenantPolicyTemplateIPSLAMonitoringPolicyName, msoTenantPolicyTemplateName)
}

func testFabricPolicyTemplateL3DomainConfig() string {
	return fmt.Sprintf(`
resource "mso_fabric_policies_l3_domain" "%[1]s" {
	template_id    = mso_template.%[2]s.id
	name           = "%[1]s"
}
`, msoFabricPolicyTemplateL3DomainName, msoFabricPolicyTemplateName)
}

func testFabricPolicyTemplateSyncEInterfacePolicyConfig() string {
	return fmt.Sprintf(`
resource "mso_fabric_policies_synce_interface_policy" "%[1]s" {
	template_id     = mso_template.%[2]s.id
	name            = "%[1]s"
}
`, msoFabricPolicyTemplateSyncEInterfacePolicyName, msoFabricPolicyTemplateName)
}

func testFabricPolicyTemplateMacsecPolicyConfig() string {
	return fmt.Sprintf(`
resource "mso_fabric_policies_macsec_policy" "%[1]s" {
	template_id            = mso_template.%[2]s.id
	name                   = "%[1]s"
	interface_type         = "access"
	cipher_suite           = "256GcmAes"
	window_size            = 128
	security_policy        = "shouldSecure"
	sak_expire_time        = 60
	confidentiality_offset = "offset30"
	key_server_priority    = 8
}
`, msoFabricPolicyTemplateMacsecPolicyName, msoFabricPolicyTemplateName)
}

func testSchemaTemplateVrfL3MulticastConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_vrf" "%[1]s" {
	name             = "%[1]s"
	display_name     = "%[1]s"
	schema_id        = mso_schema.%[2]s.id
	template         = "%[3]s"
	layer3_multicast = true
	preferred_group  = true
}
`, msoSchemaTemplateVrfL3MulticastName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateBdL3MulticastConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_bd" "%[1]s" {
	schema_id				= mso_schema.%[2]s.id
	template_name			= "%[3]s"
	name					= "%[1]s"
	display_name			= "%[1]s"
	layer2_unknown_unicast 	= "proxy"
	vrf_name				= mso_schema_template_vrf.%[4]s.name
	layer3_multicast		= true
}
`, msoSchemaTemplateBdL3MulticastName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfL3MulticastName)
}

func testSchemaTemplateFilterEntryConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_filter_entry" "%[1]s" {
	schema_id          = mso_schema.%[2]s.id
	template_name      = "%[3]s"
	name               = "%[1]s"
	display_name       = "%[1]s"
	entry_name         = "%[1]s_entry"
	entry_display_name = "%[1]s_entry"
}
`, msoSchemaTemplateFilterName, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateContractConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_contract" "%[1]s" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	contract_name = "%[1]s"
	display_name  = "%[1]s"
	filter_type   = "bothWay"
	scope         = "context"
	filter_relationship {
		filter_name = mso_schema_template_filter_entry.%[4]s.name
		filter_type = "bothWay"
	}
}
`, msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateFilterName)
}

func testSchemaTemplateAnpEpgContractConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp_epg_contract" "%[1]s_provider" {
	schema_id         = mso_schema.%[2]s.id
	template_name     = "%[3]s"
	anp_name          = "%[4]s"
	epg_name          = "%[5]s"
	contract_name     = mso_schema_template_contract.%[1]s.contract_name
	relationship_type = "provider"
}
`, msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName)
}

func testSchemaTemplateAnpEpgSubnetConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_anp_epg_subnet" "%[1]s_subnet" {
	schema_id = mso_schema.%[2]s.id
	template  = "%[3]s"
	anp_name  = "%[4]s"
	epg_name  = mso_schema_template_anp_epg.%[5]s.name
	ip        = "%[6]s"
	scope     = "private"
	shared    = false
}
`, msoSchemaTemplateAnpEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateAnpName, msoSchemaTemplateAnpEpgName, msoSchemaTemplateAnpEpgSubnetIp)
}
