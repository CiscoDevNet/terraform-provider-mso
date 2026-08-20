package mso

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// msoSchemaTemplateServiceGraphSchemaId is set during the first test step's Check to capture the dynamic schema ID for use in the manual deletion PreConfig step.
var msoSchemaTemplateServiceGraphSchemaId string

func TestAccMSOSchemaTemplateServiceGraphResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMSOSchemaTemplateServiceGraphDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Test: Create Schema Template Service Graph with a single firewall node and description")
				},
				Config: testAccMSOSchemaTemplateServiceGraphConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "template_name", msoSchemaTemplateName),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "description", "Terraform test service graph"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.#", "1"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.0.type", "firewall"),
					// Capture the dynamic schema ID from state for use in the manual deletion PreConfig step
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName]
						if !ok {
							return fmt.Errorf("Service Graph resource not found in state")
						}
						msoSchemaTemplateServiceGraphSchemaId = rs.Primary.Attributes["schema_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Schema Template Service Graph description and add a second service node (load-balancer) and third custom service node")
				},
				Config: testAccMSOSchemaTemplateServiceGraphConfigUpdateThreeNodes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "description", "Terraform test service graph updated"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.#", "3"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.0.type", "firewall"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.1.type", "load-balancer"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.2.type", "other"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Update Schema Template Service Graph to remove description")
				},
				Config: testAccMSOSchemaTemplateServiceGraphConfigRemoveDescription(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "description", ""),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.#", "3"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Schema Template Service Graph") },
				ResourceName:      "mso_schema_template_service_graph." + msoSchemaTemplateServiceGraphName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Schema Template Service Graph after manual deletion from NDO")
					msoClient := testAccProvider.Meta().(*client.Client)
					path := fmt.Sprintf("/templates/%s/serviceGraphs/%s", msoSchemaTemplateName, msoSchemaTemplateServiceGraphName)
					_, err := msoClient.PatchbyID(
						fmt.Sprintf("api/v1/schemas/%s", msoSchemaTemplateServiceGraphSchemaId),
						models.GetRemovePatchPayload(path),
					)
					if err != nil {
						t.Fatalf("Failed to manually delete Service Graph: %v", err)
					}
				},
				Config: testAccMSOSchemaTemplateServiceGraphConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "schema_id"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_graph_name", msoSchemaTemplateServiceGraphName),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "description", "Terraform test service graph"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.#", "1"),
					resource.TestCheckResourceAttr("mso_schema_template_service_graph."+msoSchemaTemplateServiceGraphName, "service_node.0.type", "firewall"),
				),
			},
		},
	})
}

func testAccCheckMSOSchemaTemplateServiceGraphDestroy(s *terraform.State) error {
	msoClient := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mso_schema_template_service_graph" {
			continue
		}
		schemaId := rs.Primary.Attributes["schema_id"]
		cont, err := msoClient.GetViaURL(fmt.Sprintf("api/v1/schemas/%s", schemaId))
		if err != nil {
			return nil
		}
		_, _, err = getTemplateServiceGraphCont(cont, rs.Primary.Attributes["template_name"], rs.Primary.Attributes["service_graph_name"])
		if err == nil {
			return fmt.Errorf("Schema Template Service Graph still exists")
		}
	}
	return nil
}

func testAccMSOSchemaTemplateServiceGraphDependencies() string {
	return fmt.Sprintf(`%s%s%s`, testSiteConfigAnsibleTest(), testTenantConfig(), testSchemaConfig())
}

func testAccMSOSchemaTemplateServiceGraphConfigCreate() string {
	return fmt.Sprintf(`%s
resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"
  description        = "Terraform test service graph"

  service_node {
    type = "firewall"
  }
}
`, testAccMSOSchemaTemplateServiceGraphDependencies(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName)
}

func testAccMSOSchemaTemplateServiceGraphConfigUpdateThreeNodes() string {
	return testAccMSOSchemaTemplateServiceGraphConfigThreeNodes("Terraform test service graph updated")
}

func testAccMSOSchemaTemplateServiceGraphConfigRemoveDescription() string {
	return testAccMSOSchemaTemplateServiceGraphConfigThreeNodes("")
}

func testAccMSOSchemaTemplateServiceGraphConfigThreeNodes(description string) string {
	descLine := ""
	if description != "" {
		descLine = fmt.Sprintf("  description        = %q\n", description)
	}
	return fmt.Sprintf(`%s

resource "mso_schema_template_service_graph" "%[2]s" {
  schema_id          = mso_schema.%[3]s.id
  template_name      = "%[4]s"
  service_graph_name = "%[2]s"
%[5]s
  service_node {
    type = "firewall"
  }

  service_node {
    type = "load-balancer"
  }

  service_node {
    type = "other"
  }
}
`, testAccMSOSchemaTemplateServiceGraphDependencies(), msoSchemaTemplateServiceGraphName, msoSchemaName, msoSchemaTemplateName, descLine)
}
