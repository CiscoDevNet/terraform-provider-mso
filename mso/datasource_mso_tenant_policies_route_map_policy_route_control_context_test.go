package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSORouteMapPolicyContextDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Task: Create Resource and Query via Data Source") },
				Config:    testAccMSORouteMapPolicyContextDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.mso_tenant_policies_route_map_policy_route_control_context.test", "parent_id",
						"mso_tenant_policies_route_map_policy_route_control_context.test", "parent_id",
					),
					resource.TestCheckResourceAttrPair(
						"data.mso_tenant_policies_route_map_policy_route_control_context.test", "name",
						"mso_tenant_policies_route_map_policy_route_control_context.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.mso_tenant_policies_route_map_policy_route_control_context.test", "description",
						"mso_tenant_policies_route_map_policy_route_control_context.test", "description",
					),
					resource.TestCheckResourceAttrPair(
						"data.mso_tenant_policies_route_map_policy_route_control_context.test", "action",
						"mso_tenant_policies_route_map_policy_route_control_context.test", "action",
					),
					resource.TestCheckResourceAttrPair(
						"data.mso_tenant_policies_route_map_policy_route_control_context.test", "order",
						"mso_tenant_policies_route_map_policy_route_control_context.test", "order",
					),
					resource.TestCheckResourceAttrPair(
						"data.mso_tenant_policies_route_map_policy_route_control_context.test", "set_rule_uuid",
						"mso_tenant_policies_route_map_policy_route_control_context.test", "set_rule_uuid",
					),
					resource.TestCheckResourceAttrPair(
						"data.mso_tenant_policies_route_map_policy_route_control_context.test", "match_rules.#",
						"mso_tenant_policies_route_map_policy_route_control_context.test", "match_rules.#",
					),
				),
			},
		},
	})
}

func testAccMSORouteMapPolicyContextDataSourceConfig() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test" {
  template_id = mso_template.template_tenant.id
  name        = "test_route_map_ds"
  description = "Route Map for data source test"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "test" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test.id
  name        = "ctx_ds_1"
  description = "Context for data source test"
  action      = "permit"
  order       = 1
}

data "mso_tenant_policies_route_map_policy_route_control_context" "test" {
  parent_id = mso_tenant_policies_route_map_policy_route_control_context.test.parent_id
  name      = "ctx_ds_1"
}
`, testAccMSOTemplateResourceTenantConfig())
}
