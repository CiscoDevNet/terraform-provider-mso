package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateBdDHCPPolicyDataSource(t *testing.T) {
	const dataSourceName = "data.mso_schema_template_bd_dhcp_policy.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "4.0")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Schema Template BD DHCP Policy Data Source - Not Found") },
				Config:      testAccMSOSchemaTemplateBdDHCPPolicyDataSourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`DHCP policy .* not found on Bridge Domain`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Schema Template BD DHCP Policy Data Source") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_option_policy_name", msoTenantPoliciesDhcpOptionPolicyName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Find a standalone DHCP policy by UUID") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyDataSourceUUIDConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr(dataSourceName, "dhcp_option_policy_name", msoTenantPoliciesDhcpOptionPolicyName),
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", "mso_schema_template_bd_dhcp_policy.test", "uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dhcp_option_policy_uuid", "mso_tenant_policies_dhcp_option_policy."+msoTenantPoliciesDhcpOptionPolicyName, "uuid"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateBdDHCPPolicyDataSourceNotFoundConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfig() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
data "mso_schema_template_bd_dhcp_policy" "test" {
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  bd_name       = mso_schema_template_bd.dhcp_parent.name
  name          = "%s"
}
`, "non_existing_dhcp_policy")
}

func testAccMSOSchemaTemplateBdDHCPPolicyDataSourceConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfig() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "test" {
  schema_id               = mso_schema_template_bd.dhcp_parent.schema_id
  template_name           = mso_schema_template_bd.dhcp_parent.template_name
  bd_name                 = mso_schema_template_bd.dhcp_parent.name
  name                    = mso_tenant_policies_dhcp_relay_policy.%s.name
  dhcp_option_policy_name = mso_tenant_policies_dhcp_option_policy.%s.name
}
data "mso_schema_template_bd_dhcp_policy" "test" {
  schema_id     = mso_schema_template_bd_dhcp_policy.test.schema_id
  template_name = mso_schema_template_bd_dhcp_policy.test.template_name
  bd_name       = mso_schema_template_bd_dhcp_policy.test.bd_name
  name          = mso_schema_template_bd_dhcp_policy.test.name
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpOptionPolicyName)
}

func testAccMSOSchemaTemplateBdDHCPPolicyDataSourceUUIDConfig() string {
	return testAccMSOSchemaTemplateBdDHCPPolicyPrerequisiteConfigWithTwoOptionPolicies() + testAccMSOSchemaTemplateBdDHCPPolicyBDConfig() + fmt.Sprintf(`
resource "mso_schema_template_bd_dhcp_policy" "test" {
  depends_on              = [mso_schema_template_bd.dhcp_parent]
  schema_id               = mso_schema_template_bd.dhcp_parent.schema_id
  template_name           = mso_schema_template_bd.dhcp_parent.template_name
  bd_name                 = mso_schema_template_bd.dhcp_parent.name
  uuid                    = mso_tenant_policies_dhcp_relay_policy.%s.uuid
  dhcp_option_policy_uuid = mso_tenant_policies_dhcp_option_policy.%s.uuid
}
data "mso_schema_template_bd_dhcp_policy" "test" {
  schema_id     = mso_schema_template_bd_dhcp_policy.test.schema_id
  template_name = mso_schema_template_bd_dhcp_policy.test.template_name
  bd_name       = mso_schema_template_bd_dhcp_policy.test.bd_name
  uuid          = mso_tenant_policies_dhcp_relay_policy.%s.uuid
}
`, msoTenantPoliciesDhcpRelayPolicyName, msoTenantPoliciesDhcpOptionPolicyName, msoTenantPoliciesDhcpRelayPolicyName)
}
