package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMSOL3OutInterfaceGroupPolicyDataSource(t *testing.T) {
	siteName := testAccSiteName()
	l3outName := testAccResourceName("terraform_interface_group_data")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccProviderPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Reject group policy datasource with missing template") },
				Config:      testAccMSOL3OutInterfaceGroupPolicyDataSourceInvalidTemplateConfig(),
				ExpectError: regexp.MustCompile(`L3Out Template Not Found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read a standalone L3Out interface group policy") },
				Config:    testAccMSOL3OutInterfaceGroupPolicyDataSourceCreateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "id", "mso_l3out_interface_group_policy.test", "id"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "template_id", "mso_template.l3out_test", "id"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "l3out_uuid", "mso_l3out.test", "uuid"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "name", "test_interface_group"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "description", "initial"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "bfd.enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "bfd_multi_hop.enabled", "false"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "ospf.enabled", "false"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Refresh interface group policy datasource after authentication is configured")
				},
				Config: testAccMSOL3OutInterfaceGroupPolicyDataSourceUpdateConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "id", "mso_l3out_interface_group_policy.test", "id"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "bfd.authentication_enabled", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "bfd.key_id", "20"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "bfd_multi_hop.authentication_enabled", "true"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "bfd_multi_hop.key_id", "30"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "ospf.authentication_type", "md5"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "ospf.key_id", "10"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Read interface group policy tenant references and all NetFlow monitor types")
				},
				Config: testAccMSOL3OutInterfaceGroupPolicyDataSourceReferencesConfig(siteName, l3outName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "id", "mso_l3out_interface_group_policy.test", "id"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "interface_routing_policy_uuid", "mso_tenant_policies_l3out_interface_routing_policy.routing", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "custom_qos_policy_uuid", "mso_tenant_policies_custom_qos_policy.qos", "uuid"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "qos_priority", "level6"),
					resource.TestCheckResourceAttr("data.mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.%", "4"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ipv4", "mso_tenant_policies_netflow_monitor.ipv4", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ipv6", "mso_tenant_policies_netflow_monitor.ipv6", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.ce", "mso_tenant_policies_netflow_monitor.ce", "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_l3out_interface_group_policy.test", "netflow_monitor_uuids.unspecified", "mso_tenant_policies_netflow_monitor.unspecified", "uuid"),
				),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject group policy datasource with missing L3Out") },
				Config:      testAccMSOL3OutInterfaceGroupPolicyDataSourceMissingParentConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Not Found`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Reject group policy datasource with missing policy") },
				Config:      testAccMSOL3OutInterfaceGroupPolicyDataSourceMissingPolicyConfig(siteName, l3outName),
				ExpectError: regexp.MustCompile(`L3Out Interface Group Policy Not Found`),
			},
		},
	})
}

func testAccMSOL3OutInterfaceGroupPolicyDataSourceInvalidTemplateConfig() string {
	return `
data "mso_l3out_interface_group_policy" "test" {
  template_id = "00000000-0000-0000-0000-000000000000"
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  name        = "missing"
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyDataSourceCreateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyCreateConfig(siteName, l3outName) + testAccMSOL3OutInterfaceGroupPolicyDataSourceConfig()
}

func testAccMSOL3OutInterfaceGroupPolicyDataSourceUpdateConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyAuthenticatedConfig(siteName, l3outName) + testAccMSOL3OutInterfaceGroupPolicyDataSourceConfig()
}

func testAccMSOL3OutInterfaceGroupPolicyDataSourceReferencesConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyAllReferencesConfig(siteName, l3outName) + testAccMSOL3OutInterfaceGroupPolicyDataSourceConfig()
}

func testAccMSOL3OutInterfaceGroupPolicyDataSourceConfig() string {
	return `
data "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = mso_l3out_interface_group_policy.test.name

  depends_on = [mso_l3out_interface_group_policy.test]
}
`
}
func testAccMSOL3OutInterfaceGroupPolicyDataSourceMissingParentConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
data "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = "00000000-0000-0000-0000-000000000000"
  name        = "missing"
  depends_on  = [mso_l3out.test]
}
`
}

func testAccMSOL3OutInterfaceGroupPolicyDataSourceMissingPolicyConfig(siteName, l3outName string) string {
	return testAccMSOL3OutInterfaceGroupPolicyBaseConfig(siteName, l3outName) + `
data "mso_l3out_interface_group_policy" "test" {
  template_id = mso_template.l3out_test.id
  l3out_uuid  = mso_l3out.test.uuid
  name        = "missing"
}
`
}
