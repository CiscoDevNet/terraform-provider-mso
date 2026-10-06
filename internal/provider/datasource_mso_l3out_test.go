package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMSOL3OutDataSource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_l3out")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Reject L3Out datasource without a name") },
				Config:    testAccMSOL3OutDataSourceMissingNameConfig(),
				ExpectError: regexp.MustCompile(
					`Missing required argument`,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reject L3Out datasource when the template does not exist") },
				Config:    testAccMSOL3OutDataSourceInvalidTemplateConfig(l3outName),
				ExpectError: regexp.MustCompile(
					`L3Out Template Not Found`,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read L3Out datasource by name") },
				Config:    testAccMSOL3OutDataSourceCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "template_id", "mso_template.l3out_test", "id"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "uuid", "mso_l3out.test", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "id", "mso_l3out.test", "id"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "name", l3outName),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "description", "initial L3Out description"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "vrf_uuid", "mso_schema_template_vrf.l3out_test_vrf_1", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "l3_domain", "mso_fabric_policies_l3_domain.l3out_test_domain_1", "name"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "target_dscp", "af11"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "pim_enabled", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "annotations.%", "1"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "annotations.owner", "network"),
					resource.TestCheckResourceAttrSet("data.mso_l3out.test", "id"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Reject L3Out datasource when the L3Out does not exist") },
				Config:    testAccMSOL3OutDataSourceMissingObjectConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(
					`L3Out Not Found`,
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read BGP-only L3Out and import route control") },
				Config:    testAccMSOL3OutDataSourceBGPOnlyConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "import_route_control_enabled", "true"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read OSPF-only L3Out") },
				Config:    testAccMSOL3OutDataSourceOSPFOnlyConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_l3out.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.area_id", "0.0.0.2"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.area_type", "regular"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Refresh L3Out datasource after the resource changes") },
				Config:    testAccMSOL3OutDataSourceUpdateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "template_id", "mso_template.l3out_test", "id"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "uuid", "mso_l3out.test", "uuid"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "name", l3outName),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "description", "updated L3Out description"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "vrf_uuid", "mso_schema_template_vrf.l3out_test_vrf_1", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "l3_domain", "mso_fabric_policies_l3_domain.l3out_test_domain_2", "name"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "target_dscp", "expedited_forwarding"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "pim_enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "import_route_control_enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "bgp.enabled", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.enabled", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.area_id", "0.0.0.1"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.area_type", "regular"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.cost", "10"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.send_redistributed_lsa", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.originate_summary_lsa", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.suppress_forwarding_address_in_translated_lsa", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "originate_default_route", "only"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "ospf.originate_default_route_always", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "annotations.%", "2"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "annotations.owner", "operations"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "annotations.purpose", "routing"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read cleared L3Out annotations from datasource") },
				Config:    testAccMSOL3OutDataSourceClearAnnotationsConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_l3out.test", "annotations.%", "0"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read L3Out interface groups from the data source") },
				Config:    testAccMSOL3OutDataSourceInterfaceGroupsConfig(testAccMSOL3OutInterfaceGroupsCreateConfig(siteName, l3outName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.%", "4"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.backup.description", "backup group"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.bfd.key_id", "20"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key_id", "30"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.ospf.authentication_type", "md5"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.middle.description", "middle group"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.tail.description", "tail group"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Refresh L3Out interface groups after an update") },
				Config:    testAccMSOL3OutDataSourceInterfaceGroupsConfig(testAccMSOL3OutInterfaceGroupsUpdateConfig(siteName, l3outName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.%", "4"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.bfd.key_id", "21"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.bfd_multi_hop.key_id", "31"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.ospf.authentication_type", "simple"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.backup.description", "updated backup"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.middle.description", "middle group"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.tail.description", "tail group"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read cleared L3Out interface groups") },
				Config:    testAccMSOL3OutDataSourceInterfaceGroupsConfig(testAccMSOL3OutInterfaceGroupsEmptyConfig(siteName, l3outName)),
				Check:     resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.%", "0"),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read all interface group tenant references and NetFlow monitor types") },
				Config:    testAccMSOL3OutDataSourceInterfaceGroupsConfig(testAccMSOL3OutInterfaceGroupAllReferencesConfig(siteName, l3outName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "interface_groups.edge.interface_routing_policy_uuid", "mso_tenant_policies_l3out_interface_routing_policy.routing", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "interface_groups.edge.custom_qos_policy_uuid", "mso_tenant_policies_custom_qos_policy.qos", "uuid"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.qos_priority", "level6"),
					resource.TestCheckResourceAttr("data.mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.%", "4"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ipv4", "mso_tenant_policies_netflow_monitor.ipv4", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ipv6", "mso_tenant_policies_netflow_monitor.ipv6", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.ce", "mso_tenant_policies_netflow_monitor.ce", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out.test", "interface_groups.edge.netflow_monitor_uuids.unspecified", "mso_tenant_policies_netflow_monitor.unspecified", "uuid"),
				),
			},
		},
	})
}

func testAccMSOL3OutDataSourceMissingNameConfig() string {
	return `
data "mso_l3out" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
}
`
}

func testAccMSOL3OutDataSourceInvalidTemplateConfig(l3outName string) string {
	return `
data "mso_l3out" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  name        = "` + l3outName + `"
}
`
}

func testAccMSOL3OutDataSourceCreateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "initial L3Out description"
  l3_domain   = mso_fabric_policies_l3_domain.l3out_test_domain_1.name
  target_dscp = "af11"
  pim_enabled = true
  annotations = { owner = "network" }
}

data "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = mso_l3out.test.name
}
`
}

func testAccMSOL3OutDataSourceMissingObjectConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description = "initial L3Out description"
  l3_domain   = mso_fabric_policies_l3_domain.l3out_test_domain_1.name
  target_dscp = "af11"
  pim_enabled = true
}

data "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = "${local.l3out_name}_missing"
}
`
}

func testAccMSOL3OutDataSourceBGPOnlyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id                  = mso_template.l3out_test.id
  name                         = local.l3out_name
  vrf_uuid                     = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  import_route_control_enabled = true

  bgp = {
    enabled = true
  }
}

data "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = mso_l3out.test.name
}
`
}

func testAccMSOL3OutDataSourceOSPFOnlyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  bgp         = {}

  ospf = {
    enabled   = true
    area_id   = "0.0.0.2"
    area_type = "regular"
  }
}

data "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = mso_l3out.test.name
}
`
}

func testAccMSOL3OutDataSourceUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id                  = mso_template.l3out_test.id
  name                         = local.l3out_name
  vrf_uuid                     = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  description                  = "updated L3Out description"
  l3_domain                    = mso_fabric_policies_l3_domain.l3out_test_domain_2.name
  target_dscp                  = "expedited_forwarding"
  pim_enabled                  = false
  import_route_control_enabled = false

  bgp = {
    enabled = true
  }
  annotations = {
    owner   = "operations"
    purpose = "routing"
  }
  ospf = {
    enabled                                        = true
    area_id                                        = "0.0.0.1"
    area_type                                      = "regular"
    cost                                           = 10
    send_redistributed_lsa                         = true
    originate_summary_lsa                          = true
    suppress_forwarding_address_in_translated_lsa = true
    originate_default_route_always                 = false
  }
  originate_default_route = "only"
}

data "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = mso_l3out.test.name
}
`
}

func testAccMSOL3OutDataSourceClearAnnotationsConfig(siteName, l3outName string) string {
	return testAccMSOL3OutPrerequisites(siteName, l3outName) + `
resource "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = local.l3out_name
  vrf_uuid    = mso_schema_template_vrf.l3out_test_vrf_1.uuid
  annotations = {}
}

data "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = mso_l3out.test.name
}
`
}

func testAccMSOL3OutDataSourceInterfaceGroupsConfig(config string) string {
	return config + `
data "mso_l3out" "test" {
  template_id = mso_template.l3out_test.id
  name        = mso_l3out.test.name
  depends_on  = [mso_l3out.test]
}
`
}
