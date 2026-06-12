package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccMSOSchemaSiteContractServiceGraphDatasource tests the datasource for
// mso_schema_site_contract_service_graph:
//
//  1. Not found — all prerequisites are created (including the template-level
//     contract service graph and the site service graph) but the
//     mso_schema_site_contract_service_graph resource is absent. NDO mirrors
//     template contracts into the site's contracts list, so the contract entry
//     exists in sites[] but has no serviceGraphRelationship, which causes the
//     datasource to return "No service graph found".
//  2. Read — the resource is created and the datasource reads it back, checking
//     all attributes via AttrPair.
func TestAccMSOSchemaSiteContractServiceGraphDatasource(t *testing.T) {
	resourceRef := "mso_schema_site_contract_service_graph." + msoSchemaTemplateContractName
	datasourceRef := "data.mso_schema_site_contract_service_graph." + msoSchemaTemplateContractName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteContractServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("DataSource: Lookup site contract with no service graph relationship (expect error)")
				},
				Config:      testAccMSOSchemaSiteContractServiceGraphDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`No service graph found`),
			},
			{
				PreConfig: func() {
					fmt.Println("DataSource: Read existing site contract service graph and verify attributes match resource")
				},
				Config: testAccMSOSchemaSiteContractServiceGraphDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceRef, "schema_id", resourceRef, "schema_id"),
					resource.TestCheckResourceAttrPair(datasourceRef, "template_name", resourceRef, "template_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "site_id", resourceRef, "site_id"),
					resource.TestCheckResourceAttrPair(datasourceRef, "contract_name", resourceRef, "contract_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_graph_name", resourceRef, "service_graph_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_graph_schema_id", resourceRef, "service_graph_schema_id"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_graph_template_name", resourceRef, "service_graph_template_name"),
					resource.TestCheckResourceAttr(datasourceRef, "node_relationship.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceRef, "node_relationship.0.provider_connector_cluster_interface", resourceRef, "node_relationship.0.provider_connector_cluster_interface"),
					resource.TestCheckResourceAttrPair(datasourceRef, "node_relationship.0.provider_connector_redirect_policy_tenant", resourceRef, "node_relationship.0.provider_connector_redirect_policy_tenant"),
					resource.TestCheckResourceAttrPair(datasourceRef, "node_relationship.0.provider_connector_redirect_policy", resourceRef, "node_relationship.0.provider_connector_redirect_policy"),
					resource.TestCheckResourceAttrPair(datasourceRef, "node_relationship.0.consumer_connector_cluster_interface", resourceRef, "node_relationship.0.consumer_connector_cluster_interface"),
					resource.TestCheckResourceAttrPair(datasourceRef, "node_relationship.0.consumer_connector_redirect_policy_tenant", resourceRef, "node_relationship.0.consumer_connector_redirect_policy_tenant"),
					resource.TestCheckResourceAttrPair(datasourceRef, "node_relationship.0.consumer_connector_redirect_policy", resourceRef, "node_relationship.0.consumer_connector_redirect_policy"),
					resource.TestCheckResourceAttrPair(datasourceRef, "node_relationship.0.consumer_subnet_ips.#", resourceRef, "node_relationship.0.consumer_subnet_ips.#"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteContractServiceGraphDatasourcePrereqConfig builds the
// self-contained prerequisite stack for the datasource test. It is identical
// to testAccMSOSchemaSiteContractServiceGraphPrereqConfig but lives here so the
// datasource test file is independent of the resource test file.
func testAccMSOSchemaSiteContractServiceGraphDatasourcePrereqConfig() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"

  service_node {
    type = "firewall"
  }
}

%[5]s%[6]s%[7]s%[8]s%[9]s

resource "mso_schema_template_contract_service_graph" "%[10]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  contract_name      = mso_schema_template_contract.%[10]s.contract_name
  service_graph_name = mso_schema_template_service_graph.%[2]s.service_graph_name

  node_relationship {
    provider_connector_bd_name = mso_schema_template_bd.%[11]s.name
    consumer_connector_bd_name = mso_schema_template_bd.%[11]s.name
  }
}

resource "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  site_id            = mso_schema_site.%[12]s.site_id
  service_graph_name = "%[2]s"

  service_node {
    device_dn = "%[13]s"
  }

  depends_on = [mso_schema_template_service_graph.%[2]s]
}
`,
		testSchemaWithAnsibleTestTenantAndSingleSiteConfig(), // %[1]s
		msoSchemaTemplateServiceGraphName,                    // %[2]s
		msoSchemaName,                                        // %[3]s
		msoSchemaTemplateName,                                // %[4]s
		testSchemaTemplateVrfConfig(),                        // %[5]s
		testSchemaTemplateBdConfig(),                         // %[6]s
		testSchemaTemplateFilterEntryConfig(),                // %[7]s
		testSchemaTemplateContractConfig(),                   // %[8]s
		testSchemaTemplateVrfContractConfig(),                // %[9]s
		msoSchemaTemplateContractName,                        // %[10]s
		msoSchemaTemplateBdName,                              // %[11]s
		msoSchemaSiteResourceLabel1,                          // %[12]s
		msoSchemaSiteContractServiceGraphDeviceDn,            // %[13]s
	)
}

// testAccMSOSchemaSiteContractServiceGraphDatasourceResourceConfig creates the
// mso_schema_site_contract_service_graph resource without any datasource block.
// Used as the prerequisite in both datasource config functions below.
func testAccMSOSchemaSiteContractServiceGraphDatasourceResourceConfig() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_site_contract_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  site_id            = mso_schema_site.%[4]s.site_id
  template_name      = "%[5]s"
  contract_name      = mso_schema_template_contract.%[2]s.contract_name
  service_graph_name = mso_schema_site_service_graph.%[6]s.service_graph_name

  node_relationship {
    provider_connector_cluster_interface = "%[7]s"
    consumer_connector_cluster_interface = "%[8]s"
  }
}
`,
		testAccMSOSchemaSiteContractServiceGraphDatasourcePrereqConfig(), // %[1]s
		msoSchemaTemplateContractName,                                    // %[2]s
		msoSchemaName,                                                    // %[3]s
		msoSchemaSiteResourceLabel1,                                      // %[4]s
		msoSchemaTemplateName,                                            // %[5]s
		msoSchemaTemplateServiceGraphName,                                // %[6]s
		msoSchemaSiteContractServiceGraphProviderClusterInterface,        // %[7]s
		msoSchemaSiteContractServiceGraphConsumerClusterInterface,        // %[8]s
	)
}

// testAccMSOSchemaSiteContractServiceGraphDatasourceNotFoundConfig creates the
// mso_schema_site_contract_service_graph resource and then queries the datasource
// with a non-existent contract name. The datasource inputs reference the resource's
// own output attributes to force an implicit dependency, ensuring the resource is
// fully applied before the datasource evaluates. The non-existent contract_name
// causes setSiteContractServiceGraphAttrs to fall through its loop and return
// "No service graph found".
func testAccMSOSchemaSiteContractServiceGraphDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_site_contract_service_graph" "not_found" {
  schema_id     = mso_schema_site_contract_service_graph.%[2]s.schema_id
  template_name = mso_schema_site_contract_service_graph.%[2]s.template_name
  site_id       = mso_schema_site_contract_service_graph.%[2]s.site_id
  contract_name = "non_existing_contract"
}
`,
		testAccMSOSchemaSiteContractServiceGraphDatasourceResourceConfig(), // %[1]s
		msoSchemaTemplateContractName,                                      // %[2]s
	)
}

// testAccMSOSchemaSiteContractServiceGraphDatasourceReadConfig creates the
// resource and a datasource that reads it back by referencing the resource's
// own output attributes, creating an implicit dependency so the datasource
// evaluates only after the resource has been applied.
func testAccMSOSchemaSiteContractServiceGraphDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_site_contract_service_graph" "%[2]s" {
  schema_id     = mso_schema_site_contract_service_graph.%[2]s.schema_id
  template_name = mso_schema_site_contract_service_graph.%[2]s.template_name
  site_id       = mso_schema_site_contract_service_graph.%[2]s.site_id
  contract_name = mso_schema_site_contract_service_graph.%[2]s.contract_name
}
`,
		testAccMSOSchemaSiteContractServiceGraphDatasourceResourceConfig(), // %[1]s
		msoSchemaTemplateContractName,                                      // %[2]s
	)
}
