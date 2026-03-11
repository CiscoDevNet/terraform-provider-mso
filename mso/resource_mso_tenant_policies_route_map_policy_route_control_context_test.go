package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSORouteMapPolicyContextResource(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_route_control_context.test_context"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Task: Creating Route Map Policy Context") },
				Config:    testAccMSORouteMapPolicyContextConfig_Create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName, "action", "permit"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Updating Context Attributes") },
				Config:    testAccMSORouteMapPolicyContextConfig_UpdateAttrs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "order", "2"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Adding Match and Set Rules") },
				Config:    testAccMSORouteMapPolicyContextConfig_WithRules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_rule_uuid", "a7f90577-fcd9-4c55-8bf4-628592922fc2"),
					resource.TestCheckResourceAttr(resourceName, "match_rules.#", "2"),
					CustomTestCheckTypeSetElemAttrs(resourceName, "match_rules", map[string]string{
						"uuid": "f016b945-83c3-49a8-ab1f-c0c69ed58ec8",
					}),
					CustomTestCheckTypeSetElemAttrs(resourceName, "match_rules", map[string]string{
						"uuid": "50f1e25e-b5b3-4a06-a1d8-a7e41120bf4e",
					}),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Removing Rules") },
				Config:    testAccMSORouteMapPolicyContextConfig_Create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "set_rule_uuid", ""),
					resource.TestCheckResourceAttr(resourceName, "match_rules.#", "0"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Task: Importing Context") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyChildWithArguments(
			"mso_tenant_policies_route_map_policy_route_control_context",
			"routeMap",
			"ctxs",
		),
	})
}

func TestAccMSORouteMapPolicyRouteControlContext_MultipleContextIndexing(t *testing.T) {
	resourceName1 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_1"
	resourceName2 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_2"
	resourceName3 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_3"
	resourceName4 := "mso_tenant_policies_route_map_policy_route_control_context.ctx_4"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Task: Creating 3 Route Control Contexts") },
				Config:    testAccRouteControlContextMultipleCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "permit"),
					resource.TestCheckResourceAttr(resourceName2, "name", "ctx_2"),
					resource.TestCheckResourceAttr(resourceName2, "description", "Second context"),
					resource.TestCheckResourceAttr(resourceName2, "order", "1"),
					resource.TestCheckResourceAttr(resourceName2, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context"),
					resource.TestCheckResourceAttr(resourceName3, "order", "2"),
					resource.TestCheckResourceAttr(resourceName3, "action", "permit"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Deleting Middle Context (ctx_2)") },
				Config:    testAccRouteControlContextMiddleDeleted(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "permit"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context"),
					resource.TestCheckResourceAttr(resourceName3, "order", "2"),
					resource.TestCheckResourceAttr(resourceName3, "action", "permit"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Updating Remaining Contexts After Middle Deletion") },
				Config:    testAccRouteControlContextUpdateAfterDeletion(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context updated"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context updated"),
					resource.TestCheckResourceAttr(resourceName3, "order", "5"),
					resource.TestCheckResourceAttr(resourceName3, "action", "deny"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Re-adding New Context (ctx_4) Alongside Existing Contexts") },
				Config:    testAccRouteControlContextReaddNew(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName1, "name", "ctx_1"),
					resource.TestCheckResourceAttr(resourceName1, "description", "First context updated"),
					resource.TestCheckResourceAttr(resourceName1, "order", "0"),
					resource.TestCheckResourceAttr(resourceName1, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName3, "name", "ctx_3"),
					resource.TestCheckResourceAttr(resourceName3, "description", "Third context updated"),
					resource.TestCheckResourceAttr(resourceName3, "order", "5"),
					resource.TestCheckResourceAttr(resourceName3, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName4, "name", "ctx_4"),
					resource.TestCheckResourceAttr(resourceName4, "description", "Fourth context new"),
					resource.TestCheckResourceAttr(resourceName4, "order", "3"),
					resource.TestCheckResourceAttr(resourceName4, "action", "permit"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Task: Importing Context (ctx_3) After Index Shifts") },
				ResourceName:      resourceName3,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyChildWithArguments(
			"mso_tenant_policies_route_map_policy_route_control_context",
			"routeMap",
			"ctxs",
		),
	})
}

func testAccMSORouteMapPolicyContextConfig_Create() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_1"
  order       = 1
  action      = "permit"
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfig_UpdateAttrs() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_1"
  description = "Updated description"
  order       = 2
  action      = "permit"
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccMSORouteMapPolicyContextConfig_WithRules() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "test_context" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_1"
  order       = 1
  action      = "permit"

  set_rule_uuid = "a7f90577-fcd9-4c55-8bf4-628592922fc2"

  match_rules {
    uuid = "f016b945-83c3-49a8-ab1f-c0c69ed58ec8"
  }
  match_rules {
    uuid = "50f1e25e-b5b3-4a06-a1d8-a7e41120bf4e"
  }
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccRouteControlContextMultipleCreate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_1"
  description = "First context"
  order       = 0
  action      = "permit"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_2" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_2"
  description = "Second context"
  order       = 1
  action      = "deny"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_3"
  description = "Third context"
  order       = 2
  action      = "permit"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control_context.ctx_2]
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccRouteControlContextMiddleDeleted() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_1"
  description = "First context"
  order       = 0
  action      = "permit"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_3"
  description = "Third context"
  order       = 2
  action      = "permit"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccRouteControlContextUpdateAfterDeletion() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_1"
  description = "First context updated"
  order       = 0
  action      = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_3"
  description = "Third context updated"
  order       = 5
  action      = "deny"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
}
`, testAccMSORouteMapPolicyConfigCreate())
}

func testAccRouteControlContextReaddNew() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_1" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_1"
  description = "First context updated"
  order       = 0
  action      = "deny"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_3" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_3"
  description = "Third context updated"
  order       = 5
  action      = "deny"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control_context.ctx_1]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "ctx_4" {
  parent_id   = mso_tenant_policies_route_map_policy_route_control.test_policy.id
  name        = "ctx_4"
  description = "Fourth context new"
  order       = 3
  action      = "permit"
  depends_on  = [mso_tenant_policies_route_map_policy_route_control_context.ctx_3]
}
`, testAccMSORouteMapPolicyConfigCreate())
}
