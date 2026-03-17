package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOVirtualPortChannelInterfaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: VPC Interface Data Source") },
				Config:    testAccMSOVirtualPortChannelInterfaceDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "name", "tf_test_vpc_if"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "description", "Terraform test VPC Interface"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1", "101"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2", "102"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.#", "2"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.0", "1/1"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.1", "1/10-11"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.#", "1"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.0", "1/2"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_descriptions.#", "1"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_policy_group_uuid"),
					customTestCheckResourceTypeSetAttr(
						"data.mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "101",
							"interface":   "1/1",
							"description": "Terraform test interface description",
						},
					),
				),
			},
		},
	})
}

func testAccMSOVirtualPortChannelInterfaceDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if.template_id
		name        = "tf_test_vpc_if"
	}
	`, testAccMSOVirtualPortChannelInterfaceConfigCreate())
}
