package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoSchemaTemplateContractServiceChainingSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateContractServiceChainingSchemaId string

func TestAccMSOSchemaTemplateContractServiceChainingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractServiceChainingDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Service Chaining with two nodes") },
				Config:    testAccMSOSchemaTemplateContractServiceChainingConfigCreateTwoNodes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract_service_chaining.chain1", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "node_filter", "allow-all"),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "service_nodes.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node1",
						"device_type":                         "loadBalancer",
						"consumer_connector.0.interface_name": "interface2",
						"consumer_connector.0.is_redirect":    "false",
						"provider_connector.0.interface_name": "interface1",
						"provider_connector.0.is_redirect":    "false",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node2",
						"device_type":                         "firewall",
						"consumer_connector.0.interface_name": "interface",
						"consumer_connector.0.is_redirect":    "false",
						"provider_connector.0.interface_name": "interface",
						"provider_connector.0.is_redirect":    "false",
					}),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_contract_service_chaining.chain1"]
						if !ok {
							return fmt.Errorf("Service Chaining resource not found in state")
						}
						msoSchemaTemplateContractServiceChainingSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Chaining to one node") },
				Config:    testAccMSOSchemaTemplateContractServiceChainingConfigUpdateOneNode(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "service_nodes.#", "1"),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node1",
						"device_type":                         "loadBalancer",
						"consumer_connector.0.interface_name": "interface2",
						"provider_connector.0.interface_name": "interface1",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Service Chaining to two nodes reordered") },
				Config:    testAccMSOSchemaTemplateContractServiceChainingConfigUpdateTwoNodesReordered(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "service_nodes.#", "2"),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node1",
						"device_type":                         "firewall",
						"consumer_connector.0.interface_name": "interface",
						"provider_connector.0.interface_name": "interface",
					}),
					CustomTestCheckTypeSetElemAttrs("mso_schema_template_contract_service_chaining.chain1", "service_nodes", map[string]string{
						"name":                                "node2",
						"device_type":                         "loadBalancer",
						"consumer_connector.0.interface_name": "interface2",
						"provider_connector.0.interface_name": "interface1",
					}),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Service Chaining") },
				ResourceName:      "mso_schema_template_contract_service_chaining.chain1",
				ImportState:       true,
				ImportStateIdFunc: testAccMSOSchemaTemplateContractServiceChainingImportStateIdFunc,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Service Chaining after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					path := fmt.Sprintf("/templates/%s/contracts/%s/serviceChaining", msoSchemaTemplateName, msoSchemaTemplateContractName)
					_, err := msoClient.PatchbyID(fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateContractServiceChainingSchemaId), models.GetRemovePatchPayload(path))
					if err != nil {
						t.Fatalf("Failed to manually delete Service Chaining: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateContractServiceChainingConfigUpdateTwoNodesReordered(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_contract_service_chaining.chain1", "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr("mso_schema_template_contract_service_chaining.chain1", "service_nodes.#", "2"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateContractServiceChainingImportStateIdFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["mso_schema_template_contract_service_chaining.chain1"]
	if !ok {
		return "", fmt.Errorf("Service Chaining resource not found in state")
	}
	return fmt.Sprintf("%s/templates/%s/contracts/%s/serviceChaining",
		rs.Primary.Attributes["schema_id"],
		rs.Primary.Attributes["template_name"],
		rs.Primary.Attributes["contract_name"],
	), nil
}

func testAccCheckMSOSchemaTemplateContractServiceChainingDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "mso_schema_template_contract_service_chaining" {
			schemaID := rs.Primary.Attributes["schema_id"]
			templateName := rs.Primary.Attributes["template_name"]
			contractName := rs.Primary.Attributes["contract_name"]

			con, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaID))
			if err != nil {
				return nil
			}

			templates := con.Search("templates").Data()
			if templates == nil {
				return nil
			}
			for _, template := range templates.([]interface{}) {
				templateMap := template.(map[string]interface{})
				if templateMap["name"].(string) != templateName {
					continue
				}
				contracts, ok := templateMap["contracts"].([]interface{})
				if !ok {
					return nil
				}
				for _, contract := range contracts {
					contractMap := contract.(map[string]interface{})
					if contractMap["name"].(string) == contractName {
						if _, exists := contractMap["serviceChaining"]; exists {
							return fmt.Errorf("Service Chaining record still exists for contract %s", contractName)
						}
					}
				}
			}
		}
	}
	return nil
}

func testAccMSOSchemaTemplateContractServiceChainingPrerequisiteConfig() string {
	return fmt.Sprintf(`%[1]s%[2]s
resource "mso_template" "%[3]s" {
	template_name = "%[3]s"
	template_type = "service_device"
	tenant_id     = mso_tenant.%[4]s.id
}

resource "mso_schema" "%[5]s" {
	name = "%[5]s"
	template {
		name          = "%[6]s"
		display_name  = "%[6]s"
		tenant_id     = mso_tenant.%[4]s.id
		template_type = "aci_multi_site"
	}
}

resource "mso_schema_template_vrf" "%[7]s" {
	schema_id    = mso_schema.%[5]s.id
	template     = "%[6]s"
	name         = "%[7]s"
	display_name = "%[7]s"
}

resource "mso_schema_template_bd" "%[8]s" {
	schema_id     = mso_schema.%[5]s.id
	template_name = "%[6]s"
	name          = "%[8]s"
	display_name  = "%[8]s"
	vrf_name      = mso_schema_template_vrf.%[7]s.name
	arp_flooding  = true
}

resource "mso_schema_template_bd" "%[9]s" {
	schema_id     = mso_schema.%[5]s.id
	template_name = "%[6]s"
	name          = "%[9]s"
	display_name  = "%[9]s"
	vrf_name      = mso_schema_template_vrf.%[7]s.name
	arp_flooding  = true
}

resource "mso_service_device_cluster" "%[10]s" {
	template_id = mso_template.%[3]s.id
	name        = "%[10]s"
	device_mode = "layer3"
	device_type = "load_balancer"

	interface_properties {
		name    = "interface1"
		bd_uuid = mso_schema_template_bd.%[8]s.uuid
	}

	interface_properties {
		name    = "interface2"
		bd_uuid = mso_schema_template_bd.%[9]s.uuid
	}
}

resource "mso_service_device_cluster" "%[11]s" {
	template_id = mso_service_device_cluster.%[10]s.template_id
	name        = "%[11]s"
	device_mode = "layer3"
	device_type = "firewall"

	interface_properties {
		name    = "interface"
		bd_uuid = mso_schema_template_bd.%[8]s.uuid
	}
}

resource "mso_schema_template_filter_entry" "%[12]s" {
	schema_id          = mso_schema.%[5]s.id
	template_name      = "%[6]s"
	name               = "%[12]s"
	display_name       = "%[12]s"
	entry_name         = "%[12]s_entry"
	entry_display_name = "%[12]s_entry"
}

resource "mso_schema_template_contract" "%[13]s" {
	schema_id     = mso_schema.%[5]s.id
	template_name = "%[6]s"
	contract_name = "%[13]s"
	display_name  = "%[13]s"
	filter_type   = "bothWay"
	scope         = "context"
	filter_relationship {
		filter_name = mso_schema_template_filter_entry.%[12]s.name
		filter_type = "bothWay"
	}
}

resource "mso_schema_template_external_epg" "%[14]s" {
	schema_id                  = mso_schema.%[5]s.id
	template_name              = "%[6]s"
	external_epg_name          = "%[14]s"
	display_name               = "%[14]s"
	vrf_name                   = mso_schema_template_vrf.%[7]s.name
	vrf_schema_id              = mso_schema_template_vrf.%[7]s.schema_id
	vrf_template_name          = "%[6]s"
	external_epg_type          = "on-premise"
	include_in_preferred_group = false
}

resource "mso_schema_template_external_epg_contract" "%[14]s_provider" {
	schema_id              = mso_schema.%[5]s.id
	template_name          = "%[6]s"
	contract_name          = mso_schema_template_contract.%[13]s.contract_name
	external_epg_name      = mso_schema_template_external_epg.%[14]s.external_epg_name
	relationship_type      = "provider"
	contract_schema_id     = mso_schema_template_contract.%[13]s.schema_id
	contract_template_name = mso_schema_template_contract.%[13]s.template_name
}
`,
		testSiteConfigAnsibleTest(),   // %[1]s
		testTenantConfig(),            // %[2]s
		msoServiceDeviceTemplateName,  // %[3]s
		msoTenantName,                 // %[4]s
		msoSchemaName,                 // %[5]s
		msoSchemaTemplateName,         // %[6]s
		msoSchemaTemplateVrfName,      // %[7]s
		msoSchemaTemplateBdName,       // %[8]s
		msoSchemaTemplateBdName2,      // %[9]s
		msoServiceDeviceClusterLbName, // %[10]s
		msoServiceDeviceClusterFwName, // %[11]s
		msoSchemaTemplateFilterName,   // %[12]s
		msoSchemaTemplateContractName, // %[13]s
		msoSchemaTemplateExtEpgName,   // %[14]s
	)
}

func testAccMSOSchemaTemplateContractServiceChainingConfigCreateTwoNodes() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_contract_service_chaining" "chain1" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	contract_name = mso_schema_template_contract.%[4]s.contract_name

	node_filter = "allow-all"

	service_nodes {
		name        = "node1"
		device_type = "loadBalancer"
		device_ref  = mso_service_device_cluster.%[5]s.uuid

		consumer_connector {
			interface_name = "interface2"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface1"
			is_redirect    = false
		}
	}

	service_nodes {
		name        = "node2"
		device_type = "firewall"
		device_ref  = mso_service_device_cluster.%[6]s.uuid

		consumer_connector {
			interface_name = "interface"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface"
			is_redirect    = false
		}
	}
}
`,
		testAccMSOSchemaTemplateContractServiceChainingPrerequisiteConfig(),
		msoSchemaName,
		msoSchemaTemplateName,
		msoSchemaTemplateContractName,
		msoServiceDeviceClusterLbName,
		msoServiceDeviceClusterFwName,
	)
}

func testAccMSOSchemaTemplateContractServiceChainingConfigUpdateOneNode() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_contract_service_chaining" "chain1" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	contract_name = mso_schema_template_contract.%[4]s.contract_name

	node_filter = "allow-all"

	service_nodes {
		name        = "node1"
		device_type = "loadBalancer"
		device_ref  = mso_service_device_cluster.%[5]s.uuid

		consumer_connector {
			interface_name = "interface2"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface1"
			is_redirect    = false
		}
	}
}
`,
		testAccMSOSchemaTemplateContractServiceChainingPrerequisiteConfig(),
		msoSchemaName,
		msoSchemaTemplateName,
		msoSchemaTemplateContractName,
		msoServiceDeviceClusterLbName,
	)
}

func testAccMSOSchemaTemplateContractServiceChainingConfigUpdateTwoNodesReordered() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_contract_service_chaining" "chain1" {
	schema_id     = mso_schema.%[2]s.id
	template_name = "%[3]s"
	contract_name = mso_schema_template_contract.%[4]s.contract_name

	node_filter = "allow-all"

	service_nodes {
		name        = "node1"
		device_type = "firewall"
		device_ref  = mso_service_device_cluster.%[6]s.uuid

		consumer_connector {
			interface_name = "interface"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface"
			is_redirect    = false
		}
	}

	service_nodes {
		name        = "node2"
		device_type = "loadBalancer"
		device_ref  = mso_service_device_cluster.%[5]s.uuid

		consumer_connector {
			interface_name = "interface2"
			is_redirect    = false
		}

		provider_connector {
			interface_name = "interface1"
			is_redirect    = false
		}
	}
}
`,
		testAccMSOSchemaTemplateContractServiceChainingPrerequisiteConfig(),
		msoSchemaName,
		msoSchemaTemplateName,
		msoSchemaTemplateContractName,
		msoServiceDeviceClusterLbName,
		msoServiceDeviceClusterFwName,
	)
}
