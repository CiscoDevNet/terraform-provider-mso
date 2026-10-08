package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOFabricResourcePhysicalInterfaceResource(t *testing.T) {
	resourceName := "mso_fabric_resource_policies_physical_interface." + msoFabricResourcePhysicalInterfaceName
	var templateID, uuid string
	breakoutResourceName := "mso_fabric_resource_policies_physical_interface." + msoFabricResourcePhysicalInterfaceName + "_breakout"
	var breakoutTemplateID, breakoutUUID string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Create Without interface_policy_group_uuid and breakout_mode (error)") },
				Config:      testAccMSOFabricResourcePhysicalInterfaceConfigErrorMissingInterfacePolicyAndBreakoutMode(),
				ExpectError: regexp.MustCompile(`Either 'interface_policy_group_uuid' or 'breakout_mode' must be specified for creating a Physical Interface`),
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
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_policy_group_uuid"),
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
					fmt.Println("Test: Recreate Physical Interface after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "fabricResourceTemplate", "template", "interfaceProfiles"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOFabricResourcePhysicalInterfaceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "name", msoFabricResourcePhysicalInterfaceName),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "policy_group_type", "physical"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "description", ""),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "nodes.#", "1"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "uuid"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "template_id"),
					resource.TestCheckResourceAttrSet("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_policy_group_uuid"),
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
				PreConfig:    func() { fmt.Println("Test: Import missing Physical Interface") },
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
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[breakoutResourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", breakoutResourceName)
						}
						breakoutTemplateID = rs.Primary.Attributes["template_id"]
						breakoutUUID = rs.Primary.Attributes["uuid"]
						if breakoutTemplateID == "" || breakoutUUID == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", breakoutResourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Physical Interface with Breakout Mode after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), breakoutTemplateID, breakoutUUID, "fabricResourceTemplate", "template", "interfaceProfiles"); err != nil {
						t.Fatalf("delete %s out of band: %v", breakoutResourceName, err)
					}
				},
				Config: testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigCreate(),
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
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Physical Interface with Breakout Mode") },
				ResourceName: breakoutResourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[breakoutResourceName]
					if !ok || rs.Primary == nil {
						return "", fmt.Errorf("resource %s not found in state", breakoutResourceName)
					}
					separator := strings.LastIndex(rs.Primary.ID, "/")
					if separator < 0 || separator == len(rs.Primary.ID)-1 {
						return "", fmt.Errorf("resource %s has an invalid import ID", breakoutResourceName)
					}
					return rs.Primary.ID[:separator+1] + "tf_missing_oob", nil
				},
				ExpectError: regexp.MustCompile(`cannot import .*: resource not found`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Duplicate interface descriptions (error)") },
				Config:      testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingDuplicateInterfaceDescription(),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(fmt.Sprintf("interface profile %s_breakout_updated have more than one description for interface 1/2", msoFabricResourcePhysicalInterfaceName))),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Invalid interface in interface descriptions (error)") },
				Config:      testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingInvalidInterfaceDescription(),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(fmt.Sprintf("interface profile %s_breakout_updated doesn't have interface 1/3 which is used in description", msoFabricResourcePhysicalInterfaceName))),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_fabric_resource_policies_physical_interface", "fabricResourceTemplate", "template", "interfaceProfiles"),
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
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s"
        nodes                       = ["101"]
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceConfigUpdateAddingInterfaceDescriptions() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s"
		description                 = "Terraform test Physical Interface updated"
        nodes                       = ["101", "102"]
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
        interface_descriptions {
            interface   = "1/1"
            description = "Interface Description 1/1"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricPolicyTemplateInterfaceSettingName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceConfigUpdateAddingExtraInterfaceDescription() string {
	return fmt.Sprintf(`%[1]s
	resource "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s"
		description                 = "Terraform test Physical Interface updated"
        nodes                       = ["101", "102"]
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
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
        template_id                 = mso_template.%[4]s.id
        name                        = "%[2]s_updated"
        nodes                       = ["101"]
        interfaces                  = ["1/1","1/2"]
        interface_policy_group_uuid = mso_fabric_policies_interface_setting.%[3]s_physical.uuid
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

func testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingDuplicateInterfaceDescription() string {
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
        interface_descriptions {
            interface   = "1/2"
            description = "Interface Description 1/2 duplicate"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricResourceTemplateName)
}

func testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingInvalidInterfaceDescription() string {
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
        interface_descriptions {
            interface   = "1/3"
            description = "Interface Description 1/2 duplicate"
        }
	}`, fabricResourcePhysicalInterfacePreConfig, msoFabricResourcePhysicalInterfaceName, msoFabricResourceTemplateName)
}
