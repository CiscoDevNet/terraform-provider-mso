package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateContractServiceGraphDatasource(t *testing.T) {
	resourceRef := "mso_schema_template_contract_service_graph." + msoSchemaTemplateContractName
	datasourceRef := "data.mso_schema_template_contract_service_graph." + msoSchemaTemplateContractName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateContractServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				// Contract exists but no service graph relationship has been created yet.
				// This directly tests the new not-found handling in the datasource.
				PreConfig:   func() { fmt.Println("DataSource: Lookup contract with no service graph relationship (expect error)") },
				Config:      testAccMSOSchemaTemplateContractServiceGraphDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`service graph relationship not found`),
			},
			{
				PreConfig: func() { fmt.Println("DataSource: Read existing contract service graph") },
				Config:    testAccMSOSchemaTemplateContractServiceGraphDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceRef, "schema_id", resourceRef, "schema_id"),
					resource.TestCheckResourceAttrPair(datasourceRef, "template_name", resourceRef, "template_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "contract_name", resourceRef, "contract_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_graph_name", resourceRef, "service_graph_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_graph_schema_id", resourceRef, "service_graph_schema_id"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_graph_template_name", resourceRef, "service_graph_template_name"),
					resource.TestCheckResourceAttr(datasourceRef, "node_relationship.#", "1"),
					resource.TestCheckResourceAttr(datasourceRef, "node_relationship.0.provider_connector_bd_name", msoSchemaTemplateBdName),
					resource.TestCheckResourceAttr(datasourceRef, "node_relationship.0.consumer_connector_bd_name", msoSchemaTemplateBdName),
				),
			},
		},
	})
}

// testAccMSOSchemaTemplateContractServiceGraphDatasourceNotFoundConfig creates
// the full resource (including the service graph relationship), then queries a
// non-existent contract name. The datasource references the resource's schema_id
// and template_name to create an implicit dependency, ensuring the relationship
// is created before the datasource runs. This prevents the datasource from
// running concurrently and failing before the relationship resource is applied.
func testAccMSOSchemaTemplateContractServiceGraphDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_template_contract_service_graph" "not_found" {
  schema_id     = mso_schema_template_contract_service_graph.%[2]s.schema_id
  template_name = mso_schema_template_contract_service_graph.%[2]s.template_name
  contract_name = "nonexistent-contract"
}
`, testAccMSOSchemaTemplateContractServiceGraphConfigCreate(),
		msoSchemaTemplateContractName, // %[2]s
	)
}

func testAccMSOSchemaTemplateContractServiceGraphDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_template_contract_service_graph" "%[2]s" {
  schema_id     = mso_schema_template_contract_service_graph.%[2]s.schema_id
  template_name = mso_schema_template_contract_service_graph.%[2]s.template_name
  contract_name = mso_schema_template_contract_service_graph.%[2]s.contract_name
}
`, testAccMSOSchemaTemplateContractServiceGraphConfigCreate(),
		msoSchemaTemplateContractName, // %[2]s
	)
}
