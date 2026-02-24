package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSONodeSettingsResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.quality_level", "option_2_generation_1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.node_domain", "25"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.priority_2", "99"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Node Settings Policy") },
				Config:    testAccMSONodeSettingsConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "name", "tf_test_node_settings"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "description", "Terraform test Node Settings Policy updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.quality_level", "option_1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.node_domain", "30"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.priority_2", "100"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Node Settings Policy Remove SyncE and PTP") },
				Config:    testAccMSONodeSettingsConfigUpdateRemove(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "name", "tf_test_node_settings"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "description", "Terraform test Node Settings Policy removed synce and ptp"),
					resource.TestCheckNoResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce"),
					resource.TestCheckNoResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Node Settings Policy") },
				ResourceName:      "mso_fabric_policies_node_settings.node_settings",
				ImportState:       true,
				ImportStateVerify: true,
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
		synce = {
			admin_state = "enabled"
			quality_level = "option_2_generation_1"
		}
		ptp = {
			node_domain = 25
			priority_2 = 99
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSONodeSettingsConfigUpdate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_node_settings" "node_settings" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_node_settings"
		description     = "Terraform test Node Settings Policy updated"
		synce = {
			admin_state = "disabled"
			quality_level = "option_1"
		}
		ptp = {
			node_domain = 30
			priority_2 = 100
		}
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSONodeSettingsConfigUpdateRemove() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_node_settings" "node_settings" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_node_settings"
		description     = "Terraform test Node Settings Policy removed synce and ptp"
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
