package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOFabricResourcePhysicalInterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Missing interface_policy_uuid and breakout_mode (error)") },
				Config:      testAccMSOFabricResourcePhysicalInterfaceConfigErrorMissingInterfacePolicyAndBreakoutMode(),
				ExpectError: regexp.MustCompile(`Either 'interface_policy_uuid' or 'breakout_mode' must be specified for creating a Physical Interface`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Physical Interface") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "name", msoFabricResourcePhysicalInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "policy_group_type", "physical"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "nodes.#", "1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "template_id"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_policy_uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Physical Interface adding interface descriptions") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceConfigUpdateAddingInterfaceDescriptions(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "name", msoFabricResourcePhysicalInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "description", "Terraform test Physical Interface updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "nodes.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Physical Interface adding extra interface description") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceConfigUpdateAddingExtraInterfaceDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "name", msoFabricResourcePhysicalInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "description", "Terraform test Physical Interface updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "nodes.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "Interface Description 1/2",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Physical Interface removing extra interface description") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceConfigUpdateRemovingExtraInterfaceDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "name", msoFabricResourcePhysicalInterfaceName+"_updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "nodes.#", "1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "",
						},
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Physical Interface") },
				ResourceName:      "mso_fabric_resource_policies_physical_interface." + msoFabricResourcePhysicalInterfaceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { fmt.Println("Test: Create Physical Interface with Breakout Mode") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "name", msoFabricResourcePhysicalInterfaceName+"_breakout"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "policy_group_type", "breakout"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "description", "Terraform test Physical Interface"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "nodes.#", "1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "breakout_mode", "4x10G"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "template_id"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Breakout Mode and add interface descriptions") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateAddingInterfaceDescriptions(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "name", msoFabricResourcePhysicalInterfaceName+"_breakout"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "description", "Terraform test Physical Interface updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "nodes.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "breakout_mode", "4x25G"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Breakout Mode and add extra interface description") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateAddingExtraInterfaceDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "name", msoFabricResourcePhysicalInterfaceName+"_breakout"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "description", "Terraform test Physical Interface updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "nodes.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "breakout_mode", "4x100G"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "Interface Description 1/2",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Breakout Mode and remove extra interface description") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingExtraInterfaceDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "name", msoFabricResourcePhysicalInterfaceName+"_breakout_updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "description", "Terraform test Physical Interface updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "nodes.#", "1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interfaces.#", "2"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "breakout_mode", "4x100G"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "Interface Description 1/2",
						},
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Physical Interface with Breakout Mode") },
				ResourceName:      "mso_fabric_resource_policies_physical_interface." + msoFabricResourcePhysicalInterfaceName + "_breakout",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_resource_policies_physical_interface", "fabricResourceTemplate", "template", "physicalInterfaces"),
	})
}

var fabricResourcePhysicalInterfacePreConfig = testFabricResourceTemplateConfig() + testFabricPolicyTemplateConfig() + testFabricPoliciesInterfaceSettingPhysicalConfig()

func testAccMSOFabricResourcePhysicalInterfaceConfigErrorMissingInterfacePolicyAndBreakoutMode() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id = mso_template.%[3]s.id
        name        = "%[2]s"
		description = "Terraform test Physical Interface"
        nodes       = ["101"]
        interfaces  = ["1/1","1/2"]
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id           = mso_template.%[4]s.id
        name                  = "%[2]s"
        nodes                 = ["101"]
        interfaces            = ["1/1","1/2"]
        interface_policy_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceConfigUpdateAddingInterfaceDescriptions() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id           = mso_template.%[4]s.id
        name                  = "%[2]s"
		description           = "Terraform test Physical Interface updated"
        nodes                 = ["101", "102"]
        interfaces            = ["1/1","1/2"]
        interface_policy_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
        interface_descriptions {
            interface   = "1/1"
            description = "Interface Description 1/1"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceConfigUpdateAddingExtraInterfaceDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id           = mso_template.%[4]s.id
        name                  = "%[2]s"
		description           = "Terraform test Physical Interface updated"
        nodes                 = ["101", "102"]
        interfaces            = ["1/1","1/2"]
        interface_policy_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
        interface_descriptions {
            interface   = "1/1"
            description = "Interface Description 1/1"
        }
        interface_descriptions {
            interface   = "1/2"
            description = "Interface Description 1/2"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceConfigUpdateRemovingExtraInterfaceDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id           = mso_template.%[4]s.id
        name                  = "%[2]s_updated"
        nodes                 = ["101"]
        interfaces            = ["1/1","1/2"]
        interface_policy_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
        interface_descriptions {
            interface   = "1/2"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigCreate() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s_breakout" {
        template_id   = mso_template.%[3]s.id
        name          = "%[2]s_breakout"
		description   = "Terraform test Physical Interface"
        nodes         = ["101"]
        interfaces    = ["1/1","1/2"]
        breakout_mode = "4x10G"
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateAddingInterfaceDescriptions() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s_breakout" {
        template_id   = mso_template.%[3]s.id
        name          = "%[2]s_breakout"
		description   = "Terraform test Physical Interface updated"
        nodes         = ["101", "102"]
        interfaces    = ["1/1","1/2"]
        breakout_mode = "4x25G"
        interface_descriptions {
            interface   = "1/1"
            description = "Interface Description 1/1"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateAddingExtraInterfaceDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s_breakout" {
        template_id   = mso_template.%[3]s.id
        name          = "%[2]s_breakout"
		description   = "Terraform test Physical Interface updated"
        nodes         = ["101", "102"]
        interfaces    = ["1/1","1/2"]
        breakout_mode = "4x100G"
        interface_descriptions {
            interface   = "1/1"
            description = "Interface Description 1/1"
        }
        interface_descriptions {
            interface   = "1/2"
            description = "Interface Description 1/2"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingExtraInterfaceDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s_breakout" {
        template_id   = mso_template.%[3]s.id
        name          = "%[2]s_breakout_updated"
		description   = "Terraform test Physical Interface updated"
        nodes         = ["101"]
        interfaces    = ["1/1","1/2"]
        breakout_mode = "4x100G"
        interface_descriptions {
            interface   = "1/2"
            description = "Interface Description 1/2"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricResourceTemplateName)
}
