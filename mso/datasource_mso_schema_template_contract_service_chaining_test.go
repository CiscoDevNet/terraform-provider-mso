package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaTemplateContractServiceChainingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMSOSchemaTemplateContractServiceChainingDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_schema_template_contract_service_chaining.chain1_data", "contract_name", "contract1"),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract_service_chaining.chain1_data", "node_filter", "allow-all"),
					resource.TestCheckResourceAttr("data.mso_schema_template_contract_service_chaining.chain1_data", "service_nodes.#", "2"),
					CustomTestCheckTypeSetElemAttrs("data.mso_schema_template_contract_service_chaining.chain1_data", "service_nodes", map[string]string{
						"name":                                "node1",
						"device_type":                         "loadBalancer",
						"consumer_connector.0.interface_name": "interface2",
						"consumer_connector.0.is_redirect":    "false",
						"provider_connector.0.interface_name": "interface1",
						"provider_connector.0.is_redirect":    "false",
					}),
					CustomTestCheckTypeSetElemAttrs("data.mso_schema_template_contract_service_chaining.chain1_data", "service_nodes", map[string]string{
						"name":                                "node2",
						"device_type":                         "firewall",
						"consumer_connector.0.interface_name": "interface",
						"consumer_connector.0.is_redirect":    "false",
						"provider_connector.0.interface_name": "interface",
						"provider_connector.0.is_redirect":    "false",
					}),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateContractServiceChainingDataSourceConfig() string {
	return fmt.Sprintf(`%s

    data "mso_schema_template_contract_service_chaining" "chain1_data" {
        schema_id     = mso_schema_template_contract_service_chaining.chain1.schema_id
        template_name = mso_schema_template_contract_service_chaining.chain1.template_name
        contract_name = mso_schema_template_contract_service_chaining.chain1.contract_name
    }
`, testAccMSOSchemaTemplateContractServiceChainingConfigCreateTwoNodes())
}
