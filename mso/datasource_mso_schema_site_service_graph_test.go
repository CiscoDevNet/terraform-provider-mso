package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaSiteServiceGraphDatasource(t *testing.T) {
	resourceRef := "mso_schema_site_service_graph." + msoSchemaTemplateServiceGraphName
	datasourceRef := "data.mso_schema_site_service_graph." + msoSchemaTemplateServiceGraphName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAPICPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaSiteServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Lookup non-existing site service graph (expect error)") },
				Config:      testAccMSOSchemaSiteServiceGraphDatasourceNotFoundConfig(),
				ExpectError: regexp.MustCompile(`Unable to find site service graph`),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Read existing site service graph and verify attributes match resource")
				},
				Config: testAccMSOSchemaSiteServiceGraphDatasourceReadConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceRef, "schema_id", resourceRef, "schema_id"),
					resource.TestCheckResourceAttrPair(datasourceRef, "template_name", resourceRef, "template_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "site_id", resourceRef, "site_id"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_graph_name", resourceRef, "service_graph_name"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_node.0.device_dn", resourceRef, "service_node.0.device_dn"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_node.0.consumer_connector_type", resourceRef, "service_node.0.consumer_connector_type"),
					resource.TestCheckResourceAttrPair(datasourceRef, "service_node.0.provider_connector_type", resourceRef, "service_node.0.provider_connector_type"),
				),
			},
		},
	})
}

// testAccMSOSchemaSiteServiceGraphDatasourcePrereqConfig returns the prerequisite
// HCL for the datasource test: schema + site association + template service graph.
// This is a self-contained copy of the resource test prereqs so the datasource
// test file is independent.
func testAccMSOSchemaSiteServiceGraphDatasourcePrereqConfig() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"

  service_node {
    type = "firewall"
  }
}
`, testSchemaWithSingleSiteAssociationConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName)
}

// testAccMSOSchemaSiteServiceGraphDatasourceResourceConfig creates the
// mso_schema_site_service_graph resource without any datasource block.
// Used as the %[1]s prerequisite in both datasource config functions below.
func testAccMSOSchemaSiteServiceGraphDatasourceResourceConfig() string {
	return fmt.Sprintf(`%[1]s
resource "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  site_id            = mso_schema_site.%[5]s.site_id
  service_graph_name = "%[2]s"

  service_node {
    device_dn = "%[6]s"
  }
}
`, testAccMSOSchemaSiteServiceGraphDatasourcePrereqConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName, msoSchemaSiteResourceLabel1, msoSchemaSiteServiceGraphDeviceDn)
}

// testAccMSOSchemaSiteServiceGraphDatasourceNotFoundConfig creates the resource
// and then queries the datasource with a non-existing service_graph_name.
// The datasource inputs reference the resource's own output attributes to
// guarantee the resource is applied before the datasource is evaluated.
func testAccMSOSchemaSiteServiceGraphDatasourceNotFoundConfig() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_site_service_graph" "not_found" {
  schema_id          = mso_schema_site_service_graph.%[2]s.schema_id
  template_name      = mso_schema_site_service_graph.%[2]s.template_name
  site_id            = mso_schema_site_service_graph.%[2]s.site_id
  service_graph_name = "non_existing_graph"
}
`, testAccMSOSchemaSiteServiceGraphDatasourceResourceConfig(), msoSchemaTemplateServiceGraphName)
}

// testAccMSOSchemaSiteServiceGraphDatasourceReadConfig creates the resource and
// a datasource that looks it up by referencing the resource's own output attributes.
func testAccMSOSchemaSiteServiceGraphDatasourceReadConfig() string {
	return fmt.Sprintf(`%[1]s
data "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema_site_service_graph.%[2]s.schema_id
  template_name      = mso_schema_site_service_graph.%[2]s.template_name
  site_id            = mso_schema_site_service_graph.%[2]s.site_id
  service_graph_name = mso_schema_site_service_graph.%[2]s.service_graph_name
}
`, testAccMSOSchemaSiteServiceGraphDatasourceResourceConfig(), msoSchemaTemplateServiceGraphName)
}
