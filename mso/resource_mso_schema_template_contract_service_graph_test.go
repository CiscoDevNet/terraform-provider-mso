package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoSchemaTemplateContractServiceGraphSchemaId is set during the first test
// step's Check to capture the dynamic schema ID for use in the manual deletion
// PreConfig step.
var msoSchemaTemplateContractServiceGraphSchemaId string

func TestAccMSOSchemaTemplateContractServiceGraphResource(t *testing.T) {
	resourceRef := "mso_schema_template_contract_service_graph." + msoSchemaTemplateContractName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create contract service graph with a firewall node using BD1 for provider and consumer")
				},
				Config: testAccMSOSchemaTemplateContractServiceGraphConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceRef, "id"),
					resource.TestCheckResourceAttr(resourceRef, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr(resourceRef, "contract_name", msoSchemaTemplateContractName),
					resource.TestCheckResourceAttr(resourceRef, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_bd_name", msoSchemaTemplateBdName),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceRef]
						if !ok {
							return fmt.Errorf("contract service graph resource not found in state")
						}
						msoSchemaTemplateContractServiceGraphSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update contract service graph node_relationship to use BD2 for provider and consumer")
				},
				Config: testAccMSOSchemaTemplateContractServiceGraphConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_bd_name", msoSchemaTemplateBdName2),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_bd_name", msoSchemaTemplateBdName2),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import contract service graph") },
				ResourceName:      resourceRef,
				ImportState:       true,
				ImportStateIdFunc: testAccMSOSchemaTemplateContractServiceGraphImportStateId(resourceRef),
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate contract service graph after manual deletion of service graph relationship")
					msoClient := testAccProvider.Meta().(*client.Client)
					path := fmt.Sprintf("/templates/%s/contracts/%s/serviceGraphRelationship", msoSchemaTemplateName, msoSchemaTemplateContractName)
					_, err := msoClient.PatchbyID(
						fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateContractServiceGraphSchemaId),
						models.GetRemovePatchPayload(path),
					)
					if err != nil {
						t.Fatalf("Failed to manually delete service graph relationship: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateContractServiceGraphConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceRef, "id"),
					resource.TestCheckResourceAttr(resourceRef, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.#", "1"),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.provider_connector_bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(resourceRef, "node_relationship.0.consumer_connector_bd_name", msoSchemaTemplateBdName),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaTemplateContractServiceGraphDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_template_contract_service_graph" {
			continue
		}
		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", rs.Primary.Attributes["schema_id"]))
		if err != nil {
			return nil
		}
		templatesCount, err := cont.ArrayCount("templates")
		if err != nil {
			return nil
		}
		for i := 0; i < templatesCount; i++ {
			templateCont, err := cont.ArrayElement(i, "templates")
			if err != nil {
				continue
			}
			if models.StripQuotes(templateCont.S("name").String()) != rs.Primary.Attributes["template_name"] {
				continue
			}
			contractCount, err := templateCont.ArrayCount("contracts")
			if err != nil {
				continue
			}
			for j := 0; j < contractCount; j++ {
				contractCont, err := templateCont.ArrayElement(j, "contracts")
				if err != nil {
					continue
				}
				if models.StripQuotes(contractCont.S("name").String()) != rs.Primary.Attributes["contract_name"] {
					continue
				}
				if contractCont.Exists("serviceGraphRelationship") {
					return fmt.Errorf("mso_schema_template_contract_service_graph (%s) still exists", rs.Primary.ID)
				}
			}
		}
	}
	return nil
}

func testAccMSOSchemaTemplateContractServiceGraphImportStateId(resourceRef string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceRef]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceRef)
		}
		return rs.Primary.ID, nil
	}
}

// testAccMSOSchemaTemplateContractServiceGraphBaseConfig builds the shared
// prerequisites: site, tenant, schema, VRF, BD1, filter entry, contract, and
// VRF-as-provider association. The VRF-contract provider is required by NDO
// when a service graph is attached to a contract.
func testAccMSOSchemaTemplateContractServiceGraphBaseConfig() string {
	return fmt.Sprintf(`%s%s%s%s%s%s%s%s`,
		testSiteConfigAnsibleTest(),
		testTenantConfig(),
		testSchemaConfig(),
		testSchemaTemplateVrfConfig(),
		testSchemaTemplateBdConfig(),
		testSchemaTemplateFilterEntryConfig(),
		testSchemaTemplateContractConfig(),
		testSchemaTemplateVrfContractConfig(),
	)
}

// testAccMSOSchemaTemplateContractServiceGraphConfig builds the full config
// with an optional extra prerequisites block (e.g. BD2) and a specific BD
// name for the node_relationship connectors.
func testAccMSOSchemaTemplateContractServiceGraphConfig(extraPrereqs, bdName string) string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"
  service_node {
    type = "firewall"
  }
}

resource "mso_schema_template_contract_service_graph" "%[5]s" {
  schema_id          = mso_schema_template_vrf_contract.%[5]s_provider.schema_id
  template_name      = "%[4]s"
  contract_name      = mso_schema_template_contract.%[5]s.contract_name
  service_graph_name = mso_schema_template_service_graph.%[2]s.service_graph_name
  node_relationship {
    provider_connector_bd_name = mso_schema_template_bd.%[6]s.name
    consumer_connector_bd_name = mso_schema_template_bd.%[6]s.name
  }
}
`, testAccMSOSchemaTemplateContractServiceGraphBaseConfig()+extraPrereqs,
		msoSchemaTemplateServiceGraphName, // %[2]s
		msoSchemaName,                     // %[3]s
		msoSchemaTemplateName,             // %[4]s
		msoSchemaTemplateContractName,     // %[5]s
		bdName,                            // %[6]s
	)
}

func testAccMSOSchemaTemplateContractServiceGraphConfigCreate() string {
	return testAccMSOSchemaTemplateContractServiceGraphConfig("", msoSchemaTemplateBdName)
}

func testAccMSOSchemaTemplateContractServiceGraphConfigUpdate() string {
	return testAccMSOSchemaTemplateContractServiceGraphConfig(testSchemaTemplateBdConfig2(), msoSchemaTemplateBdName2)
}
