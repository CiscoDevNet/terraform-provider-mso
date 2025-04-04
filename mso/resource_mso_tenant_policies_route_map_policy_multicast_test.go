package mso

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMSOTenantPoliciesMcastRouteMapPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { fmt.Println("Test: Create Route Map Policy for Multicast") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "name", "tf_test_route_map_policy_multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "description", "Terraform test Route Map Policy for Multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_entries_multicast.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_entries_multicast",
						map[string]string{
							"order":               "1",
							"group_ip":            "226.2.2.2/8",
							"source_ip":           "1.1.1.1/1",
							"rendezvous_point_ip": "1.1.1.2",
							"action":              "permit",
						},
					),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Create Route Map Policy for Multicast with Invalid order value in Route Map Entries")
				},
				Config:                    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreateWithInvalidOrder(),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				ExpectError:               regexp.MustCompile(`config is invalid: expected route_map_entries_multicast\.0\.order to be in the range \(0 - 65535\), got 65536`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy for Multicast adding extra entry") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateAddingExtraEntry(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "name", "tf_test_route_map_policy_multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "description", "Terraform test Route Map Policy for Multicast adding extra entry"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_entries_multicast.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_entries_multicast",
						map[string]string{
							"order":               "1",
							"group_ip":            "226.2.2.2/8",
							"source_ip":           "1.1.1.1/1",
							"rendezvous_point_ip": "1.1.1.2",
							"action":              "permit",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_entries_multicast",
						map[string]string{
							"order":               "2",
							"group_ip":            "226.3.3.3/24",
							"source_ip":           "2.2.2.2/2",
							"rendezvous_point_ip": "2.2.2.3",
							"action":              "deny",
						},
					),
				),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy for Multicast removing extra entry") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateRemovingExtraEntry(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "name", "tf_test_route_map_policy_multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "description", "Terraform test Route Map Policy for Multicast removing extra entry"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_entries_multicast.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_entries_multicast",
						map[string]string{
							"order":               "1",
							"group_ip":            "226.2.2.2/8",
							"source_ip":           "1.1.1.1/1",
							"rendezvous_point_ip": "1.1.1.2",
							"action":              "permit",
						},
					),
				),
			},
			{
				PreConfig:         func() { fmt.Println("Test: Import Route Map Policy for Multicast") },
				ResourceName:      "mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testCheckResourceDestroyPolicyWithArguments("mso_tenant_policies_route_map_policy_multicast", "mcastRouteMap"),
	})
}

func testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreate() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
		template_id = mso_template.template_tenant.id
		name        = "tf_test_route_map_policy_multicast"
		description = "Terraform test Route Map Policy for Multicast"
		route_map_entries_multicast {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreateWithInvalidOrder() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast_error" {
		template_id = mso_template.template_tenant.id
		name        = "tf_test_route_map_policy_multicast_error"
		description = "Terraform test Route Map Policy for Multicast with invalid order"
		route_map_entries_multicast {
			order                   = 65536
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateAddingExtraEntry() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
		template_id = mso_template.template_tenant.id
		name        = "tf_test_route_map_policy_multicast"
		description = "Terraform test Route Map Policy for Multicast adding extra entry"
		route_map_entries_multicast {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
		route_map_entries_multicast {
			order                   = 2
			group_ip                = "226.3.3.3/24"
			source_ip               = "2.2.2.2/2"
			rendezvous_point_ip     = "2.2.2.3"
			action                  = "deny"
		}
	}`, testAccMSOTemplateResourceTenantConfig())
}

func testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateRemovingExtraEntry() string {
	return fmt.Sprintf(`%s
	resource "mso_tenant_policies_route_map_policy_multicast" "route_map_policy_multicast" {
		template_id = mso_template.template_tenant.id
		name        = "tf_test_route_map_policy_multicast"
		description = "Terraform test Route Map Policy for Multicast removing extra entry"
		route_map_entries_multicast {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}`, testAccMSOTemplateResourceTenantConfig())
}
