package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOFabricResourcePortChannelInterfaceResource(t *testing.T) {
	resourceName := "mso_fabric_resource_policies_port_channel_interface." + msoFabricResourcePortChannelInterfaceName
	var templateID, uuid string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Port Channel Interface") },
				Config:    testAccMSOFabricResourcePortChannelInterfaceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "name", msoFabricResourcePortChannelInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "node", "101"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "template_id"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_policy_group_uuid"),
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
					fmt.Println("Test: Recreate Port Channel Interface after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricResourceTemplate", "template", "portChannels"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOFabricResourcePortChannelInterfaceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "name", msoFabricResourcePortChannelInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "node", "101"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "template_id"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_policy_group_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Port Channel Interface adding interface descriptions") },
				Config:    testAccMSOFabricResourcePortChannelInterfaceConfigUpdateAddingInterfaceDescriptions(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "name", msoFabricResourcePortChannelInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "description", "Terraform test Port Channel Interface updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "node", "101"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Port Channel Interface adding extra interface description") },
				Config:    testAccMSOFabricResourcePortChannelInterfaceConfigUpdateAddingExtraInterfaceDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "name", msoFabricResourcePortChannelInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "description", "Terraform test Port Channel Interface updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "node", "101"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "Interface Description 1/2",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Port Channel Interface removing extra interface description") },
				Config:    testAccMSOFabricResourcePortChannelInterfaceConfigUpdateRemovingExtraInterfaceDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "name", msoFabricResourcePortChannelInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "node", "101"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "",
						},
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Port Channel Interface") },
				ResourceName:      "mso_fabric_resource_policies_port_channel_interface." + msoFabricResourcePortChannelInterfaceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Port Channel Interface") },
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
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_resource_policies_port_channel_interface", "fabricResourceTemplate", "template", "portChannels"),
	})
}

var fabricResourcePortChannelInterfacePreConfig = testFabricResourceTemplateConfig() + testFabricPolicyTemplateConfig() + testFabricPoliciesInterfaceSettingPortChannelConfig()

func testAccMSOFabricResourcePortChannelInterfaceConfigCreate() string {
	return fmt.Sprintf(`%[1]s
    resource "mso_fabric_resource_policies_port_channel_interface" "%[2]s" {
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s"
        node                        = "101"
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_portchannel.uuid
    }`, fabricResourcePortChannelInterfacePreConfig, msoFabricResourcePortChannelInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePortChannelInterfaceConfigUpdateAddingInterfaceDescriptions() string {
	return fmt.Sprintf(`%[1]s
    resource "mso_fabric_resource_policies_port_channel_interface" "%[2]s" {
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s"
        description                 = "Terraform test Port Channel Interface updated"
        node                        = "101"
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_portchannel.uuid
        interface_descriptions {
            interface   = "1/1"
            description = "Interface Description 1/1"
        }
    }`, fabricResourcePortChannelInterfacePreConfig, msoFabricResourcePortChannelInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePortChannelInterfaceConfigUpdateAddingExtraInterfaceDescription() string {
	return fmt.Sprintf(`%[1]s
    resource "mso_fabric_resource_policies_port_channel_interface" "%[2]s" {
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s"
        description                 = "Terraform test Port Channel Interface updated"
        node                        = "101"
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_portchannel.uuid
        interface_descriptions {
            interface   = "1/1"
            description = "Interface Description 1/1"
        }
        interface_descriptions {
            interface   = "1/2"
            description = "Interface Description 1/2"
        }
    }`, fabricResourcePortChannelInterfacePreConfig, msoFabricResourcePortChannelInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePortChannelInterfaceConfigUpdateRemovingExtraInterfaceDescription() string {
	return fmt.Sprintf(`%[1]s
    resource "mso_fabric_resource_policies_port_channel_interface" "%[2]s" {
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s"
        node                        = "101"
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_portchannel.uuid
        interface_descriptions {
            interface   = "1/2"
        }
    }`, fabricResourcePortChannelInterfacePreConfig, msoFabricResourcePortChannelInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}
