package mso

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSORouteMapPolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_route_control.test_policy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Task: Create Route Map Policy") },
				Config:    testAccMSORouteMapPolicyConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_policy"),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Update Route Map Policy Description") },
				Config:    testAccMSORouteMapPolicyConfigUpdateDescription(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
				),
			},
			{
				PreConfig: func() { fmt.Println("Task: Update Route Map Policy Name") },
				Config:    testAccMSORouteMapPolicyConfigUpdateName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_policy_renamed"),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Task: Import Route Map Policy") },
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_route_map_policy_route_control", "routeMapP"),
	})
}

func testAccMSORouteMapPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = "Initial description"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigUpdateDescription() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy"
  description = "Updated description"
}
`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSORouteMapPolicyConfigUpdateName() string {
	return fmt.Sprintf(`%s
resource "mso_tenant_policies_route_map_policy_route_control" "test_policy" {
  template_id = mso_template.template_tenant.id
  name        = "test_policy_renamed"
  description = "Updated description"
}
`, testAccMSOTemplateResourceTenantConfig())
}
