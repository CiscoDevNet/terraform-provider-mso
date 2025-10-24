package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSyncEInterfacePolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create SyncE Interface Policy") },
				Config:    testAccMSOSyncEInterfacePolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "name", "tf_test_synce_interface_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "description", "Terraform test SyncE Interface Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "sync_state_msg", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "selection_input", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "src_priority", "120"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "wait_to_restore", "6"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update SyncE Interface Policy") },
				Config:    testAccMSOSyncEInterfacePolicyConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "name", "tf_test_synce_interface_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "description", "Terraform test SyncE Interface Policy updated"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "admin_state", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "sync_state_msg", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "selection_input", "disabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "src_priority", "100"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "wait_to_restore", "5"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import SyncE Interface Policy") },
				ResourceName:      "mso_fabric_policies_synce_interface_policy.synce_interface_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_policies_synce_interface_policy", "fabricPolicyTemplate", "template", "syncEthIntfPolicies"),
	})
}

func testAccMSOSyncEInterfacePolicyConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_synce_interface_policy" "synce_interface_policy" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_synce_interface_policy"
		description     = "Terraform test SyncE Interface Policy"
		admin_state     = "enabled"
		sync_state_msg  = "enabled"
		selection_input = "enabled"
		src_priority    = 120
		wait_to_restore = 6
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}

func testAccMSOSyncEInterfacePolicyConfigUpdate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_policies_synce_interface_policy" "synce_interface_policy" {
		template_id     = mso_template.template_fabric_policy.id
		name            = "tf_test_synce_interface_policy"
		description     = "Terraform test SyncE Interface Policy updated"
		admin_state     = "disabled"
		sync_state_msg  = "disabled"
		selection_input = "disabled"
		src_priority    = 100
		wait_to_restore = 5
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
