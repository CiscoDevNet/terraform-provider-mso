package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSORouteMapPolicyDataSource(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_route_control.test_policy"
	dataSourceName := "data.mso_tenant_policies_route_map_policy_route_control.test_ds"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Missing Route Map Policy data source returns an error") },
				Config:      testAccMSORouteMapPolicyDataSourceMissing(),
				ExpectError: regexp.MustCompile(`Policy name missing_route_map_policy not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Route Map Policy Data Source") },
				Config:    testAccMSORouteMapPolicyDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "template_id", resourceName, "template_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", resourceName, "uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_type", resourceName, "match_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "uuid"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "test_policy"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "Initial description"),
					resource.TestCheckResourceAttr(dataSourceName, "match_type", "prefix_and_routing_policy"),
				),
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithPathAttributesAndArguments("mso_tenant_policies_route_map_policy_route_control", "tenantPolicyTemplate", "template", "routeMapPolicies"),
	})
}

func testAccMSORouteMapPolicyDataSourceConfig() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_route_map_policy_route_control" "test_ds" {
  template_id = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  name        = mso_tenant_policies_route_map_policy_route_control.test_policy.name
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyDataSourceMissing() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_route_map_policy_route_control" "missing" {
  template_id = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  name        = "missing_route_map_policy"
}
`, testAccMSORouteMapPolicyConfigCreate())
}
