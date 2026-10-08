package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOFabricPoliciesMCPGlobalPolicyResource(t *testing.T) {
	resourceName := "mso_fabric_policies_mcp_global_policy." + msoFabricPolicyTemplateMCPGlobalPolicyName
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create MCP Global Policy") },
				Config:    testAccMSOFabricPoliciesMCPGlobalPolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "name", msoFabricPolicyTemplateMCPGlobalPolicyName),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "enable_mcp_pdu_per_vlan", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "key", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "loop_detect_multiplication_factor", "3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "port_disable_protection", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "initial_delay_time", "180"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_sec", "2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_msec", "0"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "template_id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						templateID = rs.Primary.Attributes["template_id"]
						uuid = rs.Primary.Attributes["uuid"]
						if templateID == "" || uuid == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate MCP Global Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricPolicyTemplate", "template", "mcpGlobalPolicy"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOFabricPoliciesMCPGlobalPolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "name", msoFabricPolicyTemplateMCPGlobalPolicyName),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "enable_mcp_pdu_per_vlan", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "key", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "loop_detect_multiplication_factor", "3"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "port_disable_protection", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "initial_delay_time", "180"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_sec", "2"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_msec", "0"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "template_id"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import MCP Global Policy") },
				ResourceName:      "mso_fabric_policies_mcp_global_policy." + msoFabricPolicyTemplateMCPGlobalPolicyName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MCP Global Policy") },
				Config:    testAccMSOFabricPoliciesMCPGlobalPolicyConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "name", msoFabricPolicyTemplateMCPGlobalPolicyName+"_updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "description", "Updated MCP Global Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "enable_mcp_pdu_per_vlan", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "key", "test_key_123"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "loop_detect_multiplication_factor", "5"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "port_disable_protection", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "initial_delay_time", "360"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_sec", "5"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_msec", "500"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "template_id"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MCP Global Policy without key and admin_state is disabled") },
				Config:    testAccMSOFabricPoliciesMCPGlobalPolicyConfigUpdateWithoutKeyAndDisableAdminState(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "name", msoFabricPolicyTemplateMCPGlobalPolicyName+"_updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "description", "Updated MCP Global Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "enable_mcp_pdu_per_vlan", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "key", "test_key_123"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "loop_detect_multiplication_factor", "9"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "port_disable_protection", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "initial_delay_time", "200"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_sec", "15"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_msec", "900"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "template_id"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update MCP Global Policy without key when the admin_state enabled again") },
				Config:    testAccMSOFabricPoliciesMCPGlobalPolicyConfigUpdateWithoutKeyAndEnableAdminState(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "name", msoFabricPolicyTemplateMCPGlobalPolicyName+"_updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "description", "Updated MCP Global Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "enable_mcp_pdu_per_vlan", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "key", "test_key_123"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "loop_detect_multiplication_factor", "9"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "port_disable_protection", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "initial_delay_time", "200"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_sec", "15"),
					resource.TestCheckResourceAttr("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "transmission_frequency_msec", "900"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_policies_mcp_global_policy."+msoFabricPolicyTemplateMCPGlobalPolicyName, "template_id"),
				),
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing MCP Global Policy") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok || rs.Primary == nil {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					separator := strings.LastIndex(rs.Primary.ID, "/")
					if separator < 0 || separator == len(rs.Primary.ID)-1 {
						return "", fmt.Errorf("resource %s has an invalid import ID", resourceName)
					}
					return rs.Primary.ID[:separator+1] + "tf_missing_oob", nil
				},
				ExpectError: regexp.MustCompile(`cannot import .*: resource not found`),
			},
		},
	})
}

func testFabricPolicyTemplateMCPGlobalPolicyConfig() string {
	return fmt.Sprintf(`
resource "mso_fabric_policies_mcp_global_policy" "%[1]s" {
	template_id = mso_template.%[2]s.id
	name        = "%[1]s"
}
`, msoFabricPolicyTemplateMCPGlobalPolicyName, msoFabricPolicyTemplateName)
}

func testAccMSOFabricPoliciesMCPGlobalPolicyConfigCreate() string {
	return testFabricPolicyTemplateConfig() + testFabricPolicyTemplateMCPGlobalPolicyConfig()
}

func testAccMSOFabricPoliciesMCPGlobalPolicyConfigUpdate() string {
	return testFabricPolicyTemplateConfig() + fmt.Sprintf(`
resource "mso_fabric_policies_mcp_global_policy" "%[1]s" {
	template_id                        = mso_template.%[2]s.id
	name                               = "%[1]s_updated"
	description                        = "Updated MCP Global Policy"
	admin_state                        = "enabled"
	enable_mcp_pdu_per_vlan                = "enabled"
	key                                = "test_key_123"
	loop_detect_multiplication_factor  = 5
	port_disable_protection            = "enabled"
	initial_delay_time                 = 360
	transmission_frequency_sec         = 5
	transmission_frequency_msec        = 500
}
`, msoFabricPolicyTemplateMCPGlobalPolicyName, msoFabricPolicyTemplateName)
}

func testAccMSOFabricPoliciesMCPGlobalPolicyConfigUpdateWithoutKeyAndDisableAdminState() string {
	return testFabricPolicyTemplateConfig() + fmt.Sprintf(`
resource "mso_fabric_policies_mcp_global_policy" "%[1]s" {
	template_id                        = mso_template.%[2]s.id
	name                               = "%[1]s_updated"
	admin_state                        = "disabled"
	enable_mcp_pdu_per_vlan                = "disabled"
	loop_detect_multiplication_factor  = 9
	port_disable_protection            = "disabled"
	initial_delay_time                 = 200
	transmission_frequency_sec         = 15
	transmission_frequency_msec        = 900
}
`, msoFabricPolicyTemplateMCPGlobalPolicyName, msoFabricPolicyTemplateName)
}

func testAccMSOFabricPoliciesMCPGlobalPolicyConfigUpdateWithoutKeyAndEnableAdminState() string {
	return testFabricPolicyTemplateConfig() + fmt.Sprintf(`
resource "mso_fabric_policies_mcp_global_policy" "%[1]s" {
	template_id                        = mso_template.%[2]s.id
	name                               = "%[1]s_updated"
	admin_state                        = "enabled"
	enable_mcp_pdu_per_vlan                = "disabled"
	loop_detect_multiplication_factor  = 9
	port_disable_protection            = "disabled"
	initial_delay_time                 = 200
	transmission_frequency_sec         = 15
	transmission_frequency_msec        = 900
}
`, msoFabricPolicyTemplateMCPGlobalPolicyName, msoFabricPolicyTemplateName)
}
