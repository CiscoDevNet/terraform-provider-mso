package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateBdDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "4.0")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Schema Template BD Data Source - Not Found") },
				Config:      testAccMSOSchemaTemplateBdDataSourceNotFound(),
				ExpectError: regexp.MustCompile(`Unable to find the BD`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Schema Template BD Data Source") },
				Config:    testAccMSOSchemaTemplateBdDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "display_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "description", "Terraform test BD"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_unknown_unicast", "flood"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "unknown_multicast_flooding", "flood"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "multi_destination_flooding", "flood_in_bd"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "ipv6_unknown_multicast_flooding", "flood"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "arp_flooding", "true"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "unicast_routing", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "layer3_multicast", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "intersite_bum_traffic", "true"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "optimize_wan_bandwidth", "false"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "layer2_stretch", "true"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "virtual_mac_address", ""),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "ep_move_detection_mode", "none"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd."+msoSchemaTemplateBdName, "vrf_name", msoSchemaTemplateVrfL3MulticastName),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read inline DHCP policy UUIDs from the Bridge Domain data source") },
				Config:    testAccMSOSchemaTemplateBdDHCPPolicyInlineUUIDDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_schema_template_bd.dhcp_parent", "dhcp_policies.#", "1"),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd.dhcp_parent", "dhcp_policies.0.name", msoTenantPoliciesDhcpRelayPolicyName),
					resource.TestCheckResourceAttr("data.mso_schema_template_bd.dhcp_parent", "dhcp_policies.0.dhcp_option_policy_name", msoTenantPoliciesDhcpOptionPolicyName),
					resource.TestCheckResourceAttrPair("data.mso_schema_template_bd.dhcp_parent", "dhcp_policies.0.uuid", "mso_tenant_policies_dhcp_relay_policy."+msoTenantPoliciesDhcpRelayPolicyName, "uuid"),
					resource.TestCheckResourceAttrPair("data.mso_schema_template_bd.dhcp_parent", "dhcp_policies.0.dhcp_option_policy_uuid", "mso_tenant_policies_dhcp_option_policy."+msoTenantPoliciesDhcpOptionPolicyName, "uuid"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateBdDataSource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_template_bd" "%[2]s" {
		schema_id     = mso_schema_template_bd.%[2]s.schema_id
		template_name = mso_schema_template_bd.%[2]s.template_name
		name          = "%[2]s"
	}`, testAccMSOSchemaTemplateBdConfigCreate(), msoSchemaTemplateBdName)
}

func testAccMSOSchemaTemplateBdDataSourceNotFound() string {
	return fmt.Sprintf(`%[1]s
	data "mso_schema_template_bd" "%[2]s" {
		schema_id     = mso_schema_template_bd.%[2]s.schema_id
		template_name = mso_schema_template_bd.%[2]s.template_name
		name          = "non_existing_bd"
	}`, testAccMSOSchemaTemplateBdConfigCreate(), msoSchemaTemplateBdName)
}

func testAccMSOSchemaTemplateBdDHCPPolicyInlineUUIDDataSourceConfig() string {
	return fmt.Sprintf(`%s
data "mso_schema_template_bd" "dhcp_parent" {
  schema_id     = mso_schema_template_bd.dhcp_parent.schema_id
  template_name = mso_schema_template_bd.dhcp_parent.template_name
  name          = mso_schema_template_bd.dhcp_parent.name
}
`, testAccMSOSchemaTemplateBdDHCPPolicyInlineUUIDCreateConfig())
}
