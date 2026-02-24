package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSONodeSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Node Settings Data Source") },
				Config:    testAccMSONodeSettingsDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "name", "tf_test_node_settings"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "description", "Terraform test Node Settings Policy"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.admin_state", "enabled"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "synce.quality_level", "option_2_generation_1"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.node_domain", "25"),
					resource.TestCheckResourceAttr("mso_fabric_policies_node_settings.node_settings", "ptp.priority_2", "99"),
				),
			},
		},
	})
}

func testAccMSONodeSettingsDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_policies_node_settings" "node_settings" {
	    template_id        = mso_fabric_policies_node_settings.node_settings.template_id
	    name               = "tf_test_node_settings"
    }`, testAccMSONodeSettingsConfigCreate())
}
