package mso

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMSOTenantPoliciesMcastRouteMapPolicyResource(t *testing.T) {
	resourceName := "mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast"
	var templateID, uuid string

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
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries",
						map[string]string{
							"order":               "1",
							"group_ip":            "226.2.2.2/8",
							"source_ip":           "1.1.1.1/1",
							"rendezvous_point_ip": "1.1.1.2",
							"action":              "permit",
						},
					),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok || rs.Primary == nil {
							return fmt.Errorf("resource %s not found in state", resourceName)
						}
						templateID = rs.Primary.Attributes["template_id"]
						uuid = rs.Primary.Attributes["uuid"]
						if templateID == "" || uuid == "" {
							return fmt.Errorf("resource %s is missing its template ID or UUID", resourceName)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Recreate Route Map Policy for Multicast after out-of-band deletion")
					if err := testAccDeletePolicyOutOfBand(testAccPreCheck(t), templateID, uuid, "tenantPolicyTemplate", "template", "mcastRouteMapPolicies"); err != nil {
						t.Fatalf("delete %s out of band: %v", resourceName, err)
					}
				},
				Config: testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "name", "tf_test_route_map_policy_multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "description", "Terraform test Route Map Policy for Multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries",
						map[string]string{
							"order":               "1",
							"group_ip":            "226.2.2.2/8",
							"source_ip":           "1.1.1.1/1",
							"rendezvous_point_ip": "1.1.1.2",
							"action":              "permit",
						},
					),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Test: Create Route Map Policy for Multicast with Invalid order value in Route Map Entries")
				},
				Config:                    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigCreateWithInvalidOrder(),
				Destroy:                   false,
				PreventPostDestroyRefresh: true,
				ExpectError:               regexp.MustCompile(`expected route_map_multicast_entries\.0\.order to be in the range \(0 - 65535\), got 65536`),
			},
			{
				PreConfig: func() { fmt.Println("Test: Update Route Map Policy for Multicast adding extra entry") },
				Config:    testAccMSOTenantPoliciesMcastRouteMapPolicyConfigUpdateAddingExtraEntry(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "name", "tf_test_route_map_policy_multicast"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "description", "Terraform test Route Map Policy for Multicast adding extra entry"),
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries.#", "2"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries",
						map[string]string{
							"order":               "1",
							"group_ip":            "226.2.2.2/8",
							"source_ip":           "1.1.1.1/1",
							"rendezvous_point_ip": "1.1.1.2",
							"action":              "permit",
						},
					),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries",
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
					resource.TestCheckResourceAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries.#", "1"),
					customTestCheckResourceTypeSetAttr("mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast", "route_map_multicast_entries",
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
			{
				PreConfig:    func() { fmt.Println("Test: Import missing Route Map Policy for Multicast") },
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok || rs.Primary == nil {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}
					separator := strings.LastIndex(rs.Primary.ID, "/")
					if separator < 0 || separator == len(rs.Primary.ID)-1 {
						return "", fmt.Errorf("resource %s has an invalid import ID", resourceName)
					}
					return rs.Primary.ID[:separator+1] + "tf_missing_oob", nil
				},
				ExpectError: regexp.MustCompile(`cannot import .*: resource not found`),
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
		route_map_multicast_entries {
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
		route_map_multicast_entries {
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
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
		route_map_multicast_entries {
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
		route_map_multicast_entries {
			order                   = 1
			group_ip                = "226.2.2.2/8"
			source_ip               = "1.1.1.1/1"
			rendezvous_point_ip     = "1.1.1.2"
			action                  = "permit"
		}
	}`, testAccMSOTemplateResourceTenantConfig())
}
