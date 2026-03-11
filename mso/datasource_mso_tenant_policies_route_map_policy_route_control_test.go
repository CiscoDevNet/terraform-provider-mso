package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSORouteMapPolicyDataSource(t *testing.T) {
	resourceName := "mso_route_map_policy.test_policy"
	dataSourceName := "data.mso_route_map_policy.test_ds"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Task: Create Resource and Query via Data Source") },
				Config:    testAccMSORouteMapPolicyDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "template_id", resourceName, "template_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", resourceName, "uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "test_ds_policy"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "Data source test description"),
				),
			},
		},
	})
}

func testAccMSORouteMapPolicyDataSourceConfig() string {
	return fmt.Sprintf(`%s
resource "mso_route_map_policy" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_ds_policy"
  description = "Data source test description"
}

data "mso_route_map_policy" "test_ds" {
  template_id = mso_route_map_policy.test_policy.template_id
  name        = mso_route_map_policy.test_policy.name
}
`, testAccMSOTemplateResourceTenantConfig())
}
