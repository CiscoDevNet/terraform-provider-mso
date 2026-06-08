package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOSchemaSiteServiceGraphDatasource(t *testing.T) {
	resourceRef := "mso_schema_site_service_graph." + msoSchemaTemplateServiceGraphName
	datasourceRef := "data.mso_schema_site_service_graph." + msoSchemaTemplateServiceGraphName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Read Site Service Graph via datasource and verify attributes match resource")
				},
				Config: testAccMSOSchemaSiteServiceGraphDatasourceConfig(),
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
	return fmt.Sprintf(`%s
resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"

  service_node {
    type = "firewall"
  }
}
`, testSchemaWithAnsibleTestTenantAndSingleSiteConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaSiteServiceGraphDatasourceConfig() string {
	return fmt.Sprintf(`%s
resource "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  site_id            = mso_schema_site.%[5]s.site_id
  service_graph_name = "%[2]s"

  service_node {
    device_dn = "%[6]s"
  }
}

data "mso_schema_site_service_graph" "%[2]s" {
  schema_id          = mso_schema_site_service_graph.%[2]s.schema_id
  template_name      = mso_schema_site_service_graph.%[2]s.template_name
  site_id            = mso_schema_site_service_graph.%[2]s.site_id
  service_graph_name = mso_schema_site_service_graph.%[2]s.service_graph_name
}
`, testAccMSOSchemaSiteServiceGraphDatasourcePrereqConfig(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName, msoSchemaSiteResourceLabel1, msoSchemaSiteServiceGraphDeviceDn)
}
