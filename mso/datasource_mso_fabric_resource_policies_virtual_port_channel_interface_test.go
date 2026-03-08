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
					resource.TestCheckResourceAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"name",
						"tf_test_vpc_if",
					),
					resource.TestCheckResourceAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"description",
						"Terraform test VPC Interface",
					),
					resource.TestCheckResourceAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"node_1",
						"101",
					),
					resource.TestCheckResourceAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"node_2",
						"102",
					),
					resource.TestCheckResourceAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"node_1_interfaces.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"node_2_interfaces.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions.#",
						"1",
					),
					resource.TestCheckResourceAttrSet(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"uuid",
					),
				),
			},
		},
	})
}

func testAccMSOVirtualPortChannelInterfaceDataSource() string {
	return fmt.Sprintf(`%s
	data "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
	    template_id        = mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if.template_id
	    name               = "tf_test_vpc_if"
    }`, testAccMSOVirtualPortChannelInterfaceConfigCreate())
}
