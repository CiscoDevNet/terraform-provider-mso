package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOVirtualPortChannelInterfaceResource(t *testing.T) {
	resourceName := "mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if"
	var templateID, uuid string

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
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces", "1/1"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces", "1/10-11"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.#", "1"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces", "1/2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_descriptions.#", "1"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_policy_group_uuid"),
					CustomTestCheckTypeSetElemAttrs(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "101",
							"interface":   "1/1",
							"description": "Terraform test interface description",
						},
					),
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
					fmt.Println("Test: Recreate Virtual Port Channel Interface after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricResourceTemplate", "template", "virtualPortChannels"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOVirtualPortChannelInterfaceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "name", "tf_test_vpc_if"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "description", "Terraform test VPC Interface"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1", "101"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2", "102"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.#", "2"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces", "1/1"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces", "1/10-11"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.#", "1"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces", "1/2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_descriptions.#", "1"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_policy_group_uuid"),
					CustomTestCheckTypeSetElemAttrs(
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
					fmt.Println("Test: Virtual Port Channel Interface Resource - Update Unset Interface Description")
				},
				Config: testAccMSOVirtualPortChannelInterfaceConfigUpdateUnsetDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "name", "tf_test_vpc_if"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "description", ""),
					CustomTestCheckTypeSetElemAttrs(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "101",
							"interface":   "1/1",
							"description": "",
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
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "name", "tf_test_vpc_if_new"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "description", "Terraform test VPC Interface (updated)"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1", "103"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2", "104"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces.#", "3"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces", "1/1"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces", "1/2-5"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_1_interfaces", "1/7"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces.#", "3"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces", "1/2"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces", "1/3"),
					testCheckTypeSetStringElemAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "node_2_interfaces", "1/5-7"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_descriptions.#", "2"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_policy_group_uuid"),
					CustomTestCheckTypeSetElemAttrs(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "103",
							"interface":   "1/1",
							"description": "Terraform test interface description 103",
						},
					),
					CustomTestCheckTypeSetElemAttrs(
						"mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
						"interface_descriptions",
						map[string]string{
							"node":        "104",
							"interface":   "1/3",
							"description": "",
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
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if", "interface_policy_group_uuid"),
				),
			},
			{
				ResourceName:      "mso_fabric_resource_policies_virtual_port_channel_interface.vpc_if",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Virtual Port Channel Interface") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments(
			"mso_fabric_resource_policies_virtual_port_channel_interface",
			"fabricResourceTemplate",
			"template",
			"virtualPortChannels",
		),
	})
}

func testAccMSOVirtualPortChannelInterfaceConfigCreate() string {
	return fmt.Sprintf(`%s%s
	resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_template.template_fabric_resource.id
		interface_policy_group_uuid = mso_fabric_policies_interface_setting.%s_portchannel.uuid
		name        = "tf_test_vpc_if"
		description = "Terraform test VPC Interface"
		node_1      = "101"
		node_2      = "102"
		node_1_interfaces = ["1/1", "1/10-11"]
		node_2_interfaces = ["1/2"]

		interface_descriptions {
			node        = "101"
			interface   = "1/1"
			description = "Terraform test interface description"
		}
	}
	`, testAccMSOTemplateResourceFabricResourceConfig(), testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigCreate(), msoFabricPolicyTemplateInterfaceSettingName)
}

func testAccMSOVirtualPortChannelInterfaceConfigUpdateUnsetDescription() string {
	return fmt.Sprintf(`%s%s
	resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_template.template_fabric_resource.id
		interface_policy_group_uuid = mso_fabric_policies_interface_setting.%s_portchannel.uuid
		name        = "tf_test_vpc_if"
		description = ""
		node_1      = "101"
		node_2      = "102"
		node_1_interfaces = ["1/1", "1/10-11"]
		node_2_interfaces = ["1/2"]

		interface_descriptions {
			node        = "101"
			interface   = "1/1"
			description = ""
		}
	}
	`, testAccMSOTemplateResourceFabricResourceConfig(), testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigCreate(), msoFabricPolicyTemplateInterfaceSettingName)
}

func testAccMSOVirtualPortChannelInterfaceConfigUpdate() string {
	return fmt.Sprintf(`%s%s
	resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_template.template_fabric_resource.id
		interface_policy_group_uuid = mso_fabric_policies_interface_setting.%s_portchannel.uuid
		name        = "tf_test_vpc_if_new"
		description = "Terraform test VPC Interface (updated)"
		node_1      = "103"
		node_2      = "104"
		node_1_interfaces = ["1/1", "1/2-5", "1/7"]
		node_2_interfaces = ["1/2", "1/3", "1/5-7"]

		interface_descriptions {
			node        = "103"
			interface   = "1/1"
			description = "Terraform test interface description 103"
		}

		interface_descriptions {
			node        = "104"
			interface   = "1/3"
			// No description on new description should be ""
		}
	}
	`, testAccMSOTemplateResourceFabricResourceConfig(), testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigCreate(), msoFabricPolicyTemplateInterfaceSettingName)
}

func testAccMSOVirtualPortChannelInterfaceConfigRemoveDescriptions() string {
	return fmt.Sprintf(`%s%s
	resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
		template_id = mso_template.template_fabric_resource.id
		interface_policy_group_uuid = mso_fabric_policies_interface_setting.%s_portchannel.uuid
		name        = "tf_test_vpc_if_new"
		description = "Terraform test VPC Interface (removed descriptions)"
		node_1      = "103"
		node_2      = "104"
		node_1_interfaces = ["1/1", "1/2-5", "1/7"]
		node_2_interfaces = ["1/2", "1/3", "1/5-7"]
	}
	`, testAccMSOTemplateResourceFabricResourceConfig(), testAccMSOFabricPoliciesInterfaceSettingPortChannelConfigCreate(), msoFabricPolicyTemplateInterfaceSettingName)
}
