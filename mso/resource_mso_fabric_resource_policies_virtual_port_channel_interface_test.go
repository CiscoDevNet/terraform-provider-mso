package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOVirtualPortChannelInterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Virtual Port Channel Interface Resource - Create")
				},
				Config: testAccMSOVirtualPortChannelInterfaceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "name", "tf_test_vpc_if"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "description", "Terraform test VPC Interface"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1", "101"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2", "102"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.0", "1/1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.1", "1/10-11"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.#", "1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.0", "1/2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_descriptions.#", "1"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "uuid"),
					customTestCheckResourceTypeSetAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "101",
							"interface":   "1/1",
							"description": "Terraform test interface description",
						},
					),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Virtual Port Channel Interface Resource - Update")
				},
				Config: testAccMSOVirtualPortChannelInterfaceConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "name", "tf_test_vpc_if"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "description", "Terraform test VPC Interface (updated)"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1", "103"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2", "104"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.#", "3"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.0", "1/1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.1", "1/2-5"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.2", "1/7"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.#", "3"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.0", "1/2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.1", "1/3"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.2", "1/5-7"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_descriptions.#", "2"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "uuid"),
					customTestCheckResourceTypeSetAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "103",
							"interface":   "1/1",
							"description": "Terraform test interface description 103",
						},
					),
					customTestCheckResourceTypeSetAttr(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "104",
							"interface":   "1/3",
							"description": "Terraform test interface description 104",
						},
					),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Virtual Port Channel Interface Resource - Removed Descriptions")
				},
				Config: testAccMSOVirtualPortChannelInterfaceConfigRemoveDescriptions(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "description", "Terraform test VPC Interface (removed descriptions)"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.#", "3"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.#", "3"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_descriptions.#", "0"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "uuid"),
				),
			},
			{
				ResourceName:      "mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments(
			"mso_fabric_resource_policies_virtual_port_channel_interface",
			"fabricResourceTemplate",
			"template",
			"virtualPortChannels",
		),
	})
}

func testAccMSOVirtualPortChannelInterfaceConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_template.template_fabric_resource.id
		name        = "tf_test_vpc_if"
		description = "Terraform test VPC Interface"
		node_1      = "101"
		node_2      = "102"

		node_1_interfaces = ["1/1", "1/10-11"]
		node_2_interfaces = ["1/2"]

		interface_policy_group_uuid = "c0e13c44-a1ee-4fe7-8630-c02e9b2520aa"
		// interface_policy_group_uuid will be validated/required once the
		// interface policy group resource merges (PR #449).

		interface_descriptions {
			node        = "101"
			interface   = "1/1"
			description = "Terraform test interface description"
		}
	}
	`, testAccMSOTemplateResourceFabricResourceConfig())
}

func testAccMSOVirtualPortChannelInterfaceConfigUpdate() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_template.template_fabric_resource.id
		name        = "tf_test_vpc_if"
		description = "Terraform test VPC Interface (updated)"
		node_1      = "103"
		node_2      = "104"

		node_1_interfaces = ["1/1", "1/2-5", "1/7"]
		node_2_interfaces = ["1/2", "1/3", "1/5-7"]

		interface_policy_group_uuid = "c0e13c44-a1ee-4fe7-8630-c02e9b2520aa"
		// interface_policy_group_uuid will be validated/required once the
		// interface policy group resource merges (PR #449).

		interface_descriptions {
			node        = "103"
			interface   = "1/1"
			description = "Terraform test interface description 103"
		}

		interface_descriptions {
			node        = "104"
			interface   = "1/3"
			description = "Terraform test interface description 104"
		}
	}
	`, testAccMSOTemplateResourceFabricResourceConfig())
}

func testAccMSOVirtualPortChannelInterfaceConfigRemoveDescriptions() string {
	return fmt.Sprintf(`%s
	resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_template.template_fabric_resource.id
		name        = "tf_test_vpc_if"
		description = "Terraform test VPC Interface (removed descriptions)"
		node_1      = "103"
		node_2      = "104"

		node_1_interfaces = ["1/1", "1/2-5", "1/7"]
		node_2_interfaces = ["1/2", "1/3", "1/5-7"]

		interface_policy_group_uuid = "c0e13c44-a1ee-4fe7-8630-c02e9b2520aa"
		// interface_policy_group_uuid will be validated/required once the
		// interface policy group resource lands (PR #449).
	}
	`, testAccMSOTemplateResourceFabricResourceConfig())
}
