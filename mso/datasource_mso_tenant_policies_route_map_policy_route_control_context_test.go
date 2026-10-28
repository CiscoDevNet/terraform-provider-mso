package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMSORouteMapPolicyContextDataSource(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_route_control_context.test_context"
	dataSourceName := "data.mso_tenant_policies_route_map_policy_route_control_context.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccVersionCheck(t, "5.1")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { fmt.Println("Test: Read Missing Route Map Policy Context") },
				Config:      testAccMSORouteMapPolicyContextDataSourceConfigMissingContext(),
				ExpectError: regexp.MustCompile(`(?s)Route Map Policy Context.*tf_missing_context.*not found.*parent policy`),
			},
			{
				PreConfig:   func() { fmt.Println("Test: Read Route Map Policy Context with Missing Parent Policy") },
				Config:      testAccMSORouteMapPolicyContextDataSourceConfigMissingParent(),
				ExpectError: regexp.MustCompile(`Policy uuid 00000000-0000-0000-0000-000000000000 not found`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Read Route Map Policy Context") },
				Config:    testAccMSORouteMapPolicyContextDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "template_id", resourceName, "template_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "route_map_policy_uuid", resourceName, "route_map_policy_uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "action", resourceName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceName, "order", resourceName, "order"),
					resource.TestCheckResourceAttrPair(dataSourceName, "set_rule_uuid", resourceName, "set_rule_uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_rule_uuids.#", resourceName, "match_rule_uuids.#"),
					resource.TestCheckResourceAttr(dataSourceName, "match_rule_uuids.#", "2"),
					resource.TestCheckResourceAttrPair(dataSourceName, "set_rule_uuid", "mso_tenant_policies_set_rule_policy.test_set_rule", "uuid"),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_1", "uuid"),
					resource.TestCheckTypeSetElemAttrPair(dataSourceName, "match_rule_uuids.*", "mso_tenant_policies_match_rule_policy.test_match_rule_2", "uuid"),
				),
			},
		},
		CheckDestroy: testCheckRouteMapPolicyContextDestroy,
	})
}

func testAccMSORouteMapPolicyContextDataSourceConfig() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_route_map_policy_route_control_context" "test" {
  template_id           = mso_tenant_policies_route_map_policy_route_control_context.test_context.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control_context.test_context.route_map_policy_uuid
  name                  = mso_tenant_policies_route_map_policy_route_control_context.test_context.name
}
`, testAccMSORouteMapPolicyContextConfig_WithRules())
}

func testAccMSORouteMapPolicyContextDataSourceConfigMissingContext() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_route_map_policy_route_control_context" "test" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.test_policy.uuid
  name                  = "tf_missing_context"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}

func testAccMSORouteMapPolicyContextDataSourceConfigMissingParent() string {
	return fmt.Sprintf(`%s
data "mso_tenant_policies_route_map_policy_route_control_context" "test" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.test_policy.template_id
  route_map_policy_uuid = "00000000-0000-0000-0000-000000000000"
  name                  = "tf_missing_context"
}
`, testAccMSORouteMapPolicyContextConfigBase())
}
