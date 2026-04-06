package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOFabricResourcePortChannelInterfaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Port Channel Interface Data Source - With Interface Descriptions") },
				Config:    testAccMSOFabricResourcePortChannelInterfaceDataSourceWithInterfaceDescriptions(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "name", msoFabricResourcePortChannelInterfaceName),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "description", "Terraform test Port Channel Interface updated"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "node", "101"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interfaces.#", "2"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "uuid"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "template_id"),
					resource.TestCheckResourceAttrSet("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_policy_group_uuid"),
					resource.TestCheckResourceAttr("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions.#", "2"),
					CustomTestCheckTypeSetElemAttrs("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions",
						map[string]string{
							"interface":   "1/1",
							"description": "Interface Description 1/1",
						},
					),
					CustomTestCheckTypeSetElemAttrs("data.mso_fabric_resource_policies_port_channel_interface."+msoFabricResourcePortChannelInterfaceName, "interface_descriptions",
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

func testAccMSOFabricResourcePortChannelInterfaceDataSourceWithInterfaceDescriptions() string {
	return fmt.Sprintf(`%[1]s
    data "mso_fabric_resource_policies_port_channel_interface" "%[2]s" {
        template_id = mso_template.%[3]s.id
        name        = mso_fabric_resource_policies_port_channel_interface.%[2]s.name
    }`,
		testAccMSOFabricResourcePortChannelInterfaceConfigUpdateAddingExtraInterfaceDescription(),
		msoFabricResourcePortChannelInterfaceName,
		msoFabricResourceTemplateName,
	)
}
