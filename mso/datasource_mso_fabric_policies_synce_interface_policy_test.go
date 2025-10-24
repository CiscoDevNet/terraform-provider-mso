package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSyncEInterfacePolicyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: VLAN Pool Data Source") },
				Config:    testAccMSOSyncEInterfacePolicyDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_policies_synce_interface_policy.synce_interface_policy", "name", "tf_test_synce_interface_policy"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_synce_interface_policy.synce_interface_policy", "description", "Terraform test SyncE Interface Policy"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_synce_interface_policy.synce_interface_policy", "admin_state", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_synce_interface_policy.synce_interface_policy", "sync_state_msg", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_synce_interface_policy.synce_interface_policy", "selection_input", "enabled"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_synce_interface_policy.synce_interface_policy", "src_priority", "120"),
					resource.TestCheckResourceAttr("data.mso_fabric_policies_synce_interface_policy.synce_interface_policy", "wait_to_restore", "6"),
				),
			},
		},
	})
}

func testAccMSOSyncEInterfacePolicyDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_synce_interface_policy" "synce_interface_policy" {
	    template_id        = mso_fabric_policies_synce_interface_policy.synce_interface_policy.template_id
	    name               = "tf_test_synce_interface_policy"
    }`, testAccMSOSyncEInterfacePolicyConfigCreate())
}
