package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOMCPGlobalPolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: MCP Global Policy Data Source") },
				Config:    testAccMSOMCPGlobalPolicyDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "name", msoFabricPolicyTemplateMCPGlobalPolicyName),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "description", ""),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "admin_state", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "enable_mcp_pdu_per_vlan", "disabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "key", ""),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "loop_detect_multiplication_factor", "3"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "port_disable_protection", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "initial_delay_time", "180"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "transmission_frequency_sec", "2"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "transmission_frequency_msec", "0"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_policies_mcp_global_policy.mcp_global_policy", "template_id"),
				),
			},
		},
	})
}

func testAccMSOMCPGlobalPolicyDataSource() string {
	return fmt.Sprintf(`%[1]s
data "mso_fabric_policies_mcp_global_policy" "mcp_global_policy" {
	template_id = mso_template.%[2]s.id
	name        = mso_fabric_policies_mcp_global_policy.%[3]s.name
}`, testAccMSOFabricPoliciesMCPGlobalPolicyConfigCreate(), msoFabricPolicyTemplateName, msoFabricPolicyTemplateMCPGlobalPolicyName)
}
