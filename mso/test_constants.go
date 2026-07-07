package mso

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

const msoTemplateSiteName1 = "ansible_test"
const msoTemplateSiteName2 = "ansible_test_2"

// msoSchemaSiteResourceLabel1/2 are the Terraform resource block labels used
// for the two mso_schema_site resources in shared schema_site test
// configurations (see testSchemaSiteConfig and the
// testSchemaWithSingleSiteAssociation* helpers). Reusing these constants in
// nested schema_site_* tests keeps `depends_on` references and
// TestCheckResourceAttr addresses in sync with the emitted HCL.
const msoSchemaSiteResourceLabel1 = "site_1"
const msoSchemaSiteResourceLabel2 = "site_2"

var msoTenantName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpEpgName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateVrfName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateExtEpgName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateExtEpgSubnetName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateExtEpgSubnetName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPolicyTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateMCPGlobalPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateBdName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateBdName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPolicyTemplateIPSLAMonitoringPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPolicyTemplateIPSLATrackListName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateInterfaceSettingName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateL3DomainName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateSyncEInterfacePolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricPolicyTemplateMacsecPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateFilterName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateFilterName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateFilterEntryName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateContractName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateContractOneWayName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateBdL3MulticastName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateVrfL3MulticastName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateL3outName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPoliciesDhcpRelayPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPoliciesDhcpRelayPolicyName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoTenantPoliciesDhcpOptionPolicyName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpEpgUsegAttrName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateAnpEpgUsegAttrName2 = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoServiceDeviceTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoServiceDeviceClusterLbName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoServiceDeviceClusterFwName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoSchemaTemplateServiceGraphName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoServiceNodeTypeName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

// msoSchemaSiteServiceGraphDeviceName and msoSchemaSiteServiceGraphDeviceName2
// are randomised names for L4-L7 firewall devices created in msoTenantName on
// APIC by apicSetup() in testAPICPreCheck.
var msoSchemaSiteServiceGraphDeviceName = "fw1_" + acctest.RandStringFromCharSet(6, acctest.CharSetAlpha)
var msoSchemaSiteServiceGraphDeviceName2 = "fw2_" + acctest.RandStringFromCharSet(6, acctest.CharSetAlpha)
var msoSchemaSiteServiceGraphDeviceDn = "uni/tn-" + msoTenantName + "/lDevVip-" + msoSchemaSiteServiceGraphDeviceName
var msoSchemaSiteServiceGraphDeviceDn2 = "uni/tn-" + msoTenantName + "/lDevVip-" + msoSchemaSiteServiceGraphDeviceName2

// msoSchemaSiteContractServiceGraphDeviceName is the randomised name of the
// L4-L7 firewall device created in msoTenantName on APIC by apicSetup()
// in testAPICPreCheck.
var msoSchemaSiteContractServiceGraphDeviceName = "csg_fw_" + acctest.RandStringFromCharSet(6, acctest.CharSetAlpha)
var msoSchemaSiteContractServiceGraphDeviceDn = "uni/tn-" + msoTenantName + "/lDevVip-" + msoSchemaSiteContractServiceGraphDeviceName

// msoSchemaSiteContractServiceGraphProviderClusterInterface and
// msoSchemaSiteContractServiceGraphConsumerClusterInterface are randomised
// logical interface (lIf) names on the contract service graph firewall device.
// The update test step swaps these two values to exercise a config change.
var msoSchemaSiteContractServiceGraphProviderClusterInterface = "clu_if_" + acctest.RandStringFromCharSet(4, acctest.CharSetAlpha)
var msoSchemaSiteContractServiceGraphConsumerClusterInterface = "clu_if_" + acctest.RandStringFromCharSet(4, acctest.CharSetAlpha)

// Provider and consumer redirect policies use different tenants to exercise the
// cross-tenant redirect policy path. The provider uses msoTenantName (the schema
// tenant) and the consumer uses msoTenantName2.
var msoSchemaSiteContractServiceGraphProviderRedirectPolicy = "rp1_" + acctest.RandStringFromCharSet(6, acctest.CharSetAlpha)
var msoSchemaSiteContractServiceGraphConsumerRedirectPolicy = "rp2_" + acctest.RandStringFromCharSet(6, acctest.CharSetAlpha)

// msoSchemaSiteAnpEpgStaticLeafPath is the topology path of the leaf node
// used in static leaf acceptance tests. It must correspond to a real leaf
// switch onboarded to the ansible_test site in the lab.
const msoSchemaSiteAnpEpgStaticLeafPath = "topology/pod-1/node-101"

// msoSchemaSiteAnpEpgStaticPortPod, msoSchemaSiteAnpEpgStaticPortLeaf, and
// msoSchemaSiteAnpEpgStaticPortPath identify the physical interface used in
// static port acceptance tests. They must correspond to a real interface on a
// leaf switch onboarded to the ansible_test site in the lab.
// The assembled portPath (for path_type="port") is:
//
//	topology/{pod}/paths-{leaf}/pathep-[{path}]
//
// msoSchemaSiteAnpEpgStaticPortFex is the FEX extender ID used when testing
// the fex attribute. The assembled portPath with fex is:
//
//	topology/{pod}/paths-{leaf}/extpaths-{fex}/pathep-[{path}]
const msoSchemaSiteAnpEpgStaticPortPod = "pod-1"
const msoSchemaSiteAnpEpgStaticPortLeaf = "101"
const msoSchemaSiteAnpEpgStaticPortPath = "eth1/1"
const msoSchemaSiteAnpEpgStaticPortPath2 = "eth1/2"
const msoSchemaSiteAnpEpgStaticPortFex = "101"

const msoSchemaTemplateAnpEpgUsegAttrIp = "10.0.0.10/24"
const msoSchemaTemplateAnpEpgSubnetIp = "10.0.0.1/24"
const msoSchemaTemplateAnpEpgSubnetIp2 = "10.0.0.2/24"
const msoSchemaTemplateExtEpgSubnetIp = "10.0.1.1/24"
const msoSchemaTemplateExtEpgSubnetIp2 = "10.0.1.2/24"
const msoSchemaTemplateBdSubnetIp = "10.1.0.1/24"
const msoSchemaSiteBdSubnetIp = "10.2.0.1/24"
const msoSchemaSiteBdSubnetIp2 = "10.2.0.2/24"
const msoSchemaSiteAnpEpgSubnetIp = "10.3.0.1/24"

var msoFabricResourceTemplateName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricResourcePortChannelInterfaceName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
var msoFabricResourcePhysicalInterfaceName = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

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
	name       = "%[1]s"
	depends_on = [data.mso_site.%[2]s]
}
`, msoTemplateSiteName2, msoTemplateSiteName1)
}

func testTenantConfig() string {
	return testTenantConfigOneSite(msoTenantName)
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

// testSchemaConfigIgnoreTemplates creates a schema without templates and ignores external template changes.
// This is used when testing mso_schema_template separately, because mso_schema.Read() reads all templates
// from the API and writes them to state. Without ignore_changes, templates added by mso_schema_template
// would cause perpetual drift on the mso_schema resource.
func testSchemaConfigIgnoreTemplates() string {
	return fmt.Sprintf(`
resource "mso_schema" "%[1]s" {
	name = "%[1]s"
	lifecycle {
		ignore_changes = [template]
	}
}
`, msoSchemaName)
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

func testSchemaTemplateBdConfigWithName(bdName string) string {
	return fmt.Sprintf(`
resource "mso_schema_template_bd" "%[1]s" {
	schema_id				= mso_schema.%[2]s.id
	template_name			= "%[3]s"
	name					= "%[1]s"
	display_name			= "%[1]s"
	layer2_unknown_unicast 	= "proxy"
	vrf_name				= mso_schema_template_vrf.%[4]s.name
}
`, bdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testSchemaTemplateBdConfig() string {
	return testSchemaTemplateBdConfigWithName(msoSchemaTemplateBdName)
}

func testSchemaTemplateBdConfig2() string {
	return testSchemaTemplateBdConfigWithName(msoSchemaTemplateBdName2)
}

func testSchemaTemplateBdStretchedConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_bd" "%[1]s" {
	schema_id				= mso_schema.%[2]s.id
	template_name			= "%[3]s"
	name					= "%[1]s"
	display_name			= "%[1]s"
	layer2_unknown_unicast 	= "proxy"
	layer2_stretch			= true
	vrf_name				= mso_schema_template_vrf.%[4]s.name
}
`, msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testSchemaTemplateBdSubnetConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_bd_subnet" "%[1]s_subnet" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	bd_name       = mso_schema_template_bd.%[1]s.name
	ip            = "%[4]s"
	scope         = "private"
	shared        = false
}
`, msoSchemaTemplateBdName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateBdSubnetIp)
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

func testSchemaTemplateVrfContractConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_vrf_contract" "%[1]s_provider" {
	schema_id         = mso_schema.%[2]s.id
	template_name     = "%[3]s"
	vrf_name          = mso_schema_template_vrf.%[4]s.name
	relationship_type = "provider"
	contract_name     = mso_schema_template_contract.%[1]s.contract_name
}
`, msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
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

func testSchemaTemplateFilterEntryConfig2() string {
	return fmt.Sprintf(`
resource "mso_schema_template_filter_entry" "%[1]s" {
	schema_id          = mso_schema.%[2]s.id
	template_name      = "%[3]s"
	name               = "%[1]s"
	display_name       = "%[1]s"
	entry_name         = "%[1]s_entry"
	entry_display_name = "%[1]s_entry"
}
`, msoSchemaTemplateFilterName2, msoSchemaName, msoSchemaTemplateName)
}

func testSchemaTemplateL3outConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_l3out" "%[1]s" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	l3out_name    = "%[1]s"
	display_name  = "%[1]s"
	vrf_name      = mso_schema_template_vrf.%[4]s.name
}
`, msoSchemaTemplateL3outName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateVrfName)
}

func testSchemaTemplateExtEpgContractConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_external_epg_contract" "%[1]s_provider" {
	schema_id         = mso_schema.%[2]s.id
	template_name     = "%[3]s"
	external_epg_name = "%[4]s"
	contract_name     = mso_schema_template_contract.%[1]s.contract_name
	relationship_type = "provider"
}
`, msoSchemaTemplateContractName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName)
}

func testSchemaTemplateExtEpgSubnetConfig() string {
	return fmt.Sprintf(`
resource "mso_schema_template_external_epg_subnet" "%[1]s_subnet" {
	schema_id         = mso_schema.%[2]s.id
	template_name     = "%[3]s"
	external_epg_name = "%[4]s"
	ip                = "%[5]s"
}
`, msoSchemaTemplateExtEpgName, msoSchemaName, msoSchemaTemplateName, msoSchemaTemplateExtEpgName, msoSchemaTemplateExtEpgSubnetIp)
}

// testTenantConfigTwoSites creates a tenant associated with both ansible_test
// and ansible_test_2 sites, used by the schema_site acceptance tests.
func testTenantConfigTwoSites() string {
	return fmt.Sprintf(`
resource "mso_tenant" "%[1]s" {
	name = "%[1]s"
	site_associations {
		site_id = data.mso_site.%[2]s.id
	}
	site_associations {
		site_id = data.mso_site.%[3]s.id
	}
}
`, msoTenantName, msoTemplateSiteName1, msoTemplateSiteName2)
}

// testTenantConfigOneSite creates an mso_tenant resource associated with the
// primary site for the given tenant name.
func testTenantConfigOneSite(tenantName string) string {
	return fmt.Sprintf(`
resource "mso_tenant" "%[1]s" {
	name = "%[1]s"
	site_associations {
		site_id = data.mso_site.%[2]s.id
	}
}
`, tenantName, msoTemplateSiteName1)
}

// testSchemaSiteConfig emits a single mso_schema_site block referencing the
// shared schema/template. resourceLabel is used as the Terraform resource
// label (e.g. "site_1") and siteDataSource is the name of the existing
// data.mso_site source (e.g. msoTemplateSiteName1).
func testSchemaSiteConfig(resourceLabel, siteDataSource string, undeployOnDestroy bool) string {
	return fmt.Sprintf(`
resource "mso_schema_site" "%[1]s" {
	schema_id           = mso_schema.%[2]s.id
	site_id             = data.mso_site.%[3]s.id
	template_name       = tolist(mso_schema.%[2]s.template)[0].name
	undeploy_on_destroy = %[4]t
}
`, resourceLabel, msoSchemaName, siteDataSource, undeployOnDestroy)
}

// testSchemaTemplateDeployNdoConfig emits a mso_schema_template_deploy_ndo
// resource depending on the supplied list of resource references (e.g.
// "mso_schema_site.site_1", "mso_schema_template_vrf."+msoSchemaTemplateVrfName).
//
// Note: `force_apply = ""` is set explicitly to suppress a perpetual diff.
// The schema declares a default of "always-deploy", but
// resourceNDOSchemaTemplateDeployRead writes "" back to state on every
// refresh, so an unset config attribute would produce a non-empty plan after
// each apply ("" => "always-deploy"). Pinning the config to "" matches the
// post-Read state and keeps the test plan clean.
func testSchemaTemplateDeployNdoConfig(dependsOn []string) string {
	return fmt.Sprintf(`
resource "mso_schema_template_deploy_ndo" "deploy" {
	schema_id     = mso_schema.%[1]s.id
	template_name = tolist(mso_schema.%[1]s.template)[0].name
	force_apply   = ""
	depends_on    = [%[2]s]
}
`, msoSchemaName, strings.Join(dependsOn, ", "))
}

// testSchemaWithBothSitesPrerequisiteConfig emits both `data.mso_site` blocks
// (ansible_test and ansible_test_2), the shared tenant associated with both
// sites, and the schema with one template. This is the foundation used by
// schema_site and any nested schema_site_* acceptance tests.
func testSchemaWithBothSitesPrerequisiteConfig() string {
	return fmt.Sprintf(`%s%s%s%s`,
		testSiteConfigAnsibleTest(),
		testSiteConfigAnsibleTest2(),
		testTenantConfigTwoSites(),
		testSchemaConfig(),
	)
}

// testSchemaWithSingleSiteAssociationConfig extends the prerequisite config
// with a single mso_schema_site association (resource label `site_1`,
// `undeploy_on_destroy=false`). Use this as the minimum scaffolding for tests
// that exercise nested site-scoped children which do not require the template
// to be deployed.
func testSchemaWithSingleSiteAssociationConfig() string {
	return fmt.Sprintf(`%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaSiteConfig(msoSchemaSiteResourceLabel1, msoTemplateSiteName1, false),
	)
}

// testSchemaWithSingleSiteAssociationDeployedConfig adds a VRF and a
// mso_schema_template_deploy_ndo on top of testSchemaWithSingleSiteAssociationConfig
// so callers start from a deployed-template state. The site association uses
// `undeploy_on_destroy=true` so the test framework's destroy phase can undeploy
// before removing the association.
func testSchemaWithSingleSiteAssociationDeployedConfig() string {
	return fmt.Sprintf(`%s%s%s%s`,
		testSchemaWithBothSitesPrerequisiteConfig(),
		testSchemaSiteConfig(msoSchemaSiteResourceLabel1, msoTemplateSiteName1, true),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateDeployNdoConfig([]string{
			"mso_schema_site." + msoSchemaSiteResourceLabel1,
			"mso_schema_template_vrf." + msoSchemaTemplateVrfName,
		}),
	)
}

func testFabricResourceTemplateConfig() string {
	return fmt.Sprintf(`
resource "mso_template" "%[1]s" {
	template_name = "%[1]s"
	template_type = "fabric_resource"
}
	`, msoFabricResourceTemplateName)
}

func testFabricPoliciesInterfaceSettingPortChannelConfig() string {
	return fmt.Sprintf(`
resource "mso_fabric_policies_interface_setting" "%[1]s_portchannel" {
	template_id = mso_template.%[2]s.id
	type        = "portchannel"
	name        = "%[1]s_portchannel"
}
`, msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName)
}

func testFabricPoliciesInterfaceSettingPhysicalConfig() string {
	return fmt.Sprintf(`
resource "mso_fabric_policies_interface_setting" "%[1]s_physical" {
	template_id = mso_template.%[2]s.id
	type        = "physical"
	name        = "%[1]s_physical"
}
`, msoFabricPolicyTemplateInterfaceSettingName, msoFabricPolicyTemplateName)
}
