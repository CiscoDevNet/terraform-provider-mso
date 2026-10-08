package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSONodeSettingsResource(t *testing.T) {
	resourceName := "mso_fabric_policies_node_settings.node_settings"
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Node Settings Policy") },
				Config:    testAccMSONodeSettingsConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "name", "tf_test_node_settings"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "description", "Terraform test Node Settings Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.0.quality_level", "option_2_generation_1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.0.node_domain", "25"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.0.priority_2", "99"),
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
					fmt.Println("Test: Recreate Node Settings Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricPolicyTemplate", "template", "nodePolicyGroups"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSONodeSettingsConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "name", "tf_test_node_settings"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "description", "Terraform test Node Settings Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.0.admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.0.quality_level", "option_2_generation_1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.0.node_domain", "25"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.0.priority_2", "99"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Node Settings Policy") },
				Config:    testAccMSONodeSettingsConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "name", "tf_test_node_settings_new"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "description", "Terraform test Node Settings Policy updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.0.admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.0.quality_level", "option_1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.0.node_domain", "30"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.0.priority_2", "100"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Node Settings Policy Remove SyncE and PTP") },
				Config:    testAccMSONodeSettingsConfigUpdateRemove(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "name", "tf_test_node_settings_new"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.#", "0"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Node Settings Policy") },
				ResourceName:      "mso_fabric_policies_node_settings.node_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Node Settings Policy") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_node_settings", "fabricPolicyTemplate", "template", "nodePolicyGroups"),
	})
}

func testAccMSONodeSettingsConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_node_settings" "node_settings" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_node_settings"
		description     = "Terraform test Node Settings Policy"
		synce {
			admin_state   = "enabled"
			quality_level = "option_2_generation_1"
		}
		ptp {
			node_domain = 25
			priority_2  = 99
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSONodeSettingsConfigUpdate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_node_settings" "node_settings" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_node_settings_new"
		description     = "Terraform test Node Settings Policy updated"
		synce {
			admin_state   = "disabled"
			quality_level = "option_1"
		}
		ptp {
			node_domain = 30
			priority_2  = 100
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSONodeSettingsConfigUpdateRemove() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_node_settings" "node_settings" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_node_settings_new"
		description     = ""
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
