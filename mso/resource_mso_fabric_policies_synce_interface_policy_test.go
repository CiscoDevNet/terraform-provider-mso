package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOSyncEInterfacePolicyResource(t *testing.T) {
	resourceName := "mso_fabric_policies_synce_interface_policy.synce_interface_policy"
	var templateID, uuid string

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
					fmt.Println("Test: Recreate SyncE Interface Policy after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricPolicyTemplate", "template", "syncEthIntfPolicies"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOSyncEInterfacePolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "name", "tf_test_synce_interface_policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "description", "Terraform test SyncE Interface Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "sync_state_msg", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "selection_input", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "src_priority", "120"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "wait_to_restore", "6"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
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
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "src_priority", "120"),
					resource.TestCheckResourceAttr("mso_fabric_policies_synce_interface_policy.synce_interface_policy", "wait_to_restore", "6"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import SyncE Interface Policy") },
				ResourceName:      "mso_fabric_policies_synce_interface_policy.synce_interface_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing SyncE Interface Policy") },
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
	}`, testAccMSOTemplateResourceFabricPolicyConfig())
}
