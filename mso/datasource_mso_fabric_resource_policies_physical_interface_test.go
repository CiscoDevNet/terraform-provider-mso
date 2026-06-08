package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOFabricResourcePhysicalInterfaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Physical Interface Data Source - With Interface type physical") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceTypePhysicalDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "name", msoFabricResourcePhysicalInterfaceName),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "policy_group_type", "physical"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "description", "Terraform test Physical Interface updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "nodes.#", "2"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "template_id"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_policy_group_uuid"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions.#", "2"),
					CustomTestCheckTypeSetElemAttrs("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "Interface Description 1/2",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Physical Interface Data Source - Breakout Mode setup") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingExtraInterfaceDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "name", msoFabricResourcePhysicalInterfaceName+"_breakout_updated"),
					resource.TestCheckResourceAttr("mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "policy_group_type", "breakout"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Physical Interface Data Source - Breakout Mode") },
				Config:    testAccMSOFabricResourcePhysicalInterfaceTypeBreakoutDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "name", msoFabricResourcePhysicalInterfaceName+"_breakout_updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "policy_group_type", "breakout"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "template_id"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "breakout_mode", "4x100G"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "description", "Terraform test Physical Interface updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "nodes.#", "1"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interfaces.#", "2"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions.#", "1"),
					CustomTestCheckTypeSetElemAttrs("data.mso_fabric_resource_policies_physical_interface."+msoFabricResourcePhysicalInterfaceName+"_breakout", "interface_descriptions",
						map[string]string{
							"interface":   "1/2",
							"description": "Interface Description 1/2",
						},
					),
				),
			},
		},
	})
}

func testAccMSOFabricResourcePhysicalInterfaceTypePhysicalDataSource() string {
	return fmt.Sprintf(`%[1]s
    data "mso_fabric_resource_policies_physical_interface" "%[2]s" {
        template_id = mso_template.%[3]s.id
        name        = mso_fabric_resource_policies_physical_interface.%[2]s.name
    }`,
		testAccMSOFabricResourcePhysicalInterfaceConfigUpdateAddingExtraInterfaceDescription(),
		msoFabricResourcePhysicalInterfaceName,
		msoFabricResourceTemplateName,
	)
}

func testAccMSOFabricResourcePhysicalInterfaceTypeBreakoutDataSource() string {
	return fmt.Sprintf(`%[1]s
	data "mso_fabric_resource_policies_physical_interface" "%[2]s" {
		template_id = mso_template.%[3]s.id
		name        = mso_fabric_resource_policies_physical_interface.%[2]s.name
	}`,
		testAccMSOFabricResourcePhysicalInterfaceBreakoutModeConfigUpdateRemovingExtraInterfaceDescription(),
		msoFabricResourcePhysicalInterfaceName+"_breakout",
		msoFabricResourceTemplateName,
	)
}
