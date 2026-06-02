package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOServiceNodeTypeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Service Node Type Data Source - Not Found") },
				Config:      testAccMSOServiceNodeTypeDataSourceNotFound(),
				ExpectError: regexp.MustCompile(`Unable to find service node type`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Service Node Type Data Source") },
				Config:    testAccMSOServiceNodeTypeDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mso_service_node_type.test", "name", msoServiceNodeTypeName),
					resource.TestCheckResourceAttr("data.mso_service_node_type.test", "display_name", msoServiceNodeTypeName+" display"),
					resource.TestCheckResourceAttrPair("data.mso_service_node_type.test", "id", "mso_service_node_type.test", "id"),
				),
			},
		},
	})
}

func testAccMSOServiceNodeTypeDataSourceNotFound() string {
	return fmt.Sprintf(`
resource "mso_service_node_type" "test" {
  name         = "%[1]s"
  display_name = "%[1]s display"
}

data "mso_service_node_type" "test" {
  name = "non_existing_service_node_type"
}
`, msoServiceNodeTypeName)
}

func testAccMSOServiceNodeTypeDataSourceConfig() string {
	return fmt.Sprintf(`
resource "mso_service_node_type" "test" {
  name         = "%[1]s"
  display_name = "%[1]s display"
}

data "mso_service_node_type" "test" {
  name = mso_service_node_type.test.name
}
`, msoServiceNodeTypeName)
}
