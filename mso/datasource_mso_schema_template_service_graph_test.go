package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSOSchemaTemplateServiceGraphDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Schema Template Service Graph Data Source - Not Found") },
				Config:      testAccMSOSchemaTemplateServiceGraphDataSourceNotFound(),
				ExpectError: regexp.MustCompile(`unable to find service graph`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Schema Template Service Graph Data Source") },
				Config:    testAccMSOSchemaTemplateServiceGraphDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "schema_id"),
					resource.TestCheckResourceAttr("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "description", "Terraform test service graph"),
					resource.TestCheckResourceAttr("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.#", "3"),
					resource.TestCheckResourceAttr("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.0.type", "firewall"),
					resource.TestCheckResourceAttr("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.1.type", "load-balancer"),
					resource.TestCheckResourceAttr("data.mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.2.type", "other"),
				),
			},
		},
	})
}

func testAccMSOSchemaTemplateServiceGraphDataSourceNotFound() string {
	return fmt.Sprintf(`%s

data "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema_template_service_graph.%[2]s.schema_id
  template_name      = mso_schema_template_service_graph.%[2]s.template_name
  service_graph_name = "non_existing_service_graph"
}
`, testAccMSOSchemaTemplateServiceGraphConfigCreate(), msoSchemaTemplateServiceGraphName)
}

func testAccMSOSchemaTemplateServiceGraphDataSourceConfig() string {
	return fmt.Sprintf(`%s

data "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema_template_service_graph.%[2]s.schema_id
  template_name      = mso_schema_template_service_graph.%[2]s.template_name
  service_graph_name = mso_schema_template_service_graph.%[2]s.service_graph_name
}
`, testAccMSOSchemaTemplateServiceGraphConfigThreeNodes("Terraform test service graph"), msoSchemaTemplateServiceGraphName)
}
